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

var travelRuleListVasps = cli.Command{
	Name:            "list-vasps",
	Usage:           "List Travel Rule Supporting VASPs",
	Suggest:         true,
	Flags:           []cli.Flag{},
	Action:          handleTravelRuleListVasps,
	HideHelpCommand: true,
}

var travelRuleVerifyDepositByTxid = cli.Command{
	Name:    "verify-deposit-by-txid",
	Usage:   "Verify Travel Rule by Deposit TxID",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "currency",
			Usage:    "Currency code to be queried.",
			Required: true,
			BodyPath: "currency",
		},
		&requestflag.Flag[string]{
			Name:     "net-type",
			Usage:    "Blockchain network identifier for deposit/withdrawal of digital assets.\nUsed as a filter parameter to specify the network.\n",
			Required: true,
			BodyPath: "net_type",
		},
		&requestflag.Flag[string]{
			Name:     "txid",
			Usage:    "Transaction ID of the deposit to be verified.",
			Required: true,
			BodyPath: "txid",
		},
		&requestflag.Flag[string]{
			Name:     "vasp-uuid",
			Usage:    "Unique identifier (UUID) of the counterparty exchange from which the asset was withdrawn.",
			Required: true,
			BodyPath: "vasp_uuid",
		},
	},
	Action:          handleTravelRuleVerifyDepositByTxid,
	HideHelpCommand: true,
}

var travelRuleVerifyDepositByUuid = cli.Command{
	Name:    "verify-deposit-by-uuid",
	Usage:   "Verify Travel Rule by Deposit UUID",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "deposit-uuid",
			Usage:    "Unique identifier (UUID) for the deposit to verify.",
			Required: true,
			BodyPath: "deposit_uuid",
		},
		&requestflag.Flag[string]{
			Name:     "vasp-uuid",
			Usage:    "Unique identifier (UUID) of the counterparty exchange from which the asset was withdrawn.",
			Required: true,
			BodyPath: "vasp_uuid",
		},
	},
	Action:          handleTravelRuleVerifyDepositByUuid,
	HideHelpCommand: true,
}

func handleTravelRuleListVasps(ctx context.Context, cmd *cli.Command) error {
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
	_, err = client.TravelRule.ListVasps(ctx, options...)
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
		Title:          "travel-rule list-vasps",
		Transform:      transform,
	})
}

func handleTravelRuleVerifyDepositByTxid(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.TravelRuleVerifyDepositByTxidParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.TravelRule.VerifyDepositByTxid(ctx, params, options...)
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
		Title:          "travel-rule verify-deposit-by-txid",
		Transform:      transform,
	})
}

func handleTravelRuleVerifyDepositByUuid(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.TravelRuleVerifyDepositByUuidParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.TravelRule.VerifyDepositByUuid(ctx, params, options...)
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
		Title:          "travel-rule verify-deposit-by-uuid",
		Transform:      transform,
	})
}
