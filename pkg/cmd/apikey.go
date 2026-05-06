package cmd

import (
	"context"
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/upbit-official/upbit-cli/internal/apiquery"
	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
	"github.com/urfave/cli/v3"
)

var apiKeysList = cli.Command{
	Name:            "list",
	Usage:           "List API Keys",
	Suggest:         true,
	Flags:           []cli.Flag{},
	Action:          handleAPIKeysList,
	HideHelpCommand: true,
}

func handleAPIKeysList(ctx context.Context, cmd *cli.Command) error {
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
	_, err = client.APIKeys.List(ctx, options...)
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
		Title:          "api-keys list",
		Transform:      transform,
	})
}
