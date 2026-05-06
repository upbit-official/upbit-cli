package cmd

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/upbit-official/upbit-cli/internal/autocomplete"
	"github.com/upbit-official/upbit-cli/internal/requestflag"
	docs "github.com/urfave/cli-docs/v3"
	"github.com/urfave/cli/v3"
)

var (
	Command            *cli.Command
	CommandErrorBuffer bytes.Buffer
)

func init() {
	Command = &cli.Command{
		Name:      "upbit",
		Usage:     "CLI for the upbit API",
		Suggest:   true,
		Version:   Version,
		ErrWriter: &CommandErrorBuffer,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable debug logging",
			},
			&cli.StringFlag{
				Name:        "base-url",
				DefaultText: "environment default URL",
				Usage:       "Override the base URL for API requests",
				Validator: func(baseURL string) error {
					return ValidateBaseURL(baseURL, "--base-url")
				},
			},
			&cli.StringFlag{
				Name:  "format",
				Usage: "The format for displaying response data (one of: " + strings.Join(OutputFormats, ", ") + ")",
				Value: "auto",
				Validator: func(format string) error {
					if !slices.Contains(OutputFormats, strings.ToLower(format)) {
						return fmt.Errorf("format must be one of: %s", strings.Join(OutputFormats, ", "))
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:  "format-error",
				Usage: "The format for displaying error data (one of: " + strings.Join(OutputFormats, ", ") + ")",
				Value: "auto",
				Validator: func(format string) error {
					if !slices.Contains(OutputFormats, strings.ToLower(format)) {
						return fmt.Errorf("format must be one of: %s", strings.Join(OutputFormats, ", "))
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:  "transform",
				Usage: "The GJSON transformation for data output.",
			},
			&cli.StringFlag{
				Name:  "transform-error",
				Usage: "The GJSON transformation for errors.",
			},
			&cli.BoolFlag{
				Name:    "raw-output",
				Aliases: []string{"r"},
				Usage:   "If the result is a string, print it without JSON quotes. This can be useful for making output transforms talk to non-JSON-based systems.",
			},
			&requestflag.Flag[string]{
				Name:    "access-key",
				Usage:   "The access key provided by Upbit for API authentication.\nFor more details, please refer to https://docs.upbit.com/reference/auth.\n",
				Sources: cli.EnvVars("UPBIT_ACCESS_KEY"),
			},
			&requestflag.Flag[string]{
				Name:    "secret-key",
				Usage:   "The secret key used to sign API requests for secure verification.\nFor more details, please refer to https://docs.upbit.com/reference/auth.\n",
				Sources: cli.EnvVars("UPBIT_SECRET_KEY"),
			},
			&cli.StringFlag{
				Name:  "environment",
				Usage: "Set the environment for API requests (kr, sg, id, th)",
				Validator: func(environment string) error {
					if environment == "" {
						return nil
					}

					switch strings.ToLower(environment) {
					case "kr", "sg", "id", "th":
						return nil
					default:
						return fmt.Errorf("environment must be one of: kr, sg, id, th")
					}
				},
			},
		},
		Commands: []*cli.Command{
			{
				Name:     "accounts",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&accountsList,
				},
			},
			{
				Name:     "travel-rule",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&travelRuleListVasps,
					&travelRuleVerifyDepositByTxid,
					&travelRuleVerifyDepositByUuid,
				},
			},
			{
				Name:     "wallet-status",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&walletStatusList,
				},
			},
			{
				Name:     "api-keys",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&apiKeysList,
				},
			},
			{
				Name:     "orders",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&ordersCreate,
					&ordersRetrieve,
					&ordersCancel,
					&ordersCancelAndNew,
					&ordersCancelByUuids,
					&ordersCancelOpen,
					&ordersListByUuids,
					&ordersListClosed,
					&ordersListOpen,
					&ordersRetrieveChance,
					&ordersTestCreate,
				},
			},
			{
				Name:     "withdraws",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&withdrawsRetrieve,
					&withdrawsList,
					&withdrawsCancelWithdrawal,
					&withdrawsCreateKrwWithdrawal,
					&withdrawsCreateWithdrawal,
					&withdrawsListCoinAddresses,
					&withdrawsRetrieveChance,
				},
			},
			{
				Name:     "deposits",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&depositsRetrieve,
					&depositsList,
					&depositsCreateCoinAddress,
					&depositsDepositKrw,
					&depositsListCoinAddresses,
					&depositsRetrieveChance,
					&depositsRetrieveCoinAddress,
				},
			},
			{
				Name:     "trading-pairs",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&tradingPairsList,
				},
			},
			{
				Name:     "tickers",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&tickersListByQuoteCurrencies,
					&tickersListByTradingPairs,
				},
			},
			{
				Name:     "orderbooks",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&orderbooksList,
					&orderbooksListInstruments,
				},
			},
			{
				Name:     "trades",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&tradesList,
				},
			},
			{
				Name:     "candles",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&candlesListDays,
					&candlesListMinutes,
					&candlesListMonths,
					&candlesListSeconds,
					&candlesListWeeks,
					&candlesListYears,
				},
			},
			{
				Name:            "@manpages",
				Usage:           "Generate documentation for 'man'",
				UsageText:       "upbit @manpages [-o upbit.1] [--gzip]",
				Hidden:          true,
				Action:          generateManpages,
				HideHelpCommand: true,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "write manpages to the given folder",
						Value:   "man",
					},
					&cli.BoolFlag{
						Name:    "gzip",
						Aliases: []string{"z"},
						Usage:   "output gzipped manpage files to .gz",
						Value:   true,
					},
					&cli.BoolFlag{
						Name:    "text",
						Aliases: []string{"t"},
						Usage:   "output uncompressed text files",
						Value:   false,
					},
				},
			},
			{
				Name:            "__complete",
				Hidden:          true,
				HideHelpCommand: true,
				Action:          autocomplete.ExecuteShellCompletion,
			},
			{
				Name:            "@completion",
				Hidden:          true,
				HideHelpCommand: true,
				Action:          autocomplete.OutputCompletionScript,
			},
		},
		HideHelpCommand: true,
	}
}

func generateManpages(ctx context.Context, c *cli.Command) error {
	manpage, err := docs.ToManWithSection(Command, 1)
	if err != nil {
		return err
	}
	dir := c.String("output")
	err = os.MkdirAll(filepath.Join(dir, "man1"), 0755)
	if err != nil {
		// handle error
	}
	if c.Bool("text") {
		file, err := os.Create(filepath.Join(dir, "man1", "upbit.1"))
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := file.WriteString(manpage); err != nil {
			return err
		}
	}
	if c.Bool("gzip") {
		file, err := os.Create(filepath.Join(dir, "man1", "upbit.1.gz"))
		if err != nil {
			return err
		}
		defer file.Close()
		gzWriter := gzip.NewWriter(file)
		defer gzWriter.Close()
		_, err = gzWriter.Write([]byte(manpage))
		if err != nil {
			return err
		}
	}
	fmt.Printf("Wrote manpages to %s\n", dir)
	return nil
}
