package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLatestReleaseTag(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		want       string
		wantErr    bool
		errMessage string
	}{
		{
			name: "successful release fetch",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(githubRelease{TagName: "v1.2.3"})
			},
			want:    "1.2.3",
			wantErr: false,
		},
		{
			name: "release without v prefix",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(githubRelease{TagName: "2.0.0"})
			},
			want:    "2.0.0",
			wantErr: false,
		},
		{
			name: "network error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "server error", http.StatusInternalServerError)
			},
			wantErr:    true,
			errMessage: "unexpected status 500",
		},
		{
			name: "empty tag response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(githubRelease{TagName: ""})
			},
			wantErr:    true,
			errMessage: "empty release tag",
		},
		{
			name: "invalid JSON",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte("invalid"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			originalURL := latestReleaseURL
			latestReleaseURL = server.URL
			defer func() { latestReleaseURL = originalURL }()

			got, err := getLatestReleaseTag()
			if (err != nil) != tt.wantErr {
				t.Errorf("getLatestReleaseTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMessage != "" {
				if err == nil || !contains(err.Error(), tt.errMessage) {
					t.Errorf("getLatestReleaseTag() error = %v, want error containing %v", err, tt.errMessage)
				}
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("getLatestReleaseTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
