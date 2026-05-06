package cmd

import (
	"context"
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/upbit-official/upbit-cli/internal/apiquery"
	"github.com/upbit-official/upbit-cli/internal/requestflag"
	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
	"github.com/urfave/cli/v3"
)

var tickersListByQuoteCurrencies = cli.Command{
	Name:    "list-by-quote-currencies",
	Usage:   "List Tickers by Market",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "quote-currencies",
			Usage:     "List of quote currencies to query.\nFor multiple markets, use comma-separated format.\n\n[Example] SGD,BTC\n",
			Required:  true,
			QueryPath: "quote_currencies",
		},
	},
	Action:          handleTickersListByQuoteCurrencies,
	HideHelpCommand: true,
}

var tickersListByTradingPairs = cli.Command{
	Name:    "list-by-trading-pairs",
	Usage:   "List Tickers by Pairs",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "markets",
			Usage:     "List of trading pairs to query.\n\nFor multiple pairs, use comma-separated format.\n\n[Example] SGD-BTC,SGD-ETH\n",
			Required:  true,
			QueryPath: "markets",
		},
	},
	Action:          handleTickersListByTradingPairs,
	HideHelpCommand: true,
}

func handleTickersListByQuoteCurrencies(ctx context.Context, cmd *cli.Command) error {
	client := upbit.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	params := upbit.TickerListByQuoteCurrenciesParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Tickers.ListByQuoteCurrencies(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "tickers list-by-quote-currencies",
		Transform:      transform,
	})
}

func handleTickersListByTradingPairs(ctx context.Context, cmd *cli.Command) error {
	client := upbit.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	params := upbit.TickerListByTradingPairsParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Tickers.ListByTradingPairs(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "tickers list-by-trading-pairs",
		Transform:      transform,
	})
}
