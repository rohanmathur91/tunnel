package server

import (
	"net/http/httptest"
	"testing"
)

func newTestServer() *Server {
	config := LoadConfig()
	return New(&config)
}

func TestExtractTunnelId(t *testing.T) {
	server := newTestServer()

	tests := map[string]struct {
		host     string
		expected string
	}{
		"without subdomain": {
			host:     "http://localhost:8000",
			expected: "",
		},
		"subdomain with only numbers": {
			host:     "http://123.localhost:8000",
			expected: "123",
		},
		"subdomain with only letters": {
			host:     "http://abc.localhost:8000",
			expected: "abc",
		},
		"subdomain mixed": {
			host:     "http://1a2a3a.localhost:8000",
			expected: "1a2a3a",
		},
		"subdomain with uuid": {
			host:     "http://dab09a8c-4b62-4e61-9dfa-083349f8da8f.localhost:8000",
			expected: "dab09a8c-4b62-4e61-9dfa-083349f8da8f",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.host, nil)
			got := server.extractTunnelID(req)
			if got != tt.expected {
				t.Errorf("extractTunnelID(%v) = %v; want %v", tt.host, got, tt.expected)
			}
		})
	}

}
