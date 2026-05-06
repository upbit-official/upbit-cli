package cmd

import (
	"testing"

	"github.com/upbit-official/upbit-cli/internal/mocktest"
)

func TestTravelRuleListVasps(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"travel-rule", "list-vasps",
		)
	})
}

func TestTravelRuleVerifyDepositByTxid(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"travel-rule", "verify-deposit-by-txid",
			"--currency", "ETH",
			"--net-type", "ETH",
			"--txid", "5b871d34-fe38-4025-8f5c-9b22028f85d3",
			"--vasp-uuid", "8d4fe968-82b2-42e5-822f-3840a245f802",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"currency: ETH\n" +
			"net_type: ETH\n" +
			"txid: 5b871d34-fe38-4025-8f5c-9b22028f85d3\n" +
			"vasp_uuid: 8d4fe968-82b2-42e5-822f-3840a245f802\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--access-key", "string",
			"travel-rule", "verify-deposit-by-txid",
		)
	})
}

func TestTravelRuleVerifyDepositByUuid(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"travel-rule", "verify-deposit-by-uuid",
			"--deposit-uuid", "5b871d34-fe38-4025-8f5c-9b22028f85d3",
			"--vasp-uuid", "8d4fe968-82b2-42e5-822f-3840a245f802",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"deposit_uuid: 5b871d34-fe38-4025-8f5c-9b22028f85d3\n" +
			"vasp_uuid: 8d4fe968-82b2-42e5-822f-3840a245f802\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--access-key", "string",
			"travel-rule", "verify-deposit-by-uuid",
		)
	})
}
