package main

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func TestAuthenticate_SigningMethod(t *testing.T) {
	apiSecret := "test-secret"

	// 1. Create a valid token with HS256 (expected).
	tokenHS256, err := func() (string, error) {
		createAt, expireAt := time.Now(), time.Now().Add(365*24*time.Hour)
		claims := struct {
			Version string `json:"v"`
			Nonce   string `json:"nonce"`
			jwt.RegisteredClaims
		}{
			Version: "1.0",
			Nonce:   "nonce",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expireAt),
				IssuedAt:  jwt.NewNumericDate(createAt),
			},
		}
		return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(apiSecret))
	}()
	if err != nil {
		t.Fatalf("Failed to create HS256 token: %v", err)
	}

	// 2. Create a token with 'none' alg.
	noneToken := jwt.New(jwt.SigningMethodNone)
	noneToken.Claims = jwt.MapClaims{
		"foo": "bar",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	tokenNone, _ := noneToken.SignedString(jwt.UnsafeAllowNoneSignatureType)

	tests := []struct {
		name      string
		token     string
		shouldErr bool
	}{
		{
			name:      "Valid HS256 Token",
			token:     tokenHS256,
			shouldErr: false,
		},
		{
			name:      "None Algorithm Token",
			token:     tokenNone,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := http.Header{}
			err := Authenticate(context.Background(), apiSecret, tt.token, header)

			if tt.shouldErr {
				if err == nil {
					t.Errorf("Expected error for %s, but got none", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("Expected success for %s, but got error: %v", tt.name, err)
				}
			}
		})
	}
}

func TestUtils_RebuildStreamURL(t *testing.T) {
	urlSamples := []struct {
		url     string
		rebuild string
	}{
		{url: "rtsp://121.1.2.3", rebuild: "rtsp://121.1.2.3"},
		{url: "rtsp://121.1.2.3/Streaming/Channels/101", rebuild: "rtsp://121.1.2.3/Streaming/Channels/101"},
		{url: "rtsp://121.1.2.3:554/Streaming/Channels/101", rebuild: "rtsp://121.1.2.3:554/Streaming/Channels/101"},
		{url: "rtsp://121.1.2.3:554/Streaming/Channels/101?k=v", rebuild: "rtsp://121.1.2.3:554/Streaming/Channels/101?k=v"},
		{url: "rtsp://CamViewer:abc123@121.1.2.3:554/Streaming/Channels/101", rebuild: "rtsp://CamViewer:abc123@121.1.2.3:554/Streaming/Channels/101"},
		{url: "rtsp://CamViewer:abc123?!@121.1.2.3:554/Streaming/Channels/101", rebuild: "rtsp://CamViewer:abc123%3F%21@121.1.2.3:554/Streaming/Channels/101"},
		{url: "rtsp://CamViewer:abc123@?!@121.1.2.3:554/Streaming/Channels/101", rebuild: "rtsp://CamViewer:abc123%40%3F%21@121.1.2.3:554/Streaming/Channels/101"},
		{url: "rtsp://CamViewer:abc123@?!@121.1.2.3:554/Streaming/Channels/101?k=v", rebuild: "rtsp://CamViewer:abc123%40%3F%21@121.1.2.3:554/Streaming/Channels/101?k=v"},
		{url: "rtsp://CamViewer:abc123@?!@121.1.2.3:554", rebuild: "rtsp://CamViewer:abc123%40%3F%21@121.1.2.3:554"},
		{url: "rtsp://Cam@Viewer:abc123@?!@121.1.2.3:554", rebuild: "rtsp://Cam%40Viewer:abc123%40%3F%21@121.1.2.3:554"},
		{url: "rtsp://CamViewer:abc123@?!~#$%^&*()_+-=\\|?@121.1.2.3:554/Streaming/Channels/101", rebuild: "rtsp://CamViewer:abc123%40%3F%21~%23$%25%5E&%2A%28%29_+-=%5C%7C%3F@121.1.2.3:554/Streaming/Channels/101"},
		{url: "rtsp://CamViewer:abc123@347?1!@121.1.2.3:554/Streaming/Channels/101", rebuild: "rtsp://CamViewer:abc123%40347%3F1%21@121.1.2.3:554/Streaming/Channels/101"},
		{url: "srt://213.171.194.158:10080", rebuild: "srt://213.171.194.158:10080"},
		{url: "srt://213.171.194.158:10080?streamid=#!::r=live/primary,latency=20,m=request", rebuild: "srt://213.171.194.158:10080?streamid=#!::r=live/primary,latency=20,m=request"},
	}
	for _, urlSample := range urlSamples {
		if r0, err := RebuildStreamURL(urlSample.url); err != nil {
			t.Errorf("Fail for err %+v", err)
			return
		} else if rebuild := r0.String(); rebuild != urlSample.rebuild {
			t.Errorf("rebuild url %v failed, expect %v, actual %v",
				urlSample.url, urlSample.rebuild, rebuild)
			return
		}
	}
}

func TestValidateServerURL(t *testing.T) {
	tests := []struct {
		name      string
		server    string
		shouldErr bool
	}{
		{
			name:      "Valid server",
			server:    "rtmp://localhost/live",
			shouldErr: false,
		},
		{
			name:      "Invalid server starts with dash",
			server:    "-f",
			shouldErr: true,
		},
		{
			name:      "Invalid server starts with double dash",
			server:    "--help",
			shouldErr: true,
		},
		{
			name:      "Invalid server file protocol",
			server:    "file:///tmp/output.flv",
			shouldErr: true,
		},
		{
			name:      "Invalid server http protocol",
			server:    "http://localhost/live",
			shouldErr: true,
		},
		{
			name:      "Invalid server exec protocol",
			server:    "exec://whoami",
			shouldErr: true,
		},
		{
			name:      "Valid server rtmps protocol",
			server:    "rtmps://localhost/live",
			shouldErr: false,
		},
		{
			name:      "Valid server srt protocol",
			server:    "srt://localhost/live",
			shouldErr: false,
		},
		{
			name:      "Valid server rtsp protocol",
			server:    "rtsp://localhost/live",
			shouldErr: false,
		},
		{
			name:      "Invalid server local file path",
			server:    "/etc/passwd",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateServerURL(tt.server)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("Expected error for server %v, but got none", tt.server)
				}
			} else {
				if err != nil {
					t.Errorf("Expected valid for server %v, but got error: %v", tt.server, err)
				}
			}
		})
	}
}

func TestUtils_ParseFFmpegLogs(t *testing.T) {
	for _, e := range []struct {
		log   string
		ts    string
		speed string
	}{
		{log: "time=00:10:09.138 speed=1x", ts: "00:10:09.138", speed: "1x"},
		{log: "size=18859kB time=00:10:09.138 speed=1x", ts: "00:10:09.138", speed: "1x"},
		{log: "size=18859kB time=00:10:09.138 speed=1x dup=1", ts: "00:10:09.138", speed: "1x"},
		{log: "size=18859kB time=00:10:09.138 bitrate=253.5kbits/s speed=1x dup=1", ts: "00:10:09.138", speed: "1x"},
		{log: "size=18859kB time=00:10:09.38 bitrate=253.5kbits/s speed=1x", ts: "00:10:09.38", speed: "1x"},
		{log: "frame=184 fps=9.7 q=28.0 size=364kB time=00:00:19.41 bitrate=153.7kbits/s dup=0 drop=235 speed=1.03x", ts: "00:00:19.41", speed: "1.03x"},
	} {
		if ts, speed, err := ParseFFmpegCycleLog(e.log); err != nil {
			t.Errorf("Fail parse %v for err %+v", e, err)
		} else if ts != e.ts {
			t.Errorf("Fail for ts %v of %v", ts, e)
		} else if speed != e.speed {
			t.Errorf("Fail for speed %v of %v", speed, e)
		}
	}
}

func TestAuthenticate_TimingAttack(t *testing.T) {
	// This test ensures the Authenticate function works correctly after
	// the fix for timing attack vulnerability.

	apiSecret := "supersecret"
	ctx := context.Background()

	tests := []struct {
		name        string
		apiSecret   string
		token       string
		header      http.Header
		wantErr     bool
		errContains string
	}{
		{
			name:      "No API Secret",
			apiSecret: "",
			token:     "",
			header:    http.Header{},
			wantErr:   true,
			errContains: "no api secret",
		},
		{
			name:      "No Auth",
			apiSecret: apiSecret,
			token:     "",
			header:    http.Header{},
			wantErr:   true,
			errContains: "no Authorization or token",
		},
		{
			name:      "Valid Bearer",
			apiSecret: apiSecret,
			token:     "",
			header:    http.Header{"Authorization": []string{"Bearer " + apiSecret}},
			wantErr:   false,
		},
		{
			name:      "Invalid Bearer Format",
			apiSecret: apiSecret,
			token:     "",
			header:    http.Header{"Authorization": []string{"Basic " + apiSecret}},
			wantErr:   true,
			errContains: "Invalid Authorization format",
		},
		{
			name:      "Invalid Bearer Secret",
			apiSecret: apiSecret,
			token:     "",
			header:    http.Header{"Authorization": []string{"Bearer wrongsecret"}},
			wantErr:   true,
			errContains: "invalid bearer token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Authenticate(ctx, tt.apiSecret, tt.token, tt.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Authenticate() error = %v, wantErr containing %v", err, tt.errContains)
				}
			}
		})
	}
}
