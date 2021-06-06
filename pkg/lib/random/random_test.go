package random

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBase62String(t *testing.T) {
	t.Run("short", func(t *testing.T) {
		result := Base62String(10)
		require.Len(t, result, 10)
	})
	t.Run("long", func(t *testing.T) {
		result := Base62String(1000)
		require.Len(t, result, 1000)
	})
}

func TestHexString(t *testing.T) {
	t.Run("short", func(t *testing.T) {
		result := HexString(10)
		require.Len(t, result, 10)
	})
	t.Run("long", func(t *testing.T) {
		result := HexString(1000)
		require.Len(t, result, 1000)
	})
}
