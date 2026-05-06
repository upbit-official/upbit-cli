package cmd

import (
	"testing"

	"github.com/upbit-official/upbit-cli/internal/mocktest"
)

func TestTickersListByQuoteCurrencies(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"tickers", "list-by-quote-currencies",
			"--quote-currencies", "KRW,BTC,USDT",
		)
	})
}

func TestTickersListByTradingPairs(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"tickers", "list-by-trading-pairs",
			"--markets", "KRW-BTC,KRW-ETH",
		)
	})
}
