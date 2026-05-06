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

var withdrawsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get Withdrawal",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "currency",
			Usage:     "Currency code filter.\nIf not specified, the latest withdrawal is returned.\n",
			QueryPath: "currency",
		},
		&requestflag.Flag[string]{
			Name:      "txid",
			Usage:     "Transaction ID of the withdrawal to query.\n\nIf neither uuid nor txid is provided, the latest withdrawal is returned.\n",
			QueryPath: "txid",
		},
		&requestflag.Flag[string]{
			Name:      "uuid",
			Usage:     "UUID of the withdrawal to query.\n\nIf neither uuid nor txid is provided, the latest withdrawal is returned.\n",
			QueryPath: "uuid",
		},
	},
	Action:          handleWithdrawsRetrieve,
	HideHelpCommand: true,
}

var withdrawsList = cli.Command{
	Name:    "list",
	Usage:   "List Withdrawals",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "currency",
			Usage:     "Currency code filter.\nIf not specified, all currencies are returned.\n",
			QueryPath: "currency",
		},
		&requestflag.Flag[string]{
			Name:      "from",
			Usage:     "Cursor for pagination.\nEnter a \"uuid\" from the response to retrieve \"limit\" withdrawals after that timestamp.\n",
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
			Usage:     "Withdrawal state filter.\n\n* `WAITING`: Waiting\n* `PROCESSING`: Processing\n* `DONE`: Completed\n* `FAILED`: Failed\n* `CANCELLED`: Cancelled\n* `REJECTED`: Rejected\n",
			QueryPath: "state",
		},
		&requestflag.Flag[string]{
			Name:      "to",
			Usage:     "Cursor for pagination.\nEnter a \"uuid\" from the response to retrieve \"limit\" withdrawals before that timestamp.\n",
			QueryPath: "to",
		},
		&requestflag.Flag[[]string]{
			Name:      "txid",
			Usage:     "List of transaction IDs to query. Maximum 100. Cannot be used together with uuids.\n\n[Example] txids[]=txid1&txids[]=txid2\n",
			QueryPath: "txids",
		},
		&requestflag.Flag[[]string]{
			Name:      "uuid",
			Usage:     "List of UUIDs to query. Maximum 100. Cannot be used together with txids.\n\n[Example] uuids[]=uuid1&uuids[]=uuid2\n",
			QueryPath: "uuids",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleWithdrawsList,
	HideHelpCommand: true,
}

var withdrawsCancelWithdrawal = cli.Command{
	Name:    "cancel-withdrawal",
	Usage:   "Cancel Withdrawal",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "uuid",
			Usage:     "UUID of the withdrawal to cancel.",
			Required:  true,
			QueryPath: "uuid",
		},
	},
	Action:          handleWithdrawsCancelWithdrawal,
	HideHelpCommand: true,
}

var withdrawsCreateKrwWithdrawal = cli.Command{
	Name:    "create-krw-withdrawal",
	Usage:   "Withdraw KRW",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "amount",
			Usage:    "KRW amount to withdraw.\n\nMust be a numeric string.\n",
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
	Action:          handleWithdrawsCreateKrwWithdrawal,
	HideHelpCommand: true,
}

var withdrawsCreateWithdrawal = cli.Command{
	Name:    "create-withdrawal",
	Usage:   "Withdraw Digital Asset",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "address",
			Usage:    "Recipient address for digital asset withdrawal.\n\nOnly addresses registered in the allowed withdrawal address list can be used.\n",
			Required: true,
			BodyPath: "address",
		},
		&requestflag.Flag[string]{
			Name:     "amount",
			Usage:    "Amount of the asset to withdraw.\n\nMust be a numeric string.\n",
			Required: true,
			BodyPath: "amount",
		},
		&requestflag.Flag[string]{
			Name:     "currency",
			Usage:    "Currency code of the digital asset to withdraw.",
			Required: true,
			BodyPath: "currency",
		},
		&requestflag.Flag[string]{
			Name:     "net-type",
			Usage:    "Check the allowed withdrawal addresses API response to find the available \"net_type\" value for each address.\n",
			Required: true,
			BodyPath: "net_type",
		},
		&requestflag.Flag[*string]{
			Name:     "secondary-address",
			Usage:    "Secondary withdrawal address (Destination Tag, Memo, or Message).\nIf the recipient address includes a secondary address, this field must be included.\n",
			BodyPath: "secondary_address",
		},
		&requestflag.Flag[string]{
			Name:     "transaction-type",
			Usage:    "Withdrawal transaction type.\n\n* `default`: Standard withdrawal\n* `internal`: Internal (instant) withdrawal\n",
			Default:  "default",
			BodyPath: "transaction_type",
		},
	},
	Action:          handleWithdrawsCreateWithdrawal,
	HideHelpCommand: true,
}

var withdrawsListCoinAddresses = cli.Command{
	Name:            "list-coin-addresses",
	Usage:           "List Withdrawal Allowed Addresses",
	Suggest:         true,
	Flags:           []cli.Flag{},
	Action:          handleWithdrawsListCoinAddresses,
	HideHelpCommand: true,
}

var withdrawsRetrieveChance = cli.Command{
	Name:    "retrieve-chance",
	Usage:   "Get Available Withdrawal Information",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "currency",
			Usage:     "Currency code to query withdrawal availability for.",
			Required:  true,
			QueryPath: "currency",
		},
		&requestflag.Flag[string]{
			Name:      "net-type",
			Usage:     "Blockchain network identifier for digital asset withdrawal.\nRequired for digital assets.\n",
			Required:  true,
			QueryPath: "net_type",
		},
	},
	Action:          handleWithdrawsRetrieveChance,
	HideHelpCommand: true,
}

func handleWithdrawsRetrieve(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.WithdrawGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Withdraws.Get(ctx, params, options...)
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
		Title:          "withdraws retrieve",
		Transform:      transform,
	})
}

func handleWithdrawsList(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.WithdrawListParams{}

	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Withdraws.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "withdraws list",
			Transform:      transform,
		})
	} else {
		iter := client.Withdraws.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "withdraws list",
			Transform:      transform,
		})
	}
}

func handleWithdrawsCancelWithdrawal(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.WithdrawCancelWithdrawalParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Withdraws.CancelWithdrawal(ctx, params, options...)
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
		Title:          "withdraws cancel-withdrawal",
		Transform:      transform,
	})
}

func handleWithdrawsCreateKrwWithdrawal(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.WithdrawNewKrwWithdrawalParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Withdraws.NewKrwWithdrawal(ctx, params, options...)
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
		Title:          "withdraws create-krw-withdrawal",
		Transform:      transform,
	})
}

func handleWithdrawsCreateWithdrawal(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.WithdrawNewWithdrawalParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Withdraws.NewWithdrawal(ctx, params, options...)
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
		Title:          "withdraws create-withdrawal",
		Transform:      transform,
	})
}

func handleWithdrawsListCoinAddresses(ctx context.Context, cmd *cli.Command) error {
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
	_, err = client.Withdraws.ListCoinAddresses(ctx, options...)
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
		Title:          "withdraws list-coin-addresses",
		Transform:      transform,
	})
}

func handleWithdrawsRetrieveChance(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.WithdrawGetChanceParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Withdraws.GetChance(ctx, params, options...)
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
		Title:          "withdraws retrieve-chance",
		Transform:      transform,
	})
}
