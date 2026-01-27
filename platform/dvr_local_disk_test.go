package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func TestRecordWorker_TsHandler_Vulnerability(t *testing.T) {
	// Setup context
	ctx := context.Background()

	// Use a temp dir
	tempDir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalWd)

	// Create "record" directory (valid path)
	recordDir := path.Join(tempDir, "record", "valid-uuid")
	if err := os.MkdirAll(recordDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create valid ts file
	validTsFile := path.Join(recordDir, "valid-uuid.ts")
	if err := os.WriteFile(validTsFile, []byte("valid content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create "secret" directory (target for traversal)
	secretDir := path.Join(tempDir, "secret", "data")
	if err := os.MkdirAll(secretDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create secret ts file
	secretTsFile := path.Join(secretDir, "secret.ts")
	if err := os.WriteFile(secretTsFile, []byte("secret content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Initialize RecordWorker
	worker := NewRecordWorker()
	mux := http.NewServeMux()
	if err := worker.Handle(ctx, mux); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		path           string
		expectStatus   int
		expectContent  string
		expectSuccess  bool // Whether we expect the request to succeed (200 OK)
	}{
		{
			name:          "Valid Access",
			path:          "/terraform/v1/hooks/record/hls/record/valid-uuid/valid-uuid.ts",
			expectStatus:  http.StatusOK,
			expectContent: "valid content",
			expectSuccess: true,
		},
		{
			name:          "Exploit Access",
			path:          "/terraform/v1/hooks/record/hls/secret/data/secret.ts",
			expectStatus:  http.StatusOK,
			expectContent: "secret content",
			expectSuccess: false, // Should fail after fix
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)

			if tt.expectSuccess {
				if resp.StatusCode != http.StatusOK {
					t.Errorf("Expected status OK, got %v", resp.StatusCode)
				}
				if string(body) != tt.expectContent {
					t.Errorf("Expected content %q, got %q", tt.expectContent, string(body))
				}
			} else {
				if resp.StatusCode == http.StatusOK {
					t.Errorf("Expected error status, got OK. Content: %q", string(body))
				}
			}
		})
	}
}
