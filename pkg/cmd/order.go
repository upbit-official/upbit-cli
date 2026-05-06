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

var ordersCreate = cli.Command{
	Name:    "create",
	Usage:   "Create Order",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "market",
			Usage:    "Target trading pair for the order. (required)",
			Required: true,
			BodyPath: "market",
		},
		&requestflag.Flag[string]{
			Name:     "ord-type",
			Usage:    "Order type.\n\n* `limit`: Limit order\n* `price`: Market buy order (by total price)\n* `market`: Market sell order (by volume)\n* `best`: Best available price order\n",
			Required: true,
			BodyPath: "ord_type",
		},
		&requestflag.Flag[string]{
			Name:     "side",
			Usage:    "Order direction (ask=sell, bid=buy)",
			Required: true,
			BodyPath: "side",
		},
		&requestflag.Flag[string]{
			Name:     "identifier",
			Usage:    "Client-assigned order identifier.",
			BodyPath: "identifier",
		},
		&requestflag.Flag[string]{
			Name:     "price",
			Usage:    "Order price or total amount.",
			BodyPath: "price",
		},
		&requestflag.Flag[string]{
			Name:     "smp-type",
			Usage:    "Self-Match Prevention (SMP) mode.\n\n* `cancel_maker`: Cancel maker order\n* `cancel_taker`: Cancel taker order\n* `reduce`: Reduce order quantity\n",
			BodyPath: "smp_type",
		},
		&requestflag.Flag[string]{
			Name:     "time-in-force",
			Usage:    "Time in force condition.\n\n* `fok`: Fill or Kill\n* `ioc`: Immediate or Cancel\n* `post_only`: Post only (maker only)\n",
			BodyPath: "time_in_force",
		},
		&requestflag.Flag[string]{
			Name:     "volume",
			Usage:    "Order volume.",
			BodyPath: "volume",
		},
	},
	Action:          handleOrdersCreate,
	HideHelpCommand: true,
}

var ordersRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get Order",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "identifier",
			Usage:     "Client-assigned identifier of the order to query.\n\nUsed when querying by the identifier assigned at order creation.\n",
			QueryPath: "identifier",
		},
		&requestflag.Flag[string]{
			Name:      "uuid",
			Usage:     "UUID of the order to query.",
			QueryPath: "uuid",
		},
	},
	Action:          handleOrdersRetrieve,
	HideHelpCommand: true,
}

var ordersCancel = cli.Command{
	Name:    "cancel",
	Usage:   "Cancel Order",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "identifier",
			Usage:     "Client-assigned identifier of the order to cancel.\n\nUsed when cancelling by the identifier assigned at order creation.\n",
			QueryPath: "identifier",
		},
		&requestflag.Flag[string]{
			Name:      "uuid",
			Usage:     "UUID of the order to cancel.",
			QueryPath: "uuid",
		},
	},
	Action:          handleOrdersCancel,
	HideHelpCommand: true,
}

var ordersCancelAndNew = cli.Command{
	Name:    "cancel-and-new",
	Usage:   "Cancel and New Order",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "new-ord-type",
			Usage:    "Order type.\n\n* `limit`: Limit order\n* `price`: Market buy order (by total price)\n* `market`: Market sell order (by volume)\n* `best`: Best available price order\n",
			Required: true,
			BodyPath: "new_ord_type",
		},
		&requestflag.Flag[string]{
			Name:     "new-identifier",
			Usage:    "Client-assigned identifier for the new order.",
			BodyPath: "new_identifier",
		},
		&requestflag.Flag[string]{
			Name:     "new-price",
			Usage:    "Price or total amount for the new order.",
			BodyPath: "new_price",
		},
		&requestflag.Flag[string]{
			Name:     "new-smp-type",
			Usage:    "Self-Match Prevention (SMP) mode.\n\n* `cancel_maker`: Cancel maker order\n* `cancel_taker`: Cancel taker order\n* `reduce`: Reduce order quantity\n",
			BodyPath: "new_smp_type",
		},
		&requestflag.Flag[string]{
			Name:     "new-time-in-force",
			Usage:    "Time in force condition.\n\n* `fok`: Fill or Kill\n* `ioc`: Immediate or Cancel\n* `post_only`: Post only (maker only)\n",
			BodyPath: "new_time_in_force",
		},
		&requestflag.Flag[string]{
			Name:     "new-volume",
			Usage:    "Volume for the new order.",
			BodyPath: "new_volume",
		},
		&requestflag.Flag[string]{
			Name:     "prev-order-identifier",
			Usage:    "Client-assigned identifier of the order to cancel.",
			BodyPath: "prev_order_identifier",
		},
		&requestflag.Flag[string]{
			Name:     "prev-order-uuid",
			Usage:    "UUID of the order to cancel.",
			BodyPath: "prev_order_uuid",
		},
	},
	Action:          handleOrdersCancelAndNew,
	HideHelpCommand: true,
}

var ordersCancelByUuids = cli.Command{
	Name:    "cancel-by-uuids",
	Usage:   "Cancel Orders by IDs",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[[]string]{
			Name:      "identifier",
			Usage:     "List of client-assigned identifiers of orders to cancel. Maximum 20 orders.\n\n[Example] identifiers[]=id1&identifiers[]=id2…\n",
			QueryPath: "identifiers",
		},
		&requestflag.Flag[[]string]{
			Name:      "uuid",
			Usage:     "List of UUIDs of orders to cancel. Maximum 20 orders.\n\n[Example] uuids[]=uuid1&uuids[]=uuid2…\n",
			QueryPath: "uuids",
		},
	},
	Action:          handleOrdersCancelByUuids,
	HideHelpCommand: true,
}

var ordersCancelOpen = cli.Command{
	Name:    "cancel-open",
	Usage:   "Batch Cancel Orders",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "cancel-side",
			Usage:     "Side filter. \"all\" (both), \"ask\" (sell only), \"bid\" (buy only).\n",
			Default:   "all",
			QueryPath: "cancel_side",
		},
		&requestflag.Flag[int64]{
			Name:      "count",
			Usage:     "Maximum number of orders to cancel. Max 300, default 20.\n",
			Default:   20,
			QueryPath: "count",
		},
		&requestflag.Flag[string]{
			Name:      "excluded-pairs",
			Usage:     "Trading pair exclusion filter. Cancels all open orders except those for the specified pairs. Up to 20 pairs, comma-separated.\n\n[Example] excluded_pairs=KRW-BTC,KRW-ETH\n",
			QueryPath: "excluded_pairs",
		},
		&requestflag.Flag[string]{
			Name:      "order-by",
			Usage:     `Sort order for determining which orders to cancel. "desc" (newest first) or "asc" (oldest first). Default is "desc".`,
			Default:   "desc",
			QueryPath: "order_by",
		},
		&requestflag.Flag[string]{
			Name:      "pairs",
			Usage:     "Trading pair filter. Cancels open orders only for the specified pairs. Up to 20 pairs, comma-separated.\n\n[Example] pairs=KRW-BTC,KRW-ETH\n",
			QueryPath: "pairs",
		},
		&requestflag.Flag[string]{
			Name:      "quote-currencies",
			Usage:     "Quote currency filter (KRW, BTC, USDT). Cancels all open orders in markets with the specified quote currency.\n\n[Example] \"KRW\" cancels all open orders in the KRW market.\n",
			QueryPath: "quote_currencies",
		},
	},
	Action:          handleOrdersCancelOpen,
	HideHelpCommand: true,
}

var ordersListByUuids = cli.Command{
	Name:    "list-by-uuids",
	Usage:   "List Orders by IDs",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[[]string]{
			Name:      "identifier",
			Usage:     "List of client-assigned identifiers of orders to query. Maximum 100 orders.\n\n[Example] identifiers[]=id1&identifiers[]=id2…\n",
			QueryPath: "identifiers",
		},
		&requestflag.Flag[string]{
			Name:      "market",
			Usage:     "Trading pair to filter by.",
			QueryPath: "market",
		},
		&requestflag.Flag[string]{
			Name:      "order-by",
			Usage:     "Sort order. \"desc\" (newest first) or \"asc\" (oldest first). Default is \"desc\".\n",
			Default:   "desc",
			QueryPath: "order_by",
		},
		&requestflag.Flag[[]string]{
			Name:      "uuid",
			Usage:     "List of UUIDs of orders to query. Maximum 100 orders.\n\n[Example] uuids[]=uuid1&uuids[]=uuid2…\n",
			QueryPath: "uuids",
		},
	},
	Action:          handleOrdersListByUuids,
	HideHelpCommand: true,
}

var ordersListClosed = cli.Command{
	Name:    "list-closed",
	Usage:   "List Closed Orders",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "end-time",
			Usage:     "End time of the query range. Max range is 7 days.\n\n* If only \"end_time\" is given, the range is 7 days before it.\n\nFormat: ISO 8601 (2025-06-24T13:56:53+09:00) or millisecond timestamp (1750741013000)\n",
			QueryPath: "end_time",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Number of results per request (default: 100, max: 1,000).",
			Default:   100,
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "market",
			Usage:     "Trading pair to filter by.",
			QueryPath: "market",
		},
		&requestflag.Flag[string]{
			Name:      "order-by",
			Usage:     `Sort order. "desc" (newest first) or "asc" (oldest first). Default is "desc".`,
			Default:   "desc",
			QueryPath: "order_by",
		},
		&requestflag.Flag[string]{
			Name:      "start-time",
			Usage:     "Start time of the query range. Max range is 7 days.\n\n* If only \"start_time\" is given, the range is 7 days after it.\n* If neither is given, the default is the past 7 days.\n\nFormat: ISO 8601 (2025-06-24T13:56:53+09:00) or millisecond timestamp (1750741013000)\n",
			QueryPath: "start_time",
		},
		&requestflag.Flag[string]{
			Name:      "state",
			Usage:     "Closed order state.\n\n* `done`: Fully executed\n* `cancel`: Cancelled\n",
			QueryPath: "state",
		},
		&requestflag.Flag[[]string]{
			Name:      "states",
			Usage:     "Order state filter (array form). \"done\" or \"cancel\".\n\n[Example] states[]=done&states[]=cancel\n",
			Default:   []string{"done", "cancel"},
			QueryPath: "states",
		},
	},
	Action:          handleOrdersListClosed,
	HideHelpCommand: true,
}

var ordersListOpen = cli.Command{
	Name:    "list-open",
	Usage:   "List Open Orders",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Number of results per request (default: 100, max: 100).",
			Default:   100,
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "market",
			Usage:     "Trading pair to filter by.",
			QueryPath: "market",
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
			Usage:     "Open order state.\n\n* `wait`: Pending execution\n* `watch`: Pending reservation (stop order)\n",
			QueryPath: "state",
		},
		&requestflag.Flag[[]string]{
			Name:      "states",
			Usage:     "Order state filter (array form). \"wait\" or \"watch\".\n\n[Example] states[]=wait&states[]=watch\n",
			Default:   []string{"wait"},
			QueryPath: "states",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleOrdersListOpen,
	HideHelpCommand: true,
}

var ordersRetrieveChance = cli.Command{
	Name:    "retrieve-chance",
	Usage:   "Get Available Order Info",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "market",
			Usage:     "Trading pair to query.",
			Required:  true,
			QueryPath: "market",
		},
	},
	Action:          handleOrdersRetrieveChance,
	HideHelpCommand: true,
}

var ordersTestCreate = cli.Command{
	Name:    "test-create",
	Usage:   "Test Order Creation",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "market",
			Usage:    "Target trading pair for the order. (required)",
			Required: true,
			BodyPath: "market",
		},
		&requestflag.Flag[string]{
			Name:     "ord-type",
			Usage:    "Order type.\n\n* `limit`: Limit order\n* `price`: Market buy order (by total price)\n* `market`: Market sell order (by volume)\n* `best`: Best available price order\n",
			Required: true,
			BodyPath: "ord_type",
		},
		&requestflag.Flag[string]{
			Name:     "side",
			Usage:    "Order direction (ask=sell, bid=buy)",
			Required: true,
			BodyPath: "side",
		},
		&requestflag.Flag[string]{
			Name:     "identifier",
			Usage:    "Client-assigned order identifier.",
			BodyPath: "identifier",
		},
		&requestflag.Flag[string]{
			Name:     "price",
			Usage:    "Order price or total amount.",
			BodyPath: "price",
		},
		&requestflag.Flag[string]{
			Name:     "smp-type",
			Usage:    "Self-Match Prevention (SMP) mode.\n\n* `cancel_maker`: Cancel maker order\n* `cancel_taker`: Cancel taker order\n* `reduce`: Reduce order quantity\n",
			BodyPath: "smp_type",
		},
		&requestflag.Flag[string]{
			Name:     "time-in-force",
			Usage:    "Time in force condition.\n\n* `fok`: Fill or Kill\n* `ioc`: Immediate or Cancel\n* `post_only`: Post only (maker only)\n",
			BodyPath: "time_in_force",
		},
		&requestflag.Flag[string]{
			Name:     "volume",
			Usage:    "Order volume.",
			BodyPath: "volume",
		},
	},
	Action:          handleOrdersTestCreate,
	HideHelpCommand: true,
}

func handleOrdersCreate(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.OrderNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Orders.New(ctx, params, options...)
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
		Title:          "orders create",
		Transform:      transform,
	})
}

func handleOrdersRetrieve(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.OrderGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Orders.Get(ctx, params, options...)
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
		Title:          "orders retrieve",
		Transform:      transform,
	})
}

func handleOrdersCancel(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.OrderCancelParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Orders.Cancel(ctx, params, options...)
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
		Title:          "orders cancel",
		Transform:      transform,
	})
}

func handleOrdersCancelAndNew(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.OrderCancelAndNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Orders.CancelAndNew(ctx, params, options...)
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
		Title:          "orders cancel-and-new",
		Transform:      transform,
	})
}

func handleOrdersCancelByUuids(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.OrderCancelByUuidsParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Orders.CancelByUuids(ctx, params, options...)
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
		Title:          "orders cancel-by-uuids",
		Transform:      transform,
	})
}

func handleOrdersCancelOpen(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.OrderCancelOpenParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Orders.CancelOpen(ctx, params, options...)
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
		Title:          "orders cancel-open",
		Transform:      transform,
	})
}

func handleOrdersListByUuids(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.OrderListByUuidsParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Orders.ListByUuids(ctx, params, options...)
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
		Title:          "orders list-by-uuids",
		Transform:      transform,
	})
}

func handleOrdersListClosed(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.OrderListClosedParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Orders.ListClosed(ctx, params, options...)
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
		Title:          "orders list-closed",
		Transform:      transform,
	})
}

func handleOrdersListOpen(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.OrderListOpenParams{}

	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Orders.ListOpen(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "orders list-open",
			Transform:      transform,
		})
	} else {
		iter := client.Orders.ListOpenAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "orders list-open",
			Transform:      transform,
		})
	}
}

func handleOrdersRetrieveChance(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.OrderGetChanceParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Orders.GetChance(ctx, params, options...)
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
		Title:          "orders retrieve-chance",
		Transform:      transform,
	})
}

func handleOrdersTestCreate(ctx context.Context, cmd *cli.Command) error {
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

	params := upbit.OrderTestNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Orders.TestNew(ctx, params, options...)
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
		Title:          "orders test-create",
		Transform:      transform,
	})
}
