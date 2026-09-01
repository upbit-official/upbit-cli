package cmd

import (
	"testing"

	"github.com/upbit-official/upbit-cli/internal/mocktest"
)

func TestPocketsList(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"pockets", "list",
		)
	})
}

func TestPocketsListAPIKeys(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"pockets", "list-api-keys",
			"--include-expired=true",
			"--uuid", "string",
		)
	})
}

func TestPocketsListTransfers(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"pockets", "list-transfers",
			"--currency", "XRP",
			"--direction", "out",
			"--end-time", "end_time",
			"--identifier", "string",
			"--limit", "2",
			"--order-by", "desc",
			"--start-time", "start_time",
			"--state", "processing",
			"--uuid", "string",
		)
	})
}

func TestPocketsListUniversalTransfers(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"pockets", "list-universal-transfers",
			"--currency", "XRP",
			"--end-time", "end_time",
			"--from", "from",
			"--identifier", "string",
			"--limit", "2",
			"--order-by", "desc",
			"--start-time", "start_time",
			"--state", "processing",
			"--to", "to",
			"--uuid", "string",
		)
	})
}

func TestPocketsRetrieveBalance(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"pockets", "retrieve-balance",
			"--uuid", "uuid",
		)
	})
}

func TestPocketsTransfer(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"pockets", "transfer",
			"--amount", "10.00",
			"--currency", "XRP",
			"--to", "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
			"--identifier", "unique_transfer_id_20260531_01",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"amount: '10.00'\n" +
			"currency: XRP\n" +
			"to: 9ca023a5-851b-4fec-9f0a-48cd83c2eaae\n" +
			"identifier: unique_transfer_id_20260531_01\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--access-key", "string",
			"pockets", "transfer",
		)
	})
}

func TestPocketsUniversalTransfer(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"pockets", "universal-transfer",
			"--amount", "10.00",
			"--currency", "XRP",
			"--to", "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
			"--from", "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
			"--identifier", "unique_transfer_id_20260531_01",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"amount: '10.00'\n" +
			"currency: XRP\n" +
			"to: 9ca023a5-851b-4fec-9f0a-48cd83c2eaae\n" +
			"from: 9ca023a5-851b-4fec-9f0a-48cd83c2eaae\n" +
			"identifier: unique_transfer_id_20260531_01\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--access-key", "string",
			"pockets", "universal-transfer",
		)
	})
}
