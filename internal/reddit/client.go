package reddit

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"github.com/lloydhumphreys/radcli/internal/config"
)

const (
	adsBaseURL  = "https://ads-api.reddit.com/api/v3"
	authBaseURL = "https://www.reddit.com"
	openAPIURL  = "https://ads-api.reddit.com/api/v3/openapi.json"
)

type Client struct {
	store      *config.Store
	httpClient *http.Client
}

type TokenSession struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	Scope        string
	ExpiresAt    time.Time
}

func New(store *config.Store) *Client {
	return &Client{
		store: store,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *Client) BuildAuthorizationURL() (string, string, error) {
	app := c.store.Config.App
	if app == nil {
		return "", "", fmt.Errorf("no app credentials found. run `rad auth setup` first")
	}

	state, err := randomState()
	if err != nil {
		return "", "", err
	}

	c.store.Config.Auth.PendingState = state
	if err := c.store.Save(); err != nil {
		return "", "", err
	}

	values := url.Values{}
	values.Set("client_id", app.ClientID)
	values.Set("response_type", "code")
	values.Set("state", state)
	values.Set("redirect_uri", app.RedirectURI)
	values.Set("duration", "permanent")
	values.Set("scope", strings.Join(app.Scopes, ","))

	return authBaseURL + "/api/v1/authorize?" + values.Encode(), state, nil
}

func (c *Client) ReportFields(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, openAPIURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("could not fetch OpenAPI metadata (%d)", resp.StatusCode)
	}

	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	paths, _ := payload["paths"].(map[string]any)
	reportsPath, _ := paths["/ad_accounts/{ad_account_id}/reports"].(map[string]any)
	post, _ := reportsPath["post"].(map[string]any)
	requestBody, _ := post["requestBody"].(map[string]any)
	content, _ := requestBody["content"].(map[string]any)
	applicationJSON, _ := content["application/json"].(map[string]any)
	schema, _ := applicationJSON["schema"].(map[string]any)
	properties, _ := schema["properties"].(map[string]any)
	data, _ := properties["data"].(map[string]any)
	dataProperties, _ := data["properties"].(map[string]any)
	fields, _ := dataProperties["fields"].(map[string]any)
	items, _ := fields["items"].(map[string]any)
	enumValues, _ := items["enum"].([]any)

	out := make([]string, 0, len(enumValues))
	for _, raw := range enumValues {
		value, ok := raw.(string)
		if ok && value != "" {
			out = append(out, value)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("could not extract report fields from OpenAPI metadata")
	}
	return out, nil
}

func (c *Client) ExchangeAuthorizationCode(ctx context.Context, code string) (*TokenSession, error) {
	app := c.store.Config.App
	if app == nil {
		return nil, fmt.Errorf("no app credentials found. run `rad auth setup` first")
	}

	values := url.Values{}
	values.Set("grant_type", "authorization_code")
	values.Set("code", code)
	values.Set("redirect_uri", app.RedirectURI)

	session, err := c.tokenRequest(ctx, app, values)
	if err != nil {
		return nil, err
	}

	c.store.Config.Auth.AccessToken = session.AccessToken
	c.store.Config.Auth.RefreshToken = session.RefreshToken
	c.store.Config.Auth.TokenType = session.TokenType
	c.store.Config.Auth.Scope = session.Scope
	c.store.Config.Auth.ExpiresAt = session.ExpiresAt.Format(time.RFC3339Nano)
	c.store.Config.Auth.PendingState = ""
	if err := c.store.Save(); err != nil {
		return nil, err
	}

	return session, nil
}

func (c *Client) RequestJSON(ctx context.Context, method, path string, query url.Values, body any) (map[string]any, error) {
	accessToken, err := c.validAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	u := adsBaseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	return c.doJSON(ctx, method, u, "Bearer "+accessToken, body)
}

func (c *Client) RequestPaginatedJSON(ctx context.Context, method, path string, query url.Values, body any, fetchAll bool) (map[string]any, error) {
	response, err := c.RequestJSON(ctx, method, path, query, body)
	if err != nil {
		return nil, err
	}
	if !fetchAll {
		return response, nil
	}

	for {
		next := nextURL(response)
		if next == "" {
			return response, nil
		}
		accessToken, err := c.validAccessToken(ctx)
		if err != nil {
			return nil, err
		}
		page, err := c.doJSON(ctx, method, next, "Bearer "+accessToken, body)
		if err != nil {
			return nil, err
		}
		response = mergePage(response, page)
	}
}

func (c *Client) validAccessToken(ctx context.Context) (string, error) {
	if c.store.Config.Auth.AccessToken == "" {
		return "", fmt.Errorf("no saved auth session. run `rad auth login` and `rad auth complete`")
	}

	if c.store.Config.Auth.ExpiresAt != "" {
		expiresAt, err := time.Parse(time.RFC3339Nano, c.store.Config.Auth.ExpiresAt)
		if err == nil && time.Until(expiresAt) <= time.Minute {
			session, err := c.refreshToken(ctx)
			if err != nil {
				return "", err
			}
			return session.AccessToken, nil
		}
	}

	return c.store.Config.Auth.AccessToken, nil
}

func (c *Client) refreshToken(ctx context.Context) (*TokenSession, error) {
	app := c.store.Config.App
	if app == nil {
		return nil, fmt.Errorf("no app credentials found. run `rad auth setup` first")
	}
	if c.store.Config.Auth.RefreshToken == "" {
		return nil, fmt.Errorf("access token expired and no refresh token is saved. re-authenticate with `rad auth login`")
	}

	values := url.Values{}
	values.Set("grant_type", "refresh_token")
	values.Set("refresh_token", c.store.Config.Auth.RefreshToken)

	session, err := c.tokenRequest(ctx, app, values)
	if err != nil {
		return nil, err
	}
	if session.RefreshToken == "" {
		session.RefreshToken = c.store.Config.Auth.RefreshToken
	}

	c.store.Config.Auth.AccessToken = session.AccessToken
	c.store.Config.Auth.RefreshToken = session.RefreshToken
	c.store.Config.Auth.TokenType = session.TokenType
	c.store.Config.Auth.Scope = session.Scope
	c.store.Config.Auth.ExpiresAt = session.ExpiresAt.Format(time.RFC3339Nano)
	if err := c.store.Save(); err != nil {
		return nil, err
	}

	return session, nil
}

func (c *Client) tokenRequest(ctx context.Context, app *config.AppCredentials, form url.Values) (*TokenSession, error) {
	credentials := base64.StdEncoding.EncodeToString([]byte(app.ClientID + ":" + app.ClientSecret))
	jsonBody, err := c.doFormJSON(ctx, authBaseURL+"/api/v1/access_token", "Basic "+credentials, app.UserAgent, form)
	if err != nil {
		return nil, err
	}

	accessToken, _ := jsonBody["access_token"].(string)
	tokenType, _ := jsonBody["token_type"].(string)
	scope, _ := jsonBody["scope"].(string)
	refreshToken, _ := jsonBody["refresh_token"].(string)
	expiresIn, ok := numberAsFloat(jsonBody["expires_in"])
	if !ok || accessToken == "" || tokenType == "" {
		return nil, fmt.Errorf("could not decode token response")
	}

	return &TokenSession{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    tokenType,
		Scope:        scope,
		ExpiresAt:    time.Now().Add(time.Duration(expiresIn) * time.Second),
	}, nil
}

func (c *Client) doJSON(ctx context.Context, method, rawURL, authorization string, body any) (map[string]any, error) {
	userAgent := ""
	if c.store.Config.App != nil {
		userAgent = c.store.Config.App.UserAgent
	}
	return c.doRequestJSON(ctx, method, rawURL, authorization, userAgent, "application/json", body)
}

func (c *Client) doFormJSON(ctx context.Context, rawURL, authorization, userAgent string, form url.Values) (map[string]any, error) {
	return c.doRequestJSON(ctx, http.MethodPost, rawURL, authorization, userAgent, "application/x-www-form-urlencoded", form.Encode())
}

func (c *Client) doRequestJSON(ctx context.Context, method, rawURL, authorization, userAgent, contentType string, body any) (map[string]any, error) {
	var reader io.Reader
	switch v := body.(type) {
	case nil:
		reader = nil
	case string:
		reader = strings.NewReader(v)
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, rawURL, reader)
	if err != nil {
		return nil, err
	}
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}
	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(payload) == 0 {
		return map[string]any{}, nil
	}

	var out map[string]any
	if err := json.Unmarshal(payload, &out); err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if message := detailedErrorMessage(out); message != "" {
			return nil, fmt.Errorf("api request failed (%d): %s", resp.StatusCode, message)
		}
		return nil, fmt.Errorf("api request failed (%d)", resp.StatusCode)
	}

	return out, nil
}

func OpenBrowser(rawURL string) error {
	cmd := exec.Command("/usr/bin/open", rawURL)
	return cmd.Run()
}

func nextURL(response map[string]any) string {
	pagination, _ := response["pagination"].(map[string]any)
	next, _ := pagination["next_url"].(string)
	return next
}

func mergePage(original, next map[string]any) map[string]any {
	if next == nil {
		return original
	}

	if lhs, ok := original["data"].([]any); ok {
		if rhs, ok := next["data"].([]any); ok {
			original["data"] = append(lhs, rhs...)
			original["pagination"] = next["pagination"]
			return original
		}
	}

	if lhs, ok := original["data"].(map[string]any); ok {
		if rhs, ok := next["data"].(map[string]any); ok {
			if leftMetrics, ok := lhs["metrics"].([]any); ok {
				if rightMetrics, ok := rhs["metrics"].([]any); ok {
					lhs["metrics"] = append(leftMetrics, rightMetrics...)
					original["data"] = lhs
					original["pagination"] = next["pagination"]
				}
			}
		}
	}

	return original
}

func detailedErrorMessage(payload map[string]any) string {
	errMap, _ := payload["error"].(map[string]any)
	message, _ := errMap["message"].(string)
	fields, _ := errMap["fields"].([]any)
	if len(fields) == 0 {
		return message
	}

	details := make([]string, 0, len(fields))
	for _, raw := range fields {
		fieldMap, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		field, _ := fieldMap["field"].(string)
		fieldMessage, _ := fieldMap["message"].(string)
		switch {
		case field != "" && fieldMessage != "":
			details = append(details, field+": "+fieldMessage)
		case fieldMessage != "":
			details = append(details, fieldMessage)
		case field != "":
			details = append(details, field)
		}
	}
	if len(details) == 0 {
		return message
	}
	if message == "" {
		return strings.Join(details, "; ")
	}
	return message + " (" + strings.Join(details, "; ") + ")"
}

func numberAsFloat(value any) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	default:
		return 0, false
	}
}

func randomState() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
