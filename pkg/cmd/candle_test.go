package cmd

import (
	"testing"

	"github.com/upbit-official/upbit-cli/internal/mocktest"
)

func TestCandlesListDays(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"candles", "list-days",
			"--market", "KRW-BTC",
			"--converting-price-unit", "KRW",
			"--count", "200",
			"--to", "2024-01-01T00:00:00",
		)
	})
}

func TestCandlesListMinutes(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"candles", "list-minutes",
			"--unit", "15",
			"--market", "KRW-BTC",
			"--count", "200",
			"--to", "2024-01-01T00:00:00",
		)
	})
}

func TestCandlesListMonths(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"candles", "list-months",
			"--market", "KRW-BTC",
			"--count", "200",
			"--to", "2024-01-01T00:00:00",
		)
	})
}

func TestCandlesListSeconds(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"candles", "list-seconds",
			"--market", "KRW-BTC",
			"--count", "200",
			"--to", "2024-01-01T00:00:00",
		)
	})
}

func TestCandlesListWeeks(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"candles", "list-weeks",
			"--market", "KRW-BTC",
			"--count", "200",
			"--to", "2024-01-01T00:00:00",
		)
	})
}

func TestCandlesListYears(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--access-key", "string",
			"candles", "list-years",
			"--market", "KRW-BTC",
			"--count", "200",
			"--to", "2024-01-01T00:00:00",
		)
	})
}
