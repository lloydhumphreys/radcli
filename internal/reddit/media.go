package reddit

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const oauthBaseURL = "https://oauth.reddit.com"

// MediaUploadLease holds the data returned by Reddit's media asset endpoint.
type MediaUploadLease struct {
	UploadURL string            // S3 URL to POST the file to
	Fields    map[string]string // form fields to include in the S3 upload
	AssetID   string
}

// RequestMediaUploadLease requests a presigned upload lease from Reddit's core API.
func (c *Client) RequestMediaUploadLease(ctx context.Context, filename, mimeType string) (*MediaUploadLease, error) {
	accessToken, err := c.validAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Set("filepath", filename)
	form.Set("mimetype", mimeType)

	userAgent := ""
	if c.store.Config.App != nil {
		userAgent = c.store.Config.App.UserAgent
	}

	reqURL := oauthBaseURL + "/api/media/asset.json"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var parsed map[string]any
		if json.Unmarshal(body, &parsed) == nil {
			if message := detailedErrorMessage(parsed); message != "" {
				return nil, fmt.Errorf("media lease request failed (%d): %s", resp.StatusCode, message)
			}
		}
		return nil, fmt.Errorf("media lease request failed (%d): %s", resp.StatusCode, truncate(string(body), 500))
	}

	return parseUploadLease(body)
}

// UploadMediaToS3 uploads a local file to the presigned S3 URL and returns
// the resulting hosted media URL.
func (c *Client) UploadMediaToS3(ctx context.Context, lease *MediaUploadLease, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Write all S3 form fields first (order matters for S3 presigned POSTs).
	for key, value := range lease.Fields {
		if err := writer.WriteField(key, value); err != nil {
			return "", fmt.Errorf("writing form field %q: %w", key, err)
		}
	}

	// Write the file field last.
	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, lease.UploadURL, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if message := s3ErrorMessage(respBody); message != "" {
			return "", fmt.Errorf("S3 upload failed (%d): %s", resp.StatusCode, message)
		}
		return "", fmt.Errorf("S3 upload failed (%d): %s", resp.StatusCode, truncate(string(respBody), 500))
	}

	return parseS3Location(respBody)
}

func parseUploadLease(body []byte) (*MediaUploadLease, error) {
	// Reddit returns JSON like:
	// {
	//   "args": {
	//     "action": "//reddit-uploaded-media.s3-accelerate.amazonaws.com/...",
	//     "fields": [{"name":"key","value":"..."}, ...]
	//   },
	//   "asset": { "asset_id": "...", "websocket_url": "wss://..." }
	// }
	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("could not parse upload lease response: %w", err)
	}

	args, _ := raw["args"].(map[string]any)
	if args == nil {
		return nil, fmt.Errorf("upload lease response missing 'args': %s", truncate(string(body), 300))
	}

	action, _ := args["action"].(string)
	if action == "" {
		return nil, fmt.Errorf("upload lease response missing 'args.action'")
	}

	// action is typically "//reddit-uploaded-media.s3-accelerate.amazonaws.com/..."
	if strings.HasPrefix(action, "//") {
		action = "https:" + action
	}

	fields := map[string]string{}
	rawFields, _ := args["fields"].([]any)
	for _, raw := range rawFields {
		entry, _ := raw.(map[string]any)
		name, _ := entry["name"].(string)
		value, _ := entry["value"].(string)
		if name != "" {
			fields[name] = value
		}
	}

	asset, _ := raw["asset"].(map[string]any)
	assetID, _ := asset["asset_id"].(string)

	return &MediaUploadLease{
		UploadURL: action,
		Fields:    fields,
		AssetID:   assetID,
	}, nil
}

func parseS3Location(body []byte) (string, error) {
	// S3 returns XML like:
	// <PostResponse>
	//   <Location>https://reddit-uploaded-media.s3-accelerate.amazonaws.com/...</Location>
	//   <Bucket>...</Bucket>
	//   <Key>...</Key>
	//   <ETag>...</ETag>
	// </PostResponse>
	var result struct {
		Location string `xml:"Location"`
	}
	if err := xml.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("could not parse S3 response: %w", err)
	}
	if result.Location == "" {
		return "", fmt.Errorf("S3 response missing Location: %s", truncate(string(body), 300))
	}

	// The Location URL is often URL-encoded; decode it.
	decoded, err := url.QueryUnescape(result.Location)
	if err != nil {
		return result.Location, nil // fallback to raw
	}
	return decoded, nil
}

func s3ErrorMessage(body []byte) string {
	var s3Err struct {
		Code    string `xml:"Code"`
		Message string `xml:"Message"`
	}
	if xml.Unmarshal(body, &s3Err) != nil || s3Err.Message == "" {
		return ""
	}
	if s3Err.Code != "" {
		return s3Err.Code + ": " + s3Err.Message
	}
	return s3Err.Message
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
