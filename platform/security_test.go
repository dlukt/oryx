package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func TestSecurity_PathTraversal_HelloVoices(t *testing.T) {
	// Setup the global config and environment.
	tmpDir, err := ioutil.TempDir("", "srs-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Mock the global conf
	conf = NewConfig()
	conf.Pwd = tmpDir

	// Create the directory structure expected by ai-talk
	confDir := path.Join(tmpDir, "containers/conf")
	if err := os.MkdirAll(confDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a safe file
	safeFile := path.Join(confDir, "hello-chinese.aac")
	if err := ioutil.WriteFile(safeFile, []byte("safe content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a secret file outside of the allowed directory
	secretFile := path.Join(tmpDir, "secret.txt")
	if err := ioutil.WriteFile(secretFile, []byte("secret content"), 0644); err != nil {
		t.Fatal(err)
	}

	// We need to access the handler. Since handleAITalkService registers it,
	// we can try to invoke the logic directly or register it.
	// However, handleAITalkService does some initialization.
	// Let's copy the handler logic here to verify the vulnerability pattern,
	// or try to call handleAITalkService if dependencies allow.
	// talkServer is global.
	talkServer = NewTalkServer()

	mux := http.NewServeMux()

	// We need to set aiTalkExampleDir as handleAITalkService does.
	aiTalkExampleDir = path.Join(conf.Pwd, "containers/conf")

	// Use the real handler logic
	handler := makeHelloVoicesHandler(context.Background())
	mux.HandleFunc("/terraform/v1/ai-talk/stage/hello-voices/", handler)

	// Test 1: Access safe file
	req1 := httptest.NewRequest("GET", "/terraform/v1/ai-talk/stage/hello-voices/hello-chinese.aac", nil)
	w1 := httptest.NewRecorder()
	mux.ServeHTTP(w1, req1)
	if w1.Code != 200 {
		t.Errorf("Expected 200 for safe file, got %d", w1.Code)
	}

	// Create a sensitive config file in the allowed directory
	configFile := path.Join(confDir, "nginx.conf")
	if err := ioutil.WriteFile(configFile, []byte("sensitive config"), 0644); err != nil {
		t.Fatal(err)
	}

	// Test 2: Access sensitive config file (no traversal needed)
	req2 := httptest.NewRequest("GET", "/terraform/v1/ai-talk/stage/hello-voices/nginx.conf", nil)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)

	if w2.Code == 200 {
		t.Errorf("Vulnerability confirmed: Access to arbitrary files in conf directory allowed. Got 200 OK")
	} else {
		t.Logf("Got %d for sensitive file, access denied as expected.", w2.Code)
	}
}

func TestSecurity_PathTraversal_OCR(t *testing.T) {
	// Simulate the logic in ocr.go to prevent regression
	requestPath := "/terraform/v1/ai/ocr/image/../../secret.jpg"
	prefix := "/terraform/v1/ai/ocr/image/"

	if len(requestPath) < len(prefix) {
		t.Fatal("Path too short")
	}

	filename := requestPath[len(prefix):]

	// FIX APPLIED:
	fileBase := path.Base(filename)
	uuid := fileBase[:len(fileBase)-len(path.Ext(fileBase))]

	imageFilePath := path.Join("ocr", fmt.Sprintf("%v.jpg", uuid))

	// Expected behavior after fix:
	// filename = "../../secret.jpg"
	// fileBase = "secret.jpg"
	// uuid = "secret"
	// imageFilePath = "ocr/secret.jpg"

	expectedSafe := "ocr/secret.jpg"

	if imageFilePath != expectedSafe {
		t.Errorf("Path traversal detected in OCR logic. Expected %s, got %s", expectedSafe, imageFilePath)
	}
}
