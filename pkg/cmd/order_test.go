package cmd

import (
	"testing"

	"github.com/upbit-official/upbit-cli/internal/mocktest"
)

func TestOrdersCreate(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"orders", "create",
			"--market", "KRW-BTC",
			"--ord-type", "limit",
			"--side", "bid",
			"--identifier", "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
			"--price", "14000000",
			"--smp-type", "cancel_maker",
			"--time-in-force", "fok",
			"--volume", "1",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"market: KRW-BTC\n" +
			"ord_type: limit\n" +
			"side: bid\n" +
			"identifier: 9ca023a5-851b-4fec-9f0a-48cd83c2eaae\n" +
			"price: '14000000'\n" +
			"smp_type: cancel_maker\n" +
			"time_in_force: fok\n" +
			"volume: '1'\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--access-key", "string",
			"orders", "create",
		)
	})
}

func TestOrdersRetrieve(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"orders", "retrieve",
			"--identifier", "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
			"--uuid", "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
		)
	})
}

func TestOrdersCancel(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"orders", "cancel",
			"--identifier", "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
			"--uuid", "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
		)
	})
}

func TestOrdersCancelAndNew(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"orders", "cancel-and-new",
			"--new-ord-type", "limit",
			"--new-identifier", "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
			"--new-price", "100000000",
			"--new-smp-type", "cancel_maker",
			"--new-time-in-force", "fok",
			"--new-volume", "remain_only",
			"--prev-order-identifier", "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
			"--prev-order-uuid", "ad217e24-ed02-469c-9b30-c08dbbda6908",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"new_ord_type: limit\n" +
			"new_identifier: 9ca023a5-851b-4fec-9f0a-48cd83c2eaae\n" +
			"new_price: '100000000'\n" +
			"new_smp_type: cancel_maker\n" +
			"new_time_in_force: fok\n" +
			"new_volume: remain_only\n" +
			"prev_order_identifier: 9ca023a5-851b-4fec-9f0a-48cd83c2eaae\n" +
			"prev_order_uuid: ad217e24-ed02-469c-9b30-c08dbbda6908\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--access-key", "string",
			"orders", "cancel-and-new",
		)
	})
}

func TestOrdersCancelByUuids(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"orders", "cancel-by-uuids",
			"--identifier", "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
			"--uuid", "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
		)
	})
}

func TestOrdersCancelOpen(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"orders", "cancel-open",
			"--cancel-side", "bid",
			"--count", "300",
			"--excluded-pairs", "KRW-BTC,KRW-ETH",
			"--order-by", "desc",
			"--pairs", "KRW-BTC,KRW-ETH",
			"--quote-currencies", "KRW",
		)
	})
}

func TestOrdersListByUuids(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"orders", "list-by-uuids",
			"--identifier", "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
			"--market", "KRW-BTC",
			"--order-by", "desc",
			"--uuid", "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
		)
	})
}

func TestOrdersListClosed(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("state flag", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"orders", "list-closed",
			"--end-time", "2024-12-09T13:56:53+09:00",
			"--limit", "1000",
			"--market", "KRW-BTC",
			"--order-by", "desc",
			"--start-time", "2024-12-09T13:56:53+09:00",
			"--state", "done",
		)
	})
	t.Run("states flag", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"orders", "list-closed",
			"--end-time", "2024-12-09T13:56:53+09:00",
			"--limit", "1000",
			"--market", "KRW-BTC",
			"--order-by", "desc",
			"--start-time", "2024-12-09T13:56:53+09:00",
			"--states", "done",
			"--states", "cancel",
		)
	})
}

func TestOrdersListOpen(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("state flag", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"orders", "list-open",
			"--max-items", "10",
			"--limit", "100",
			"--market", "KRW-BTC",
			"--order-by", "desc",
			"--page", "1",
			"--state", "wait",
		)
	})
	t.Run("states flag", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"orders", "list-open",
			"--max-items", "10",
			"--limit", "100",
			"--market", "KRW-BTC",
			"--order-by", "desc",
			"--page", "1",
			"--states", "wait",
			"--states", "watch",
		)
	})
}

func TestOrdersRetrieveChance(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"orders", "retrieve-chance",
			"--market", "KRW-BTC",
		)
	})
}

func TestOrdersTestCreate(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"orders", "test-create",
			"--market", "KRW-BTC",
			"--ord-type", "limit",
			"--side", "bid",
			"--identifier", "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
			"--price", "14000000",
			"--smp-type", "cancel_maker",
			"--time-in-force", "fok",
			"--volume", "1",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"market: KRW-BTC\n" +
			"ord_type: limit\n" +
			"side: bid\n" +
			"identifier: 9ca023a5-851b-4fec-9f0a-48cd83c2eaae\n" +
			"price: '14000000'\n" +
			"smp_type: cancel_maker\n" +
			"time_in_force: fok\n" +
			"volume: '1'\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--access-key", "string",
			"orders", "test-create",
		)
	})
}
