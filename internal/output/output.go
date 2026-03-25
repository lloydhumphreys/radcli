package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

func PrintJSON(w io.Writer, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "%s\n", data)
	return err
}

func PrintTable(w io.Writer, rows []map[string]string, preferred []string) error {
	if len(rows) == 0 {
		_, err := fmt.Fprintln(w, "No results.")
		return err
	}

	headers := headersFor(rows, preferred)
	widths := make(map[string]int, len(headers))
	for _, h := range headers {
		widths[h] = len(h)
		for _, row := range rows {
			if l := len(row[h]); l > widths[h] {
				widths[h] = l
			}
		}
	}

	if _, err := fmt.Fprintln(w, formatValues(headers, headers, widths)); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, divider(headers, widths)); err != nil {
		return err
	}
	for _, row := range rows {
		if _, err := fmt.Fprintln(w, formatRow(headers, row, widths)); err != nil {
			return err
		}
	}
	return nil
}

func PrintCSV(w io.Writer, rows []map[string]string, preferred []string) error {
	if len(rows) == 0 {
		writer := csv.NewWriter(w)
		writer.Flush()
		return writer.Error()
	}

	headers := headersFor(rows, preferred)
	writer := csv.NewWriter(w)
	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, row := range rows {
		record := make([]string, 0, len(headers))
		for _, header := range headers {
			record = append(record, row[header])
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	writer.Flush()
	return writer.Error()
}

func headersFor(rows []map[string]string, preferred []string) []string {
	if len(preferred) > 0 {
		headers := make([]string, 0, len(preferred))
		for _, key := range preferred {
			for _, row := range rows {
				if _, ok := row[key]; ok {
					headers = append(headers, key)
					break
				}
			}
		}
		if len(headers) > 0 {
			return headers
		}
	}

	seen := map[string]struct{}{}
	for _, row := range rows {
		for key := range row {
			seen[key] = struct{}{}
		}
	}
	headers := make([]string, 0, len(seen))
	for key := range seen {
		headers = append(headers, key)
	}
	sort.Strings(headers)
	return headers
}

func formatRow(headers []string, row map[string]string, widths map[string]int) string {
	if len(headers) == 0 {
		return ""
	}
	values := make([]string, 0, len(headers))
	for _, header := range headers {
		values = append(values, row[header])
	}
	return formatValues(headers, values, widths)
}

func formatValues(headers []string, values []string, widths map[string]int) string {
	parts := make([]string, 0, len(values))
	for index, value := range values {
		header := headers[index]
		width := widths[header]
		parts = append(parts, value+strings.Repeat(" ", max(0, width-len(value))))
	}
	return strings.Join(parts, "  ")
}

func divider(headers []string, widths map[string]int) string {
	parts := make([]string, 0, len(headers))
	for _, header := range headers {
		parts = append(parts, strings.Repeat("-", widths[header]))
	}
	return strings.Join(parts, "  ")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
