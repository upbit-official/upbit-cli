package debugmiddleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"reflect"
	"strings"
)

// For the time being these type definitions are duplicated here so that we can
// test this file in a non-generated context.
type (
	Middleware     = func(*http.Request, MiddlewareNext) (*http.Response, error)
	MiddlewareNext = func(*http.Request) (*http.Response, error)
)

const redactedPlaceholder = "<REDACTED>"

// Headers known to contain sensitive information like an API key. Note that this exclude `Authorization`,
// which is handled specially in `redactHeaders` below.
var sensitiveHeaders = []string{"api-key", "x-api-key", "cookie", "set-cookie"}

// RequestLogger is a middleware that logs HTTP requests and responses.
type RequestLogger struct {
	logger           interface{ Printf(string, ...any) } // field for testability; usually log.Default()
	sensitiveHeaders []string                            // field for testability; usually sensitiveHeaders
}

// NewRequestLogger returns a new RequestLogger instance with default options.
func NewRequestLogger() *RequestLogger {
	return &RequestLogger{
		logger:           log.Default(),
		sensitiveHeaders: sensitiveHeaders,
	}
}

func (m *RequestLogger) Middleware() Middleware {
	return func(req *http.Request, mn MiddlewareNext) (*http.Response, error) {
		redacted, err := m.redactRequest(req)
		if err != nil {
			return nil, err
		}
		if reqBytes, err := httputil.DumpRequest(redacted, true); err == nil {
			m.logger.Printf("Request Content:\n%s\n", reqBytes)
		}

		resp, err := mn(req)
		if err != nil {
			return resp, err
		}

		if respBytes, err := m.dumpRedactedResponse(resp); err == nil {
			m.logger.Printf("Response Content:\n%s\n", respBytes)
		}

		return resp, err
	}
}

// redactRequest redacts sensitive information from the request for logging
// purposes. If redaction is necessary, the request is cloned before mutating
// the original and that clone is returned. As a small optimization, the
// original is request is returned unchanged if no redaction is necessary.
func (m *RequestLogger) redactRequest(req *http.Request) (*http.Request, error) {
	redactedHeaders := m.redactHeaders(req.Header)

	if reflect.DeepEqual(req.Header, redactedHeaders) {
		return req, nil
	}

	redacted := req.Clone(req.Context())
	redacted.Header = redactedHeaders
	var err error
	redacted.Body, req.Body, err = cloneBody(req.Body)
	return redacted, err
}

// DumpResponse restores the response body after reading it. Use the original
// response so that callers keep that readable body, and restore its headers even
// when dumping fails or panics.
func (m *RequestLogger) dumpRedactedResponse(resp *http.Response) ([]byte, error) {
	original := resp.Header
	resp.Header = m.redactHeaders(original)
	defer func() { resp.Header = original }()
	return httputil.DumpResponse(resp, true)
}

func (m *RequestLogger) redactHeaders(headers http.Header) http.Header {
	redactedHeaders := headers.Clone()

	// Inspect actual map keys: Header.Values only finds canonicalized keys,
	// while custom transports can supply lowercase or mixed-case duplicates.
	for key, values := range redactedHeaders {
		if strings.EqualFold(key, "Authorization") {
			for i, value := range values {
				// In case we're using something like a bearer token (e.g. `Bearer
				// <my_token>`), keep the `Bearer` part for more debugging
				// information.
				if authKind, _, ok := strings.Cut(value, " "); ok {
					values[i] = authKind + " " + redactedPlaceholder
				} else {
					values[i] = redactedPlaceholder
				}
			}
			continue
		}
		for _, header := range m.sensitiveHeaders {
			if !strings.EqualFold(key, header) {
				continue
			}
			for i := range values {
				values[i] = redactedPlaceholder
			}
			break
		}
	}

	return redactedHeaders
}

// This function returns two copies of an HTTP request body that can each be
// read independently without affecting the other.
// This logic is taken from `drainBody` in net/http/httputil.
func cloneBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == nil || b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return io.NopCloser(&buf), io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
