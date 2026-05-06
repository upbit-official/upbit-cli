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

var depositsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get Deposit",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "currency",
			Usage:     "Currency code filter.",
			QueryPath: "currency",
		},
		&requestflag.Flag[string]{
			Name:      "txid",
			Usage:     "Transaction ID of the deposit to query.\n\nIf neither uuid nor txid is provided, the latest deposit is returned.\n",
			QueryPath: "txid",
		},
		&requestflag.Flag[string]{
			Name:      "uuid",
			Usage:     "UUID of the deposit to query.\n\nIf neither uuid nor txid is provided, the latest deposit is returned.\n",
			QueryPath: "uuid",
		},
	},
	Action:          handleDepositsRetrieve,
	HideHelpCommand: true,
}

var depositsList = cli.Command{
	Name:    "list",
	Usage:   "List Deposits",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "currency",
			Usage:     "Currency code filter.\nIf not specified, the latest deposits are returned.\n",
			QueryPath: "currency",
		},
		&requestflag.Flag[string]{
			Name:      "from",
			Usage:     "Cursor for pagination.\nEnter a \"uuid\" from the response to retrieve \"limit\" deposits after that timestamp.\n",
			QueryPath: "from",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Number of results per request (default: 100, max: 100).",
			Default:   100,
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "order-by",
			Usage:     `Sort order. "desc" (newest first) or "asc" (oldest first). Default is "desc".`,
			Default:   "desc",
			QueryPath: "order_by",
		},
		&requestflag.Flag[int64]{
			Name:      "page",
			Usage:     "Page number for pagination. Default is 1.",
			Default:   1,
			QueryPath: "page",
		},
		&requestflag.Flag[string]{
			Name:      "state",
			Usage:     "Deposit state filter.\n\n* `PROCESSING`: Processing\n* `ACCEPTED`: Completed\n* `CANCELLED`: Cancelled\n* `REJECTED`: Rejected\n* `TRAVEL_RULE_SUSPECTED`: Pending Travel Rule verification\n* `REFUNDING`: Refund in progress\n* `REFUNDED`: Refunded\n",
			QueryPath: "state",
		},
		&requestflag.Flag[string]{
			Name:      "to",
			Usage:     "Cursor for pagination.\nEnter a \"uuid\" from the response to retrieve \"limit\" deposits before that timestamp.\n",
			QueryPath: "to",
		},
		&requestflag.Flag[[]string]{
			Name:      "txid",
			Usage:     "List of transaction IDs to query.\n\n[Example] txids[]=txid1&txids[]=txid2\n",
			QueryPath: "txids",
		},
		&requestflag.Flag[[]string]{
			Name:      "uuid",
			Usage:     "List of UUIDs to query.\n\n[Example] uuids[]=uuid1&uuids[]=uuid2\n",
			QueryPath: "uuids",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleDepositsList,
	HideHelpCommand: true,
}

var depositsCreateCoinAddress = cli.Command{
	Name:    "create-coin-address",
	Usage:   "Create Deposit Address",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "currency",
			Usage:    "Currency code for which to create a deposit address.",
			Required: true,
			BodyPath: "currency",
		},
		&requestflag.Flag[string]{
			Name:     "net-type",
			Usage:    "Network type.",
			Required: true,
			BodyPath: "net_type",
		},
	},
	Action:          handleDepositsCreateCoinAddress,
	HideHelpCommand: true,
}

var depositsDepositKrw = cli.Command{
	Name:    "deposit-krw",
	Usage:   "Deposit KRW",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "amount",
			Usage:    "KRW amount to deposit.",
			Required: true,
			BodyPath: "amount",
		},
		&requestflag.Flag[string]{
			Name:     "two-factor-type",
			Usage:    "Two-factor authentication method for KRW transactions.\n\n* `kakao`: Kakao authentication\n* `naver`: Naver authentication\n* `hana`: Hana certificate authentication\n",
			Required: true,
			BodyPath: "two_factor_type",
		},
	},
	Action:          handleDepositsDepositKrw,
	HideHelpCommand: true,
}

var depositsListCoinAddresses = cli.Command{
	Name:            "list-coin-addresses",
	Usage:           "List Deposit Addresses",
	Suggest:         true,
	Flags:           []cli.Flag{},
	Action:          handleDepositsListCoinAddresses,
	HideHelpCommand: true,
}

var depositsRetrieveChance = cli.Command{
	Name:    "retrieve-chance",
	Usage:   "Get Available Deposit Information",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "currency",
			Usage:     "Currency code to query deposit availability for.",
			Required:  true,
			QueryPath: "currency",
		},
		&requestflag.Flag[string]{
			Name:      "net-type",
			Usage:     "Blockchain network identifier for digital asset deposit.\n\nUsed to filter by network.\n",
			Required:  true,
			QueryPath: "net_type",
		},
	},
	Action:          handleDepositsRetrieveChance,
	HideHelpCommand: true,
}

var depositsRetrieveCoinAddress = cli.Command{
	Name:    "retrieve-coin-address",
	Usage:   "Get Deposit Address",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "currency",
			Usage:     "Currency code to query.",
			Required:  true,
			QueryPath: "currency",
		},
		&requestflag.Flag[string]{
			Name:      "net-type",
			Usage:     "Blockchain network identifier.\n\nUsed to filter by network.\n",
			Required:  true,
			QueryPath: "net_type",
		},
	},
	Action:          handleDepositsRetrieveCoinAddress,
	HideHelpCommand: true,
}

func handleDepositsRetrieve(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.DepositGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Deposits.Get(ctx, params, options...)
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
		Title:          "deposits retrieve",
		Transform:      transform,
	})
}

func handleDepositsList(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.DepositListParams{}

	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Deposits.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "deposits list",
			Transform:      transform,
		})
	} else {
		iter := client.Deposits.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "deposits list",
			Transform:      transform,
		})
	}
}

func handleDepositsCreateCoinAddress(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.DepositNewCoinAddressParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Deposits.NewCoinAddress(ctx, params, options...)
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
		Title:          "deposits create-coin-address",
		Transform:      transform,
	})
}

func handleDepositsDepositKrw(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.DepositDepositKrwParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Deposits.DepositKrw(ctx, params, options...)
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
		Title:          "deposits deposit-krw",
		Transform:      transform,
	})
}

func handleDepositsListCoinAddresses(ctx context.Context, cmd *cli.Command) error {
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
	_, err = client.Deposits.ListCoinAddresses(ctx, options...)
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
		Title:          "deposits list-coin-addresses",
		Transform:      transform,
	})
}

func handleDepositsRetrieveChance(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.DepositGetChanceParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Deposits.GetChance(ctx, params, options...)
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
		Title:          "deposits retrieve-chance",
		Transform:      transform,
	})
}

func handleDepositsRetrieveCoinAddress(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.DepositGetCoinAddressParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Deposits.GetCoinAddress(ctx, params, options...)
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
		Title:          "deposits retrieve-coin-address",
		Transform:      transform,
	})
}
