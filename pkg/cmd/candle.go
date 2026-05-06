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

var candlesListDays = cli.Command{
	Name:    "list-days",
	Usage:   "List Day Candles",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "market",
			Usage:     "Trading pair code to query.",
			Required:  true,
			QueryPath: "market",
		},
		&requestflag.Flag[string]{
			Name:      "converting-price-unit",
			Usage:     "Currency to convert the closing price into.\nWhen specified, the response includes a \"converted_trade_price\" field with the closing price converted to the given currency.\n\n[Example] Specifying \"KRW\" returns the closing price converted to KRW.\n",
			QueryPath: "converting_price_unit",
		},
		&requestflag.Flag[int64]{
			Name:      "count",
			Usage:     "Number of candles to retrieve.\n\nUp to 200 candles are supported. Default: 1.\n",
			Default:   1,
			QueryPath: "count",
		},
		&requestflag.Flag[string]{
			Name:      "to",
			Usage:     "End time of the query range.\nRetrieves candles before the specified time. If not specified, the most recent candles are returned.\n\nAccepts ISO 8601 datetime format.\n",
			QueryPath: "to",
		},
	},
	Action:          handleCandlesListDays,
	HideHelpCommand: true,
}

var candlesListMinutes = cli.Command{
	Name:    "list-minutes",
	Usage:   "List Minute Candles",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[int64]{
			Name:      "unit",
			Usage:     "Candle unit in minutes.\n\nSpecify the candle unit to retrieve. Up to 240 minutes (4 hours) is supported.\n",
			Required:  true,
			PathParam: "unit",
		},
		&requestflag.Flag[string]{
			Name:      "market",
			Usage:     "Trading pair code to query.",
			Required:  true,
			QueryPath: "market",
		},
		&requestflag.Flag[int64]{
			Name:      "count",
			Usage:     "Number of candles to retrieve.\n\nUp to 200 candles are supported. Default: 1.\n",
			Default:   1,
			QueryPath: "count",
		},
		&requestflag.Flag[string]{
			Name:      "to",
			Usage:     "End time of the query range.\nRetrieves candles before the specified time. If not specified, the most recent candles are returned.\n\nAccepts ISO 8601 datetime format.\n\n[Example]\n2025-06-24T04:56:53Z\n2025-06-24 04:56:53\n2025-06-24T13:56:53+08:00\n",
			QueryPath: "to",
		},
	},
	Action:          handleCandlesListMinutes,
	HideHelpCommand: true,
}

var candlesListMonths = cli.Command{
	Name:    "list-months",
	Usage:   "List Month Candles",
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
			Usage:     "Number of candles to retrieve. Up to 200, default: 1.",
			Default:   1,
			QueryPath: "count",
		},
		&requestflag.Flag[string]{
			Name:      "to",
			Usage:     "End time of the query range. If not specified, the most recent candles are returned.\n\nAccepts ISO 8601 datetime format.\n",
			QueryPath: "to",
		},
	},
	Action:          handleCandlesListMonths,
	HideHelpCommand: true,
}

var candlesListSeconds = cli.Command{
	Name:    "list-seconds",
	Usage:   "List Second Candles",
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
			Usage:     "Number of candles to retrieve.\n\nUp to 200 candles are supported. Default: 1.\n",
			Default:   1,
			QueryPath: "count",
		},
		&requestflag.Flag[string]{
			Name:      "to",
			Usage:     "End time of the query range.\nRetrieves candles before the specified time. If not specified, the most recent candles are returned.\n\nAccepts ISO 8601 datetime format. URL encoding is required for spaces and special characters.\n\n[Example] 2025-06-24T04:56:53Z\n2025-06-24 04:56:53\n2025-06-24T13:56:53+08:00\n\nSecond candles only support data up to 3 months prior. An empty array is returned for older times.\n",
			QueryPath: "to",
		},
	},
	Action:          handleCandlesListSeconds,
	HideHelpCommand: true,
}

var candlesListWeeks = cli.Command{
	Name:    "list-weeks",
	Usage:   "List Week Candles",
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
			Usage:     "Number of candles to retrieve. Up to 200, default: 1.\n",
			Default:   1,
			QueryPath: "count",
		},
		&requestflag.Flag[string]{
			Name:      "to",
			Usage:     "End time of the query range.\nRetrieves candles before the specified time. If not specified, the most recent candles are returned.\n",
			QueryPath: "to",
		},
	},
	Action:          handleCandlesListWeeks,
	HideHelpCommand: true,
}

var candlesListYears = cli.Command{
	Name:    "list-years",
	Usage:   "List Year Candles",
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
			Usage:     "Number of candles to retrieve. Up to 200, default: 1.",
			Default:   1,
			QueryPath: "count",
		},
		&requestflag.Flag[string]{
			Name:      "to",
			Usage:     "End time of the query range. ISO 8601 format.",
			QueryPath: "to",
		},
	},
	Action:          handleCandlesListYears,
	HideHelpCommand: true,
}

func handleCandlesListDays(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.CandleListDaysParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Candles.ListDays(ctx, params, options...)
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
		Title:          "candles list-days",
		Transform:      transform,
	})
}

func handleCandlesListMinutes(ctx context.Context, cmd *cli.Command) error {
	client := upbit.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("unit") && len(unusedArgs) > 0 {
		cmd.Set("unit", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
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

	params := upbit.CandleListMinutesParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Candles.ListMinutes(
		ctx,
		cmd.Value("unit").(int64),
		params,
		options...,
	)
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
		Title:          "candles list-minutes",
		Transform:      transform,
	})
}

func handleCandlesListMonths(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.CandleListMonthsParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Candles.ListMonths(ctx, params, options...)
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
		Title:          "candles list-months",
		Transform:      transform,
	})
}

func handleCandlesListSeconds(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.CandleListSecondsParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Candles.ListSeconds(ctx, params, options...)
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
		Title:          "candles list-seconds",
		Transform:      transform,
	})
}

func handleCandlesListWeeks(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.CandleListWeeksParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Candles.ListWeeks(ctx, params, options...)
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
		Title:          "candles list-weeks",
		Transform:      transform,
	})
}

func handleCandlesListYears(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.CandleListYearsParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Candles.ListYears(ctx, params, options...)
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
		Title:          "candles list-years",
		Transform:      transform,
	})
}
