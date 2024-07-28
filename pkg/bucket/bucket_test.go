package bucket

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBucket(t *testing.T) {
	b := New(10, 3, 300)

	t.Run("Check maxRetry", func(t *testing.T) {
		require.True(t, b.Check("test"))
		require.True(t, b.Check("test"))
		require.True(t, b.Check("test"))
		require.False(t, b.Check("test"))

		b.ResetKey("test")
	})

	t.Run("Check reset key", func(t *testing.T) {
		require.True(t, b.Check("test"))
		require.True(t, b.Check("test2"))
		require.True(t, b.Check("test2"))
		require.True(t, b.Check("test2"))
		require.True(t, b.Check("test3"))
		require.True(t, b.Check("test"))
		require.True(t, b.Check("test"))
		require.False(t, b.Check("test"))
		require.False(t, b.Check("test2"))
		b.ResetKey("test")
		require.True(t, b.Check("test"))
		require.False(t, b.Check("test2"))

		b.ResetKey("test")
		b.ResetKey("test2")
		b.ResetKey("test3")
	})

	t.Run("Check max tokens", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			require.True(t, b.Check(strconv.Itoa(i)))
		}
		require.False(t, b.Check("test"))

		for i := 0; i < 10; i++ {
			b.ResetKey(strconv.Itoa(i))
		}
		b.ResetKey("test")
	})

	t.Run("Check cleanup", func(t *testing.T) {
		b := New(10, 3, 1)

		require.True(t, b.Check("test"))
		require.True(t, b.Check("test"))
		require.True(t, b.Check("test"))
		require.False(t, b.Check("test"))
		time.Sleep(1 * time.Second)
		b.Cleanup()
		require.True(t, b.Check("test"))
	})
}
