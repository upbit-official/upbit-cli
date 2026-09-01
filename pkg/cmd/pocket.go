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

var pocketsList = cli.Command{
	Name:            "list",
	Usage:           "List Pockets",
	Suggest:         true,
	Flags:           []cli.Flag{},
	Action:          handlePocketsList,
	HideHelpCommand: true,
}

var pocketsListAPIKeys = cli.Command{
	Name:    "list-api-keys",
	Usage:   "List Pocket API Keys",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[bool]{
			Name:      "include-expired",
			Usage:     "Whether to include expired API keys in the results.",
			Default:   false,
			QueryPath: "include_expired",
		},
		&requestflag.Flag[[]string]{
			Name:      "uuid",
			Usage:     "Array of pocket UUIDs to query.\n\n[Example] uuids[]=uuid1&uuids[]=uuid2\n",
			QueryPath: "uuids",
		},
	},
	Action:          handlePocketsListAPIKeys,
	HideHelpCommand: true,
}

var pocketsListTransfers = cli.Command{
	Name:    "list-transfers",
	Usage:   "List Transfers",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "currency",
			Usage:     "Asset code to query.",
			QueryPath: "currency",
		},
		&requestflag.Flag[string]{
			Name:      "direction",
			Usage:     "Transfer direction.\n\n* `in`: Transfers into the sub-pocket\n* `out`: Transfers out of the sub-pocket\n* `all`: Transfers in all directions\n",
			QueryPath: "direction",
		},
		&requestflag.Flag[string]{
			Name:      "end-time",
			Usage:     "End time of the query period.\nRetrieves records created from `start_time` to the specified time.\n\nDefault: Current time.\nLimit: Maximum query range is 7 days.\nAllowed format: ISO 8601 with timezone or millisecond timestamp.\n",
			QueryPath: "end_time",
		},
		&requestflag.Flag[[]string]{
			Name:      "identifier",
			Usage:     "List of client-provided identifiers specified when requesting transfers.\n\n[Example] identifiers[]=id1&identifiers[]=id2\n\nLimit: Up to 20 items.\n",
			QueryPath: "identifiers",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Number of results to request (default 20, max 100).",
			Default:   20,
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "order-by",
			Usage:     "Sort order. `asc` for oldest first, `desc` for newest first.",
			Default:   "desc",
			QueryPath: "order_by",
		},
		&requestflag.Flag[string]{
			Name:      "start-time",
			Usage:     "Start time of the query period.\nRetrieves records created from the specified time to `end_time`.\n\nAllowed format: ISO 8601 with timezone or millisecond timestamp.\n",
			QueryPath: "start_time",
		},
		&requestflag.Flag[[]string]{
			Name:      "state",
			Usage:     "Asset transfer request states.\nReturns only transfers in the specified states. If omitted, returns transfers in all states.\n\n[Example] states[]=done&states[]=failed\n",
			QueryPath: "states",
		},
		&requestflag.Flag[[]string]{
			Name:      "uuid",
			Usage:     "List of UUIDs identifying asset transfer requests.\nReturns only transfers matching the specified UUIDs.\n\n[Example] uuids[]=uuid1&uuids[]=uuid2\n\nLimit: Up to 20 items.\n",
			QueryPath: "uuids",
		},
	},
	Action:          handlePocketsListTransfers,
	HideHelpCommand: true,
}

var pocketsListUniversalTransfers = cli.Command{
	Name:    "list-universal-transfers",
	Usage:   "List Universal Transfers",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "currency",
			Usage:     "Asset code to query.",
			QueryPath: "currency",
		},
		&requestflag.Flag[string]{
			Name:      "end-time",
			Usage:     "End time of the query period.\nRetrieves records created from `start_time` to the specified time.\n\nDefault: Current time.\nLimit: Maximum query range is 7 days.\nAllowed format: ISO 8601 with timezone or millisecond timestamp.\n",
			QueryPath: "end_time",
		},
		&requestflag.Flag[string]{
			Name:      "from",
			Usage:     "UUID of the pocket that sent the asset.",
			QueryPath: "from",
		},
		&requestflag.Flag[[]string]{
			Name:      "identifier",
			Usage:     "List of client-provided identifiers specified when requesting transfers.\n\n[Example] identifiers[]=id1&identifiers[]=id2\n\nLimit: Up to 20 items.\n",
			QueryPath: "identifiers",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Number of results to request (default 20, max 100).",
			Default:   20,
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "order-by",
			Usage:     "Sort order. `asc` for oldest first, `desc` for newest first.",
			Default:   "desc",
			QueryPath: "order_by",
		},
		&requestflag.Flag[string]{
			Name:      "start-time",
			Usage:     "Start time of the query period.\nRetrieves records created from the specified time to `end_time`.\n\nAllowed format: ISO 8601 with timezone or millisecond timestamp.\n",
			QueryPath: "start_time",
		},
		&requestflag.Flag[[]string]{
			Name:      "state",
			Usage:     "Asset transfer request states. Available values: `submitted`, `processing`, `done`, `failed`.\nReturns only transfers in the specified states. If omitted, returns transfers in all states.\n\n[Example] states[]=done&states[]=failed\n",
			QueryPath: "states",
		},
		&requestflag.Flag[string]{
			Name:      "to",
			Usage:     "UUID of the pocket that received the asset.",
			QueryPath: "to",
		},
		&requestflag.Flag[[]string]{
			Name:      "uuid",
			Usage:     "List of UUIDs identifying asset transfer requests.\nReturns only transfers matching the specified UUIDs.\n\n[Example] uuids[]=uuid1&uuids[]=uuid2\n\nLimit: Up to 20 items.\n",
			QueryPath: "uuids",
		},
	},
	Action:          handlePocketsListUniversalTransfers,
	HideHelpCommand: true,
}

var pocketsRetrieveBalance = cli.Command{
	Name:    "retrieve-balance",
	Usage:   "Get Sub-Pocket Balance",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "uuid",
			Usage:     "UUID of the sub-pocket to query.",
			Required:  true,
			QueryPath: "uuid",
		},
	},
	Action:          handlePocketsRetrieveBalance,
	HideHelpCommand: true,
}

var pocketsTransfer = cli.Command{
	Name:    "transfer",
	Usage:   "Transfer",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "amount",
			Usage:    "Asset quantity to transfer.",
			Required: true,
			BodyPath: "amount",
		},
		&requestflag.Flag[string]{
			Name:     "currency",
			Usage:    "Asset code to transfer.\n\n[Example] KRW, BTC, ETH\n",
			Required: true,
			BodyPath: "currency",
		},
		&requestflag.Flag[string]{
			Name:     "to",
			Usage:    "UUID of the receiving pocket.\n",
			Required: true,
			BodyPath: "to",
		},
		&requestflag.Flag[string]{
			Name:     "identifier",
			Usage:    "Client-provided identifier for the asset transfer request.\nReturns null if not specified.\n\nAllowed format: 1 to 64 letters, numbers, underscores (_), periods (.), or hyphens (-).\n",
			BodyPath: "identifier",
		},
	},
	Action:          handlePocketsTransfer,
	HideHelpCommand: true,
}

var pocketsUniversalTransfer = cli.Command{
	Name:    "universal-transfer",
	Usage:   "Universal Transfer",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "amount",
			Usage:    "Asset quantity to transfer.",
			Required: true,
			BodyPath: "amount",
		},
		&requestflag.Flag[string]{
			Name:     "currency",
			Usage:    "Asset code to transfer.\n\n[Example] KRW, BTC, ETH\n",
			Required: true,
			BodyPath: "currency",
		},
		&requestflag.Flag[string]{
			Name:     "to",
			Usage:    "UUID of the receiving pocket.",
			Required: true,
			BodyPath: "to",
		},
		&requestflag.Flag[string]{
			Name:     "from",
			Usage:    "UUID of the sending pocket.\nIf omitted, defaults to the pocket associated with the API key used for the request.\n",
			BodyPath: "from",
		},
		&requestflag.Flag[string]{
			Name:     "identifier",
			Usage:    "Client-provided identifier for the asset transfer request.\nReturns null if not specified.\n\nAllowed format: 1 to 64 letters, numbers, underscores (_), periods (.), or hyphens (-).\n",
			BodyPath: "identifier",
		},
	},
	Action:          handlePocketsUniversalTransfer,
	HideHelpCommand: true,
}

func handlePocketsList(ctx context.Context, cmd *cli.Command) error {
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

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Pockets.List(ctx, options...)
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
		Title:          "pockets list",
		Transform:      transform,
	})
}

func handlePocketsListAPIKeys(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.PocketListAPIKeysParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Pockets.ListAPIKeys(ctx, params, options...)
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
		Title:          "pockets list-api-keys",
		Transform:      transform,
	})
}

func handlePocketsListTransfers(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.PocketListTransfersParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Pockets.ListTransfers(ctx, params, options...)
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
		Title:          "pockets list-transfers",
		Transform:      transform,
	})
}

func handlePocketsListUniversalTransfers(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.PocketListUniversalTransfersParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Pockets.ListUniversalTransfers(ctx, params, options...)
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
		Title:          "pockets list-universal-transfers",
		Transform:      transform,
	})
}

func handlePocketsRetrieveBalance(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.PocketGetBalanceParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Pockets.GetBalance(ctx, params, options...)
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
		Title:          "pockets retrieve-balance",
		Transform:      transform,
	})
}

func handlePocketsTransfer(ctx context.Context, cmd *cli.Command) error {
	client := upbit.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		ApplicationJSON,
		false,
	)
	if err != nil {
		return err
	}

	params := upbit.PocketTransferParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Pockets.Transfer(ctx, params, options...)
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
		Title:          "pockets transfer",
		Transform:      transform,
	})
}

func handlePocketsUniversalTransfer(ctx context.Context, cmd *cli.Command) error {
	client := upbit.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatBrackets,
		ApplicationJSON,
		false,
	)
	if err != nil {
		return err
	}

	params := upbit.PocketUniversalTransferParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Pockets.UniversalTransfer(ctx, params, options...)
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
		Title:          "pockets universal-transfer",
		Transform:      transform,
	})
}
