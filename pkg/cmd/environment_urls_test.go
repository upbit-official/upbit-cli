package cmd

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

type environmentTransport struct {
	requests []*http.Request
}

func (transport *environmentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	transport.requests = append(transport.requests, req)
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": {"application/json"}},
		Body:       io.NopCloser(strings.NewReader("[]")),
	}, nil
}

func TestEnvironmentRequestURLs(t *testing.T) {
	for _, tc := range []struct {
		name  string
		flags []string
		env   string
		host  string
	}{
		{"default", nil, "", "api.upbit.com"},
		{"kr", []string{"--environment", "kr"}, "", "api.upbit.com"},
		{"sg", []string{"--environment", "sg"}, "", "sg-api.upbit.com"},
		{"id", []string{"--environment", "id"}, "", "api.upbit.co.id"},
		{"th", []string{"--environment", "th"}, "", "th-api.upbit.com"},
		{"environment override", []string{"--environment", "id"}, "https://custom.example", "custom.example"},
		{"flag overrides environment", []string{"--environment", "id", "--base-url", "https://api.upbit.co.id"}, "https://custom.example", "api.upbit.co.id"},
		{"explicit old address preserved", []string{"--environment", "id", "--base-url", "https://id-api.upbit.com"}, "", "id-api.upbit.com"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("UPBIT_BASE_URL", tc.env)
			if tc.env == "" {
				if err := os.Unsetenv("UPBIT_BASE_URL"); err != nil {
					t.Fatal(err)
				}
			}
			transport := &environmentTransport{}
			original := http.DefaultTransport
			http.DefaultTransport = transport
			t.Cleanup(func() { http.DefaultTransport = original })
			flags := []string{"upbit", "--access-key", "test-key", "--secret-key", strings.Repeat("x", 32), "--format", "json"}
			flags = append(flags, tc.flags...)
			flags = append(flags, "trading-pairs", "list", "--header", "X-Upbit-Initiator: upbit-cli-skill/test")
			if err := Command.Run(context.Background(), flags); err != nil {
				t.Fatal(err)
			}
			if len(transport.requests) != 1 {
				t.Fatalf("expected one request, got %d", len(transport.requests))
			}
			req := transport.requests[0]
			if req.URL.Host != tc.host || req.URL.Path != "/v1/market/all" {
				t.Fatalf("unexpected URL: %s", req.URL)
			}
			if req.Header.Get("X-Upbit-Initiator") != "upbit-cli-skill/test" {
				t.Fatal("Skills initiator header was not preserved")
			}
		})
	}
}
