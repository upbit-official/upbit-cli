package cmd

import (
	"testing"

	"github.com/upbit-official/upbit-cli/internal/mocktest"
)

func TestDepositsRetrieve(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"deposits", "retrieve",
			"--currency", "BTC",
			"--txid", "98c15999f0bdc4ae0e8a-ed35868bb0c204fe6ec29e4058a3451e-88636d1040f4baddf943274ce37cf9cc",
			"--uuid", "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
		)
	})
}

func TestDepositsList(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"deposits", "list",
			"--max-items", "10",
			"--currency", "KRW",
			"--from", "from",
			"--limit", "100",
			"--order-by", "desc",
			"--page", "1",
			"--state", "PROCESSING",
			"--to", "to",
			"--txid", "98c15999f0bdc4ae0e8a-ed35868bb0c204fe6ec29e4058a3451e-88636d1040f4baddf943274ce37cf9cc",
			"--uuid", "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
		)
	})
}

func TestDepositsCreateCoinAddress(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"deposits", "create-coin-address",
			"--currency", "BTC",
			"--net-type", "BTC",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"currency: BTC\n" +
			"net_type: BTC\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--access-key", "string",
			"deposits", "create-coin-address",
		)
	})
}

func TestDepositsDepositKrw(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"deposits", "deposit-krw",
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
			"deposits", "deposit-krw",
		)
	})
}

func TestDepositsListCoinAddresses(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"deposits", "list-coin-addresses",
		)
	})
}

func TestDepositsRetrieveChance(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"deposits", "retrieve-chance",
			"--currency", "BTC",
			"--net-type", "BTC",
		)
	})
}

func TestDepositsRetrieveCoinAddress(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"deposits", "retrieve-coin-address",
			"--currency", "BTC",
			"--net-type", "BTC",
		)
	})
}
