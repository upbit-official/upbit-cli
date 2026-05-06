package cmd

import (
	"testing"

	"github.com/upbit-official/upbit-cli/internal/mocktest"
)

func TestTradesList(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"trades", "list",
			"--market", "KRW-BTC",
			"--count", "500",
			"--cursor", "cursor",
			"--days-ago", "1",
			"--to", "134501",
		)
	})
}
