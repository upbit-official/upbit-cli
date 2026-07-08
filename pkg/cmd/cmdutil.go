package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/upbit-official/upbit-cli/internal/jsonview"
	"github.com/upbit-official/upbit-cli/pkg/config"
	"github.com/upbit-official/upbit-sdk-go/option"

	"github.com/charmbracelet/x/term"
	"github.com/itchyny/json2yaml"
	"github.com/muesli/reflow/wrap"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"github.com/urfave/cli/v3"
)

var OutputFormats = []string{"auto", "explore", "json", "jsonl", "pretty", "raw", "yaml"}

// CredentialsMissing is set to true by getDefaultRequestOptions when both
// access-key and secret-key are absent from all sources (flags, env, file).
// main.go reads this to print a user-friendly setup guide.
var CredentialsMissing bool

var (
	credentialsOnce    sync.Once
	cachedFileCreds    config.Credentials
	cachedFileCredsErr error
)

func loadFileCreds() {
	cachedFileCreds, cachedFileCredsErr = config.LoadFile()
}

// ValidateBaseURL checks that a base URL is correctly prefixed with a protocol scheme and produces a better
// error message than the person would see otherwise if it doesn't.
func ValidateBaseURL(value, source string) error {
	if value != "" && !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
		return fmt.Errorf("%s %q is missing a scheme (expected http:// or https://)", source, value)
	}
	return nil
}

func getDefaultRequestOptions(cmd *cli.Command) []option.RequestOption {
	opts := []option.RequestOption{
		option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			req.Header.Set("User-Agent", fmt.Sprintf("Upbit/CLI %s", Version))
			req.Header.Set("X-Upbit-Client-Type", "cli")
			req.Header.Set("X-Upbit-Client-Name", "upbit-cli")
			req.Header.Set("X-Upbit-Client-Version", Version)
			return next(req)
		}),
	}

	accessSet := cmd.IsSet("access-key")
	secretSet := cmd.IsSet("secret-key")

	if accessSet {
		opts = append(opts, option.WithAccessKey(cmd.String("access-key")))
	}
	if secretSet {
		opts = append(opts, option.WithSecretKey(cmd.String("secret-key")))
	}

	// Fall back to config file when env/flag not set.
	if !accessSet || !secretSet {
		credentialsOnce.Do(loadFileCreds)
		if cachedFileCredsErr != nil {
			fmt.Fprintf(os.Stderr, "warning: could not load credentials from config file: %v\n", cachedFileCredsErr)
		} else {
			if !accessSet && cachedFileCreds.AccessKey != "" {
				opts = append(opts, option.WithAccessKey(cachedFileCreds.AccessKey))
				accessSet = true
			}
			if !secretSet && cachedFileCreds.SecretKey != "" {
				opts = append(opts, option.WithSecretKey(cachedFileCreds.SecretKey))
				secretSet = true
			}
		}
	}

	if !accessSet || !secretSet {
		CredentialsMissing = true
	}

	// Override base URL if the --base-url flag is provided
	if baseURL := cmd.String("base-url"); baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}

	// Set environment if the --environment flag is provided
	if environment := cmd.String("environment"); environment != "" {
		switch strings.ToLower(environment) {
		case "kr":
			opts = append(opts, option.WithEnvironmentKr())
		case "sg":
			opts = append(opts, option.WithEnvironmentSg())
		case "id":
			opts = append(opts, option.WithEnvironmentId())
		case "th":
			opts = append(opts, option.WithEnvironmentTh())
		default:
			// The environment flag validator in cmd.go rejects invalid values.
		}
	}

	// Apply any custom headers passed via --header / -H. Header options must be applied
	// as a middleware (not via option.WithHeader directly): WithHeader touches the request
	// immediately, but at client-construction time there is no request yet, so it would
	// panic. The middleware runs per request, once the *http.Request exists.
	if raw := cmd.StringSlice("header"); len(raw) > 0 {
		if parsed, err := parseHeaderFlags(raw); err == nil {
			if len(parsed) > 0 {
				opts = append(opts, headerMiddlewareOption(parsed))
			}
		} else {
			// The --header flag validator in cmd.go rejects malformed values, so this
			// branch should be unreachable in practice; warn rather than fail silently.
			fmt.Fprintf(os.Stderr, "warning: could not apply --header values: %v\n", err)
		}
	}

	return opts
}

// reservedHeaderKeys are the exact header keys the CLI reserves for itself. They must not be
// overridable via --header: a user-supplied value for any of these keys is silently ignored so
// the CLI's own value always wins. Keys are stored in canonical form (http.CanonicalHeaderKey)
// for case-insensitive matching.
//
// Authorization is included because the SDK signs each request with a JWT derived from the
// credentials; letting --header replace it would break authentication. Any header in the
// X-Upbit-Client-* namespace is also reserved — see isReservedHeader / reservedHeaderPrefix.
var reservedHeaderKeys = map[string]bool{
	http.CanonicalHeaderKey("User-Agent"):    true,
	http.CanonicalHeaderKey("Authorization"): true,
}

// reservedHeaderPrefix reserves the X-Upbit-Client-* header namespace for the CLI, so
// client-identity headers (X-Upbit-Client-Type, -Name, -Version) and any future X-Upbit-Client-*
// header cannot be overridden via --header. Other X-Upbit-* headers (e.g. X-Upbit-Initiator) are
// deliberately NOT reserved and remain user-settable. Stored in canonical form for
// case-insensitive matching.
var reservedHeaderPrefix = http.CanonicalHeaderKey("X-Upbit-Client-")

// isReservedHeader reports whether canonical (already run through http.CanonicalHeaderKey) is a
// header the CLI owns and thus must not be overridable via --header.
func isReservedHeader(canonical string) bool {
	return reservedHeaderKeys[canonical] || strings.HasPrefix(canonical, reservedHeaderPrefix)
}

// headerFlag is a single parsed --header value. add is false for the first occurrence of a
// given (canonical) key — applied with Header.Set, overwriting any prior value — and true for
// subsequent occurrences of the same key, applied with Header.Add so repeated --header flags
// build a multi-value header.
type headerFlag struct {
	key   string
	value string
	add   bool
}

// parseHeaderFlags parses raw "Key: Value" header flag values. The first occurrence of a given
// canonical key is marked for Set; later occurrences of the same key are marked for Add.
//
// Values for reserved keys (see reservedHeaderKeys) are silently skipped so the CLI's own
// client-identity and auth headers cannot be overridden by --header.
func parseHeaderFlags(headers []string) ([]headerFlag, error) {
	var parsed []headerFlag
	seen := make(map[string]bool)
	for _, h := range headers {
		key, value, found := strings.Cut(h, ":")
		if !found {
			return nil, fmt.Errorf("invalid header %q: expected format 'Key: Value'", h)
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" {
			return nil, fmt.Errorf("invalid header %q: header key must not be empty", h)
		}
		canonical := http.CanonicalHeaderKey(key)
		// Reserved headers are owned by the CLI and cannot be overridden; skip silently.
		if isReservedHeader(canonical) {
			continue
		}
		parsed = append(parsed, headerFlag{key: key, value: value, add: seen[canonical]})
		seen[canonical] = true
	}
	return parsed, nil
}

// headerMiddlewareOption returns a request option that applies the parsed headers to each
// outgoing request. Reserved headers (client-identity and auth) are already filtered out by
// parseHeaderFlags, so nothing here can override a CLI-owned header.
func headerMiddlewareOption(parsed []headerFlag) option.RequestOption {
	return option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		for _, hf := range parsed {
			if hf.add {
				req.Header.Add(hf.key, hf.value)
			} else {
				req.Header.Set(hf.key, hf.value)
			}
		}
		return next(req)
	})
}

var debugMiddlewareOption = option.WithMiddleware(
	func(r *http.Request, mn option.MiddlewareNext) (*http.Response, error) {
		logger := log.Default()

		if reqBytes, err := httputil.DumpRequest(r, true); err == nil {
			logger.Printf("Request Content:\n%s\n", reqBytes)
		}

		resp, err := mn(r)
		if err != nil {
			return resp, err
		}

		if respBytes, err := httputil.DumpResponse(resp, true); err == nil {
			logger.Printf("Response Content:\n%s\n", respBytes)
		}

		return resp, err
	},
)

// isInputPiped tries to check for input being piped into the CLI which tells us that we should try to read
// from stdin. This can be a bit tricky in some cases like when an stdin is connected to a pipe but nothing is
// being piped in (this may happen in some environments like Cursor's integration terminal or CI), which is
// why this function is a little more elaborate than it'd be otherwise.
func isInputPiped() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	mode := stat.Mode()

	// Regular file (redirect like < file.txt) — only if non-empty.
	//
	// Notably, on Unix the case like `< /dev/null` is handled below because `/dev/null` is not a regular
	// file. On Windows, NUL appears as a regular file with size 0, so it's also handled correctly.
	if mode.IsRegular() && stat.Size() > 0 {
		return true
	}

	// For pipes/sockets (e.g. `echo foo | stainlesscli`), use an OS-specific check to determine whether
	// data is actually available. Some environments like Cursor's integrated terminal connect stdin as a
	// pipe even when nothing is being piped.
	if mode&(os.ModeNamedPipe|os.ModeSocket) != 0 {
		// Defined in either cmdutil_unix.go or cmdutil_windows.go.
		return isPipedDataAvailableOSSpecific()
	}

	return false
}

func isTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return term.IsTerminal(v.Fd())
	default:
		return false
	}
}

func streamOutput(label string, generateOutput func(w *os.File) error) error {
	// For non-tty output (probably a pipe), write directly to stdout
	if !isTerminal(os.Stdout) {
		return streamToStdout(generateOutput)
	}

	// When streaming output on Unix-like systems, there's a special trick involving creating two socket pairs
	// that we prefer because it supports small buffer sizes which results in less pagination per buffer. The
	// constructs needed to run it don't exist on Windows builds, so we have this function broken up into
	// OS-specific files with conditional build comments. Under Windows (and in case our fancy constructs fail
	// on Unix), we fall back to using pipes (`streamToPagerWithPipe`), which are OS agnostic.
	//
	// Defined in either cmdutil_unix.go or cmdutil_windows.go.
	return streamOutputOSSpecific(label, generateOutput)
}

func streamToPagerWithPipe(label string, generateOutput func(w *os.File) error) error {
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	defer r.Close()
	defer w.Close()

	pagerProgram := os.Getenv("PAGER")
	if pagerProgram == "" {
		pagerProgram = "less"
	}

	if _, err := exec.LookPath(pagerProgram); err != nil {
		return err
	}

	cmd := exec.Command(pagerProgram)
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		"LESS=-X -r -P "+label,
		"MORE=-r -P "+label,
	)

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := r.Close(); err != nil {
		return err
	}

	// If we would be streaming to a terminal and aren't forcing color one way
	// or the other, we should configure things to use color so the pager gets
	// colorized input.
	if isTerminal(os.Stdout) && os.Getenv("FORCE_COLOR") == "" {
		os.Setenv("FORCE_COLOR", "1")
	}

	if err := generateOutput(w); err != nil && !strings.Contains(err.Error(), "broken pipe") {
		return err
	}

	w.Close()
	return cmd.Wait()
}

func streamToStdout(generateOutput func(w *os.File) error) error {
	signal.Ignore(syscall.SIGPIPE)
	err := generateOutput(os.Stdout)
	if err != nil && strings.Contains(err.Error(), "broken pipe") {
		return nil
	}
	return err
}

// writeBinaryResponse writes a binary response to stdout or a file.
//
// Takes in a stdout reference so we can test this function without overriding os.Stdout in tests.
func writeBinaryResponse(response *http.Response, stdout io.Writer, outfile string) (string, error) {
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	switch outfile {
	case "-", "/dev/stdout":
		_, err := stdout.Write(body)
		return "", err
	case "":
		// If output file is unspecified, then print to stdout for plain text or
		// if stdout is not a terminal:
		if !isTerminal(os.Stdout) || isUTF8TextFile(body) {
			_, err := stdout.Write(body)
			return "", err
		}

		// If response has a suggested filename in the content-disposition
		// header, then use that (with an optional suffix to ensure uniqueness):
		file, err := createDownloadFile(response, body)
		if err != nil {
			return "", err
		}
		defer file.Close()
		if _, err := file.Write(body); err != nil {
			return "", err
		}
		return fmt.Sprintf("Wrote output to: %s", file.Name()), nil
	default:
		if err := os.WriteFile(outfile, body, 0644); err != nil {
			return "", err
		}
		return fmt.Sprintf("Wrote output to: %s", outfile), nil
	}
}

// Return a writable file handle to a new file, which attempts to choose a good filename
// based on the Content-Disposition header or sniffing the MIME filetype of the response.
func createDownloadFile(response *http.Response, data []byte) (*os.File, error) {
	filename := "file"
	// If the header provided an output filename, use that
	disp := response.Header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(disp)
	if err == nil {
		if dispFilename, ok := params["filename"]; ok {
			// Only use the last path component to prevent directory traversal
			filename = filepath.Base(dispFilename)
			// Try to create the file with exclusive flag to avoid race conditions
			file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
			if err == nil {
				return file, nil
			}
		}
	}

	// If file already exists, create a unique filename using CreateTemp
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = guessExtension(data)
	}
	base := strings.TrimSuffix(filename, ext)
	return os.CreateTemp(".", base+"-*"+ext)
}

func guessExtension(data []byte) string {
	ct := http.DetectContentType(data)

	// Prefer common extensions over obscure ones
	switch ct {
	case "application/gzip":
		return ".gz"
	case "application/pdf":
		return ".pdf"
	case "application/zip":
		return ".zip"
	case "audio/mpeg":
		return ".mp3"
	case "image/bmp":
		return ".bmp"
	case "image/gif":
		return ".gif"
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "video/mp4":
		return ".mp4"
	}

	exts, err := mime.ExtensionsByType(ct)
	if err == nil && len(exts) > 0 {
		return exts[0]
	} else if isUTF8TextFile(data) {
		return ".txt"
	} else {
		return ".bin"
	}
}

func shouldUseColors(w io.Writer) bool {
	force, ok := os.LookupEnv("FORCE_COLOR")
	if ok {
		if force == "1" {
			return true
		}
		if force == "0" {
			return false
		}
	}
	return isTerminal(w)
}

func formatJSON(res gjson.Result, opts ShowJSONOpts) ([]byte, error) {
	if opts.Transform != "" {
		transformed := res.Get(opts.Transform)
		if transformed.Exists() {
			res = transformed
		}
	}
	// Modeled after `jq -r` (`--raw-output`): if the result is a string, print it without JSON quotes so that
	// it's easier to pipe into other programs.
	if opts.RawOutput && res.Type == gjson.String {
		return []byte(res.Str + "\n"), nil
	}
	switch strings.ToLower(opts.Format) {
	case "auto":
		autoOpts := opts
		autoOpts.Format = "json"
		autoOpts.Transform = ""
		return formatJSON(res, autoOpts)
	case "pretty":
		return []byte(jsonview.RenderJSON(opts.Title, res) + "\n"), nil
	case "json":
		prettyJSON := pretty.Pretty([]byte(res.Raw))
		if shouldUseColors(opts.Stdout) {
			return pretty.Color(prettyJSON, pretty.TerminalStyle), nil
		} else {
			return prettyJSON, nil
		}
	case "jsonl":
		// @ugly is gjson syntax for "no whitespace", so it fits on one line
		oneLineJSON := res.Get("@ugly").Raw
		if shouldUseColors(opts.Stdout) {
			bytes := append(pretty.Color([]byte(oneLineJSON), pretty.TerminalStyle), '\n')
			return bytes, nil
		} else {
			return []byte(oneLineJSON + "\n"), nil
		}
	case "raw":
		return []byte(res.Raw + "\n"), nil
	case "yaml":
		input := strings.NewReader(res.Raw)
		var yaml strings.Builder
		if err := json2yaml.Convert(&yaml, input); err != nil {
			return nil, err
		}
		_, err := opts.Stdout.Write([]byte(yaml.String()))
		return nil, err
	default:
		return nil, fmt.Errorf("Invalid format: %s, valid formats are: %s", opts.Format, strings.Join(OutputFormats, ", "))
	}
}

const warningExploreNotSupported = "Warning: Output format 'explore' not supported for non-terminal output; falling back to 'json'\n"

// ShowJSONOpts configures how JSON output is displayed.
type ShowJSONOpts struct {
	ExplicitFormat bool      // true if the user explicitly passed --format
	Format         string    // output format (auto, explore, json, jsonl, pretty, raw, yaml)
	RawOutput      bool      // like jq -r: print strings without JSON quotes
	Stderr         io.Writer // stderr for warnings; injectable for testing; defaults to os.Stderr
	Stdout         *os.File  // stdout (or pager); injectable for testing; defaults to os.Stdout
	Title          string    // display title
	Transform      string    // GJSON path to extract before displaying
}

func (o *ShowJSONOpts) setDefaults() {
	if o.Stderr == nil {
		o.Stderr = os.Stderr
	}
	if o.Stdout == nil {
		o.Stdout = os.Stdout
	}
}

// ShowJSON displays a single JSON result to the user.
func ShowJSON(res gjson.Result, opts ShowJSONOpts) error {
	opts.setDefaults()

	switch strings.ToLower(opts.Format) {
	case "auto":
		autoOpts := opts
		autoOpts.Format = "json"
		return ShowJSON(res, autoOpts)
	case "explore":
		if !isTerminal(opts.Stdout) {
			if opts.ExplicitFormat {
				fmt.Fprint(opts.Stderr, warningExploreNotSupported)
			}
			jsonOpts := opts
			jsonOpts.Format = "json"
			return ShowJSON(res, jsonOpts)
		}
		if opts.Transform != "" {
			transformed := res.Get(opts.Transform)
			if transformed.Exists() {
				res = transformed
			}
		}
		return jsonview.ExploreJSON(opts.Title, res)
	default:
		bytes, err := formatJSON(res, opts)
		if err != nil {
			return err
		}

		_, err = opts.Stdout.Write(bytes)
		return err
	}
}

// Get the number of lines that would be output by writing the data to the terminal
func countTerminalLines(data []byte, terminalWidth int) int {
	return bytes.Count([]byte(wrap.String(string(data), terminalWidth)), []byte("\n"))
}

type hasRawJSON interface {
	RawJSON() string
}

// ShowJSONIterator displays an iterator of values to the user. Use itemsToDisplay = -1 for no limit.
func ShowJSONIterator[T any](iter jsonview.Iterator[T], itemsToDisplay int64, opts ShowJSONOpts) error {
	opts.setDefaults()

	if opts.Format == "explore" {
		if isTerminal(opts.Stdout) {
			return jsonview.ExploreJSONStream(opts.Title, iter)
		}
		if opts.ExplicitFormat {
			fmt.Fprint(opts.Stderr, warningExploreNotSupported)
		}
		opts.Format = "json"
	}

	terminalWidth, terminalHeight, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		terminalWidth = 100
		terminalHeight = 40
	}

	// Decide whether or not to use a pager based on whether it's a short output or a long output
	usePager := false
	output := []byte{}
	numberOfNewlines := 0
	// -1 is used to signal no limit of items to display
	for itemsToDisplay != 0 && iter.Next() {
		item := iter.Current()
		var obj gjson.Result
		if hasRaw, ok := any(item).(hasRawJSON); ok {
			obj = gjson.Parse(hasRaw.RawJSON())
		} else {
			jsonData, err := json.Marshal(item)
			if err != nil {
				return err
			}
			obj = gjson.ParseBytes(jsonData)
		}
		json, err := formatJSON(obj, opts)
		if err != nil {
			return err
		}

		output = append(output, json...)
		itemsToDisplay -= 1
		numberOfNewlines += countTerminalLines(json, terminalWidth)

		// If the output won't fit in the terminal window, stream it to a pager
		if numberOfNewlines >= terminalHeight-3 {
			usePager = true
			break
		}
	}

	if !usePager {
		_, err := opts.Stdout.Write(output)
		if err != nil {
			return err
		}

		return iter.Err()
	}

	return streamOutput(opts.Title, func(pager *os.File) error {
		_, err := pager.Write(output)
		if err != nil {
			return err
		}

		pagerOpts := opts
		pagerOpts.Stdout = pager

		for iter.Next() {
			if itemsToDisplay == 0 {
				break
			}
			item := iter.Current()
			var obj gjson.Result
			if hasRaw, ok := any(item).(hasRawJSON); ok {
				obj = gjson.Parse(hasRaw.RawJSON())
			} else {
				jsonData, err := json.Marshal(item)
				if err != nil {
					return err
				}
				obj = gjson.ParseBytes(jsonData)
			}
			if err := ShowJSON(obj, pagerOpts); err != nil {
				return err
			}
			itemsToDisplay -= 1
		}
		return iter.Err()
	})
}
