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

var orderbooksList = cli.Command{
	Name:    "list",
	Usage:   "Get Orderbook",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "markets",
			Usage:     "List of trading pairs to query.\n\nFor multiple pairs, use comma-separated format.\n",
			Required:  true,
			QueryPath: "markets",
		},
		&requestflag.Flag[int64]{
			Name:      "count",
			Usage:     "Number of orderbook entries to retrieve.\n\nBased on the best bid-ask pair, returns the specified number of pairs. Default: 30.\n",
			Default:   30,
			QueryPath: "count",
		},
		&requestflag.Flag[string]{
			Name:      "level",
			Usage:     "Orderbook aggregation level. Only supported for KRW markets.\nGroups ask/bid price and size by the specified unit. Provide as a numeric string.\nUse an integer string for units >= 1, or a double string for fractional units. Defaults to 0 if not specified.\n",
			Default:   "0",
			QueryPath: "level",
		},
	},
	Action:          handleOrderbooksList,
	HideHelpCommand: true,
}

var orderbooksListInstruments = cli.Command{
	Name:    "list-instruments",
	Usage:   "List Orderbook Instruments",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "markets",
			Usage:     "List of trading pairs to query.\n\nFor multiple pairs, use comma-separated format.\n",
			Required:  true,
			QueryPath: "markets",
		},
	},
	Action:          handleOrderbooksListInstruments,
	HideHelpCommand: true,
}

func handleOrderbooksList(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.OrderbookListParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Orderbooks.List(ctx, params, options...)
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
		Title:          "orderbooks list",
		Transform:      transform,
	})
}

func handleOrderbooksListInstruments(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.OrderbookListInstrumentsParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Orderbooks.ListInstruments(ctx, params, options...)
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
		Title:          "orderbooks list-instruments",
		Transform:      transform,
	})
}
