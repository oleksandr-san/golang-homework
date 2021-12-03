package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("clear cache removes all elements", func(t *testing.T) {
		c := NewCache(10)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		c.Clear()

		val, ok = c.Get("bbb")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge removes last element", func(t *testing.T) {
		c := NewCache(3)

		keys := []Key{"aaa", "bbb", "ccc", "ddd"}
		for value, key := range keys {
			require.False(t, c.Set(key, value))
		}

		val, ok := c.Get(keys[0])
		require.False(t, ok)
		require.Nil(t, val)

		for _, key := range keys[1:] {
			_, ok = c.Get(key)
			require.True(t, ok)
		}
	})

	t.Run("purge removes most unused element", func(t *testing.T) {
		c := NewCache(3)

		keys := []Key{"aaa", "bbb", "ccc"}
		for value, key := range keys {
			require.False(t, c.Set(key, value))
		}

		_, ok := c.Get("bbb")
		require.True(t, ok)

		_, ok = c.Get("aaa")
		require.True(t, ok)

		require.True(t, c.Set("ccc", 4))

		require.False(t, c.Set("ddd", 5))

		val, ok := c.Get("bbb")
		require.False(t, ok)
		require.Nil(t, val)

		for _, key := range []Key{"aaa", "ccc", "ddd"} {
			_, ok = c.Get(key)
			require.True(t, ok)
		}
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
