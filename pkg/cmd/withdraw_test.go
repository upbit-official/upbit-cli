package cmd

import (
	"testing"

	"github.com/upbit-official/upbit-cli/internal/mocktest"
)

func TestWithdrawsRetrieve(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"withdraws", "retrieve",
			"--currency", "KRW",
			"--txid", "98c15999f0bdc4ae0e8a-ed35868bb0c204fe6ec29e4058a3451e-88636d1040f4baddf943274ce37cf9cc",
			"--uuid", "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
		)
	})
}

func TestWithdrawsList(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"withdraws", "list",
			"--max-items", "10",
			"--currency", "KRW",
			"--from", "from",
			"--limit", "100",
			"--order-by", "desc",
			"--page", "1",
			"--state", "WAITING",
			"--to", "to",
			"--txid", "98c15999f0bdc4ae0e8a-ed35868bb0c204fe6ec29e4058a3451e-88636d1040f4baddf943274ce37cf9cc",
			"--uuid", "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
		)
	})
}

func TestWithdrawsCancelWithdrawal(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"withdraws", "cancel-withdrawal",
			"--uuid", "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
		)
	})
}

func TestWithdrawsCreateKrwWithdrawal(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"withdraws", "create-krw-withdrawal",
			"--amount", "10000",
			"--two-factor-type", "naver",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"amount: '10000'\n" +
			"two_factor_type: naver\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--access-key", "string",
			"withdraws", "create-krw-withdrawal",
		)
	})
}

func TestWithdrawsCreateWithdrawal(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"withdraws", "create-withdrawal",
			"--address", "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
			"--amount", "0.01",
			"--currency", "BTC",
			"--net-type", "BTC",
			"--secondary-address", "secondary_address",
			"--transaction-type", "default",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"address: 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa\n" +
			"amount: '0.01'\n" +
			"currency: BTC\n" +
			"net_type: BTC\n" +
			"secondary_address: secondary_address\n" +
			"transaction_type: default\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--access-key", "string",
			"withdraws", "create-withdrawal",
		)
	})
}

func TestWithdrawsListCoinAddresses(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"withdraws", "list-coin-addresses",
		)
	})
}

func TestWithdrawsRetrieveChance(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"withdraws", "retrieve-chance",
			"--currency", "BTC",
			"--net-type", "BTC",
		)
	})
}
