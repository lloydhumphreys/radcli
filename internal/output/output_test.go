package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintTableNoResults(t *testing.T) {
	var buf bytes.Buffer
	if err := PrintTable(&buf, nil, nil); err != nil {
		t.Fatalf("PrintTable() error = %v", err)
	}
	if got := buf.String(); got != "No results.\n" {
		t.Fatalf("PrintTable() output = %q, want %q", got, "No results.\n")
	}
}

func TestPrintTableUsesPreferredHeaderOrder(t *testing.T) {
	rows := []map[string]string{
		{"id": "1", "name": "horseflaps", "status": "ACTIVE"},
		{"id": "22", "name": "beta", "status": "PAUSED"},
	}

	var buf bytes.Buffer
	if err := PrintTable(&buf, rows, []string{"name", "id", "status"}); err != nil {
		t.Fatalf("PrintTable() error = %v", err)
	}

	lines := strings.Split(strings.TrimSuffix(buf.String(), "\n"), "\n")
	if len(lines) < 4 {
		t.Fatalf("expected at least 4 lines, got %d: %q", len(lines), buf.String())
	}

	nameIndex := strings.Index(lines[0], "name")
	idIndex := strings.Index(lines[0], "id")
	statusIndex := strings.Index(lines[0], "status")
	if !(nameIndex >= 0 && idIndex > nameIndex && statusIndex > idIndex) {
		t.Fatalf("header order mismatch: %q", lines[0])
	}
}

func TestPrintCSVUsesPreferredHeaderOrder(t *testing.T) {
	rows := []map[string]string{
		{"id": "1", "name": "horseflaps", "status": "ACTIVE"},
	}

	var buf bytes.Buffer
	if err := PrintCSV(&buf, rows, []string{"name", "id", "status"}); err != nil {
		t.Fatalf("PrintCSV() error = %v", err)
	}

	want := "name,id,status\nhorseflaps,1,ACTIVE\n"
	if got := buf.String(); got != want {
		t.Fatalf("PrintCSV() output = %q, want %q", got, want)
	}
}
