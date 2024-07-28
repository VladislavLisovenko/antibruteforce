package list

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

const (
	redisKey           = "test"
	addrWithMask       = "10.2.2.0/24"
	trueIP             = "10.2.2.10"
	falseIP            = "10.2.3.10"
	presetAddrWithMask = "10.2.4.0/25"
	presetIP           = "10.2.4.125"
	falsePresetIP      = "10.2.4.128"
)

func TestList(t *testing.T) {
	s := miniredis.RunT(t)
	_, err := s.SAdd(redisKey, presetAddrWithMask)
	require.NoError(t, err)
	r := redis.NewClient(&redis.Options{Addr: s.Addr()})
	ctx := context.Background()
	l, err := New(ctx, redisKey, *r)
	require.NoError(t, err)

	t.Run("Check preset mask", func(t *testing.T) {
		require.True(t, l.Check(presetIP))
		require.False(t, l.Check(falsePresetIP))
	})

	t.Run("Check ip by mask", func(t *testing.T) {
		err = l.Add(ctx, addrWithMask)
		require.NoError(t, err)
		require.True(t, l.Check(trueIP))
		require.False(t, l.Check(falseIP))
	})

	t.Run("Check delete mask", func(t *testing.T) {
		err = l.Delete(ctx, addrWithMask)
		require.NoError(t, err)
		require.False(t, l.Check(trueIP))
	})

	t.Run("Check reset cache", func(t *testing.T) {
		err = l.Add(ctx, addrWithMask)
		require.NoError(t, err)
		require.True(t, l.Check(trueIP))
		err = l.Reset(ctx)
		require.NoError(t, err)
		require.False(t, l.Check(trueIP))
	})

	t.Run("Check redis", func(t *testing.T) {
		err = l.Add(ctx, addrWithMask)
		elements, err := r.SMembers(ctx, redisKey).Result()
		require.NoError(t, err)
		require.Equal(t, 1, len(elements))
		require.Equal(t, addrWithMask, elements[0])
	})

	t.Run("Check invalid ip", func(t *testing.T) {
		err := l.Add(ctx, "test")
		require.Error(t, err)

		err = l.Delete(ctx, "test")
		require.Error(t, err)

		require.False(t, l.Check("test"))
	})
}
