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

var tradesList = cli.Command{
	Name:    "list",
	Usage:   "Recent Trades History",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "market",
			Usage:     "Trading pair code to query.",
			Required:  true,
			QueryPath: "market",
		},
		&requestflag.Flag[int64]{
			Name:      "count",
			Usage:     "Number of trade records to retrieve.\n\nUp to 500 supported. Default: 1.\n",
			Default:   1,
			QueryPath: "count",
		},
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Cursor for pagination.\nEnter the \"sequential_id\" from the response to retrieve the previous \"count\" trade records prior to that trade.\n",
			QueryPath: "cursor",
		},
		&requestflag.Flag[int64]{
			Name:      "days-ago",
			Usage:     "Day offset between the query date and the request date.\nMust specify the target date; up to 7 days of history is supported (UTC-based).\n\nInteger between 1 and 7. If omitted, returns trades for the current date. If 7, returns trades from 7 days ago in reverse chronological order.\n",
			QueryPath: "days_ago",
		},
		&requestflag.Flag[string]{
			Name:      "to",
			Usage:     "End time within the query date range (UTC).\nOptional parameter to retrieve trades within a specific time of the query date.\n\nAccepts HHmmss or HH:mm:ss format. Trade list is returned in reverse chronological order from the specified time.\n",
			QueryPath: "to",
		},
	},
	Action:          handleTradesList,
	HideHelpCommand: true,
}

func handleTradesList(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.TradeListParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Trades.List(ctx, params, options...)
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
		Title:          "trades list",
		Transform:      transform,
	})
}
