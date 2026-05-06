package cmd

import (
	"testing"

	"github.com/upbit-official/upbit-cli/internal/mocktest"
)

func TestOrderbooksList(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"orderbooks", "list",
			"--markets", "KRW-BTC,KRW-ETH",
			"--count", "10",
			"--level", "1",
		)
	})
}

func TestOrderbooksListInstruments(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"orderbooks", "list-instruments",
			"--markets", "KRW-BTC,KRW-ETH",
		)
	})
}
