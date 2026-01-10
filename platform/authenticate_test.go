package main

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

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
