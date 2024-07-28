package ratelimit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRateLimit(t *testing.T) {
	rt := New(2, 5, 10, 10, 10)

	t.Run("Check max retry", func(t *testing.T) {
		require.True(t, rt.Check("test", "test2", "127.0.0.1"))
		require.True(t, rt.Check("test", "test2", "127.0.0.1"))
		require.False(t, rt.Check("test", "test2", "127.0.0.1"))
		require.True(t, rt.Check("test2", "test2", "127.0.0.1"))
		require.True(t, rt.Check("test3", "test2", "127.0.0.1"))
		require.False(t, rt.Check("test4", "test2", "127.0.0.1"))
		rt.Reset("test", "127.0.0.1")
		rt.Reset("test2", "127.0.0.1")
		rt.Reset("test3", "127.0.0.1")
		rt.Reset("test4", "127.0.0.1")
	})

	t.Run("Check reset", func(t *testing.T) {
		// Check reset login
		require.True(t, rt.Check("test", "test3", "127.0.0.1"))
		require.True(t, rt.Check("test", "test3", "127.0.0.1"))
		require.False(t, rt.Check("test", "test3", "127.0.0.1"))
		rt.Reset("test", "127.0.0.1")
		require.True(t, rt.Check("test", "test3", "127.0.0.1"))

		require.True(t, rt.Check("test2", "test4", "127.0.0.1"))
		require.True(t, rt.Check("test2", "test4", "127.0.0.1"))
		require.True(t, rt.Check("test3", "test4", "127.0.0.1"))
		require.True(t, rt.Check("test3", "test4", "127.0.0.1"))
		require.True(t, rt.Check("test4", "test4", "127.0.0.1"))
		require.True(t, rt.Check("test4", "test5", "127.0.0.1"))
		require.True(t, rt.Check("test5", "test5", "127.0.0.1"))
		require.True(t, rt.Check("test5", "test5", "127.0.0.1"))
		require.True(t, rt.Check("test6", "test5", "127.0.0.1"))
		require.False(t, rt.Check("test6", "test5", "127.0.0.1"))
		rt.Reset("test", "127.0.0.1")
		require.True(t, rt.Check("test7", "test6", "127.0.0.1"))
	})

	t.Run("Check cleanup", func(t *testing.T) {
		rt := New(2, 5, 10, 10, 1)

		require.True(t, rt.Check("test", "test2", "127.0.0.1"))
		require.True(t, rt.Check("test", "test2", "127.0.0.1"))
		require.False(t, rt.Check("test", "test2", "127.0.0.1"))
		time.Sleep(1 * time.Second)
		rt.Cleanup()
		require.True(t, rt.Check("test", "test2", "127.0.0.1"))
	})
}
