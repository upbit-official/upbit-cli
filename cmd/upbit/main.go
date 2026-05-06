package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"slices"

	"github.com/tidwall/gjson"
	"github.com/upbit-official/upbit-cli/pkg/cmd"
	"github.com/upbit-official/upbit-sdk-go"
	"github.com/urfave/cli/v3"
)

func main() {
	app := cmd.Command

	if slices.Contains(os.Args, "__complete") {
		prepareForAutocomplete(app)
	}

	if baseURL, ok := os.LookupEnv("UPBIT_BASE_URL"); ok {
		if err := cmd.ValidateBaseURL(baseURL, "UPBIT_BASE_URL"); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		exitCode := 1

		// Check if error has a custom exit code
		if exitErr, ok := err.(cli.ExitCoder); ok {
			exitCode = exitErr.ExitCode()
		}

		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			fmt.Fprintf(os.Stderr, "%s %q: %d %s\n", apierr.Request.Method, apierr.Request.URL, apierr.Response.StatusCode, http.StatusText(apierr.Response.StatusCode))
			format := app.String("format-error")
			json := gjson.Parse(apierr.RawJSON())
			show_err := cmd.ShowJSON(json, cmd.ShowJSONOpts{
				ExplicitFormat: app.IsSet("format-error"),
				Format:         format,
				Title:          "Error",
				Transform:      app.String("transform-error"),
			})
			if show_err != nil {
				// Just print the original error:
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
			if apierr.Response.StatusCode == 401 && cmd.CredentialsMissing {
				fmt.Fprintf(os.Stderr, "\nNo credentials configured. Run 'upbit config set' to set up your API keys, or set UPBIT_ACCESS_KEY / UPBIT_SECRET_KEY environment variables.\n")
			}
		} else {
			if cmd.CommandErrorBuffer.Len() > 0 {
				os.Stderr.Write(cmd.CommandErrorBuffer.Bytes())
			} else {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
		}
		os.Exit(exitCode)
	}
}

func prepareForAutocomplete(cmd *cli.Command) {
	// urfave/cli does not handle flag completions and will print an error if we inspect a command with invalid flags.
	// This skips that sort of validation
	cmd.SkipFlagParsing = true
	for _, child := range cmd.Commands {
		prepareForAutocomplete(child)
	}
}
