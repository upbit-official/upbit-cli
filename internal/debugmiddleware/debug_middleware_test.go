package debugmiddleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/require"
)

func TestDebugMiddleware(t *testing.T) {
	t.Parallel()

	setup := func() (*RequestLogger, *bytes.Buffer) {
		var (
			logBuf     bytes.Buffer
			middleware = NewRequestLogger()
		)
		middleware.logger = log.New(&logBuf, "", 0)
		return middleware, &logBuf
	}

	t.Run("DoesNotRedactMostHeaders", func(t *testing.T) {
		t.Parallel()

		middleware, logBuf := setup()

		const stainlessUserAgent = "Stainless"

		req := httptest.NewRequest("GET", "https://example.com", nil)
		req.Header.Set("User-Agent", stainlessUserAgent)

		var nextMiddlewareRan bool
		middleware.Middleware()(req, func(req *http.Request) (*http.Response, error) {
			nextMiddlewareRan = true

			// The request sent down through middleware shouldn't be mutated.
			require.Equal(t, stainlessUserAgent, req.Header.Get("User-Agent"))

			return &http.Response{}, nil
		})

		require.True(t, nextMiddlewareRan)
		require.Contains(t, logBuf.String(), "User-Agent: "+stainlessUserAgent)
	})

	const secretToken = "secret-token"

	t.Run("RedactsAuthorizationHeader", func(t *testing.T) {
		t.Parallel()

		middleware, logBuf := setup()

		req := httptest.NewRequest("GET", "https://example.com", nil)
		req.Header.Set("Authorization", secretToken)

		var nextMiddlewareRan bool
		middleware.Middleware()(req, func(req *http.Request) (*http.Response, error) {
			nextMiddlewareRan = true

			// The request sent down through middleware shouldn't be mutated.
			require.Equal(t, secretToken, req.Header.Get("Authorization"))

			return &http.Response{}, nil
		})

		require.True(t, nextMiddlewareRan)
		require.Contains(t, logBuf.String(), "Authorization: "+redactedPlaceholder)
	})

	t.Run("RedactsOnlySecretInAuthorizationHeader", func(t *testing.T) {
		t.Parallel()

		middleware, logBuf := setup()

		req := httptest.NewRequest("GET", "https://example.com", nil)
		req.Header.Set("Authorization", "Bearer "+secretToken)

		var nextMiddlewareRan bool
		middleware.Middleware()(req, func(req *http.Request) (*http.Response, error) {
			nextMiddlewareRan = true

			return &http.Response{}, nil
		})

		require.True(t, nextMiddlewareRan)
		require.Contains(t, logBuf.String(), "Authorization: Bearer "+redactedPlaceholder)
	})

	t.Run("RedactsMultipleAuthorizationHeaders", func(t *testing.T) {
		t.Parallel()

		middleware, logBuf := setup()

		req := httptest.NewRequest("GET", "https://example.com", nil)
		req.Header.Add("Authorization", secretToken+"1")
		req.Header.Add("Authorization", secretToken+"2")

		var nextMiddlewareRan bool
		middleware.Middleware()(req, func(req *http.Request) (*http.Response, error) {
			nextMiddlewareRan = true

			// The request sent down through middleware shouldn't be mutated.
			require.Equal(t, []string{secretToken + "1", secretToken + "2"}, req.Header.Values("Authorization"))

			return &http.Response{}, nil
		})

		require.True(t, nextMiddlewareRan)

		if strings.Count(logBuf.String(), "Authorization: "+redactedPlaceholder) != 2 {
			t.Error("expected exactly two redacted placeholders in authorization headers")
		}
	})

	const customAPIKeyHeader = "X-My-Api-Key"

	t.Run("RedactsSensitiveHeaders", func(t *testing.T) {
		t.Parallel()

		middleware, logBuf := setup()

		middleware.sensitiveHeaders = []string{customAPIKeyHeader}

		req := httptest.NewRequest("GET", "https://example.com", nil)
		req.Header.Set(customAPIKeyHeader, secretToken)

		var nextMiddlewareRan bool
		middleware.Middleware()(req, func(req *http.Request) (*http.Response, error) {
			nextMiddlewareRan = true

			// The request sent down through middleware shouldn't be mutated.
			require.Equal(t, secretToken, req.Header.Get(customAPIKeyHeader))

			return &http.Response{}, nil
		})

		require.True(t, nextMiddlewareRan)
		require.Contains(t, logBuf.String(), customAPIKeyHeader+": "+redactedPlaceholder)
	})

	t.Run("RedactsMultipleSensitiveHeaders", func(t *testing.T) {
		t.Parallel()

		middleware, logBuf := setup()

		middleware.sensitiveHeaders = []string{customAPIKeyHeader}

		req := httptest.NewRequest("GET", "https://example.com", nil)
		req.Header.Add(customAPIKeyHeader, secretToken+"1")
		req.Header.Add(customAPIKeyHeader, secretToken+"2")

		var nextMiddlewareRan bool
		middleware.Middleware()(req, func(req *http.Request) (*http.Response, error) {
			nextMiddlewareRan = true

			// The request sent down through middleware shouldn't be mutated.
			require.Equal(t, []string{secretToken + "1", secretToken + "2"}, req.Header.Values(customAPIKeyHeader))

			return &http.Response{}, nil
		})

		require.True(t, nextMiddlewareRan)
		require.Equal(t, 2, strings.Count(logBuf.String(), customAPIKeyHeader+": "+redactedPlaceholder))
	})

	t.Run("DoesNotConsumeRequestBodyWhenIoReader", func(t *testing.T) {
		t.Parallel()

		middleware, logBuf := setup()
		middleware.sensitiveHeaders = []string{customAPIKeyHeader}

		const bodyContent = "test request body content"
		bodyReader := strings.NewReader(bodyContent)

		req := httptest.NewRequest("POST", "https://example.com", bodyReader)
		req.Header.Set("Authorization", secretToken)

		var nextMiddlewareRan bool
		middleware.Middleware()(req, func(req *http.Request) (*http.Response, error) {
			nextMiddlewareRan = true

			// The request body should still be fully readable after the middleware runs
			body, err := io.ReadAll(req.Body)
			require.NoError(t, err)
			require.Equal(t, bodyContent, string(body))

			// The request sent down through middleware shouldn't be mutated.
			require.Equal(t, secretToken, req.Header.Get("Authorization"))

			return &http.Response{}, nil
		})

		require.True(t, nextMiddlewareRan)
		require.Contains(t, logBuf.String(), "Authorization: "+redactedPlaceholder)
	})
}

func TestDefaultSensitiveHeadersRedactLogsAndPreserveWire(t *testing.T) {
	t.Parallel()
	var logs bytes.Buffer
	middleware := NewRequestLogger()
	middleware.logger = log.New(&logs, "", 0)
	req := httptest.NewRequest(http.MethodPost, "https://example.invalid", strings.NewReader("request-body"))
	req.Header.Add("aUtHoRiZaTiOn", "Bearer request-authorization-secret")
	req.Header.Add("AUTHORIZATION", "Bearer request-second-authorization-secret")
	for _, name := range []string{"aPi-KeY", "X-aPi-kEy", "cOoKiE", "sEt-CoOkIe"} {
		req.Header.Add(name, "request-header-secret-"+name)
		req.Header.Add(name, "request-second-header-secret-"+name)
	}
	req.Header.Set("X-Public", "public-request")
	original := req.Header.Clone()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader("response-body")),
	}
	for _, name := range []string{"aUtHoRiZaTiOn", "aPi-KeY", "X-aPi-kEy", "cOoKiE", "sEt-CoOkIe"} {
		resp.Header.Add(name, "response-header-secret-"+name)
		resp.Header.Add(name, "response-second-header-secret-"+name)
	}
	resp.Header.Set("X-Public", "public-response")
	originalResponse := resp.Header.Clone()
	got, err := middleware.Middleware()(req, func(r *http.Request) (*http.Response, error) {
		require.Equal(t, original, r.Header)
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, "request-body", string(body))
		return resp, nil
	})
	require.NoError(t, err)
	require.Equal(t, original, req.Header)
	require.Equal(t, originalResponse, got.Header)
	body, err := io.ReadAll(got.Body)
	require.NoError(t, err)
	require.Equal(t, "response-body", string(body))
	require.NotContains(t, logs.String(), "request-authorization-secret")
	require.NotContains(t, logs.String(), "request-second-authorization-secret")
	require.NotContains(t, logs.String(), "request-header-secret-")
	require.NotContains(t, logs.String(), "request-second-header-secret-")
	require.NotContains(t, logs.String(), "response-header-secret-")
	require.NotContains(t, logs.String(), "response-second-header-secret-")
	require.Equal(t, 20, strings.Count(logs.String(), redactedPlaceholder))
	require.Equal(t, 2, strings.Count(logs.String(), "Authorization: Bearer "+redactedPlaceholder))
	for _, content := range []string{"public-request", "public-response", "request-body", "response-body"} {
		require.Contains(t, logs.String(), content)
	}
}

func TestDumpRedactedResponseRestoresHeaders(t *testing.T) {
	t.Parallel()
	middleware := NewRequestLogger()
	t.Run("dump error", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Set-Cookie": {"response-cookie"}},
			Body:       io.NopCloser(iotest.ErrReader(io.ErrUnexpectedEOF)),
		}
		original := resp.Header.Clone()
		_, err := middleware.dumpRedactedResponse(resp)
		require.ErrorIs(t, err, io.ErrUnexpectedEOF)
		require.Equal(t, original, resp.Header)
	})
	t.Run("nil headers", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("body")),
		}
		_, err := middleware.dumpRedactedResponse(resp)
		require.NoError(t, err)
		require.Nil(t, resp.Header)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "body", string(body))
	})
}

func TestDebugLogRawHeaderKeysOverHTTP(t *testing.T) {
	var logs bytes.Buffer
	sensitive := []string{"Authorization", "Api-Key", "X-Api-Key", "Cookie", "Set-Cookie"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil || string(body) != "request-body" {
			t.Errorf("wire body = %q, %v", body, err)
		}
		for _, name := range sensitive {
			if len(r.Header.Values(name)) != 4 {
				t.Errorf("wire header %s = %v", name, r.Header.Values(name))
			}
			w.Header().Add(name, "Bearer response-fake-secret-"+name)
			w.Header().Add(name, "Basic response-fake-secret-second-"+name)
		}
		w.Header().Set("X-Public", "public-response")
		_, _ = w.Write([]byte("response-body"))
	}))
	defer server.Close()
	req, err := http.NewRequest(http.MethodPost, server.URL, strings.NewReader("request-body"))
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range sensitive {
		// Raw map assignment retains case variants, unlike Header.Add/Set.
		req.Header[name] = []string{"Bearer request-fake-secret-" + name, "Basic request-fake-secret-second-" + name}
		req.Header[strings.ToLower(name)] = []string{"Bearer request-fake-secret-lower-" + name}
		req.Header[strings.ToUpper(name)] = []string{"Bearer request-fake-secret-upper-" + name}
	}
	req.Header.Set("X-Public", "public-request")
	original := req.Header.Clone()
	logger := NewRequestLogger()
	logger.logger = log.New(&logs, "", 0)
	middleware := logger.Middleware()
	var originalResponse http.Header
	got, err := middleware(req, func(r *http.Request) (*http.Response, error) {
		if !reflect.DeepEqual(r.Header, original) {
			t.Error("original request headers changed")
		}
		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			return nil, err
		}
		// A wrapper transport may supply raw lowercase/mixed-case header keys.
		for _, name := range sensitive {
			values := resp.Header.Values(name)
			resp.Header.Del(name)
			resp.Header[strings.ToLower(name)] = values
			resp.Header[strings.ToUpper(name)] = []string{"Bearer response-fake-secret-upper-" + name}
		}
		originalResponse = resp.Header.Clone()
		return resp, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	defer got.Body.Close()
	if !reflect.DeepEqual(req.Header, original) || !reflect.DeepEqual(got.Header, originalResponse) {
		t.Error("original headers changed")
	}
	body, err := io.ReadAll(got.Body)
	if err != nil || string(body) != "response-body" {
		t.Fatalf("response body = %q, %v", body, err)
	}
	if strings.Contains(logs.String(), "fake-secret") {
		t.Fatalf("sensitive headers leaked: %s", logs.String())
	}
	for _, public := range []string{"public-request", "public-response", "request-body", "response-body"} {
		if !strings.Contains(logs.String(), public) {
			t.Errorf("missing public content %q", public)
		}
	}
	if !strings.Contains(logs.String(), "Bearer "+redactedPlaceholder) || !strings.Contains(logs.String(), "Basic "+redactedPlaceholder) {
		t.Error("authorization scheme was not preserved")
	}
}
