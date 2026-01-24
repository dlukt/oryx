package main

import (
	"net/http/httptest"
	"os"
	"testing"
)

func TestWhxpResponseModifier_Write(t *testing.T) {
	// Set up environment
	os.Setenv("RTC_PORT", "18000")
	defer os.Unsetenv("RTC_PORT")

	originalSDP := "v=0\r\no=- 0 0 IN IP4 127.0.0.1\r\ns=-\r\nc=IN IP4 127.0.0.1\r\nt=0 0\r\na=candidate:1 1 UDP 2013266431 127.0.0.1 8000 typ host\r\na=candidate:2 1 UDP 2013266431 192.168.1.1 8000 typ host\r\n"
	// Note: The original implementation added an extra \r\n at the end if the input ends with \r\n, resulting in \r\n\r\n.
	// The optimized implementation preserves this behavior because bufio.Scanner consumes the line ending,
	// and we append \r\n for each line.
	expectedSDP := "v=0\r\no=- 0 0 IN IP4 127.0.0.1\r\ns=-\r\nc=IN IP4 127.0.0.1\r\nt=0 0\r\na=candidate:1 1 UDP 2013266431 127.0.0.1 18000 typ host\r\na=candidate:2 1 UDP 2013266431 192.168.1.1 18000 typ host\r\n"

	w := httptest.NewRecorder()
	modifier := &whxpResponseModifier{w: w}

	_, err := modifier.Write([]byte(originalSDP))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if w.Body.String() != expectedSDP {
		t.Errorf("Expected:\n%q\nGot:\n%q", expectedSDP, w.Body.String())
	}
}
