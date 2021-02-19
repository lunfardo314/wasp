package trie

import (
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetPut(t *testing.T) {
	t.Run("should get nothing if key does not exist", func(t *testing.T) {
		trie := NewTrie()
		v, err := trie.Get("notexist")
		require.NoError(t, err)
		require.Nil(t, v)
	})

	t.Run("should get value if key exist", func(t *testing.T) {
		trie := NewTrie()
		key := kv.Key([]byte{1, 2, 3, 4})
		trie.Set(key, []byte("hello"))
		v, err := trie.Get(key)
		require.NoError(t, err)
		require.EqualValues(t, v, []byte("hello"))
	})

	t.Run("should get updated value", func(t *testing.T) {
		trie := NewTrie()
		key := kv.Key([]byte{1, 2, 3, 4})
		trie.Set(key, []byte("hello"))
		trie.Set(key, []byte("world"))
		val, err := trie.Get(key)
		require.NoError(t, err)
		require.NotNil(t, val)
		require.EqualValues(t, val, []byte("world"))
	})
}

// verify data integrity
func TestDataIntegrity(t *testing.T) {
	t.Run("should get a different hash if a new key-value pair was added or updated", func(t *testing.T) {
		trie := NewTrie()
		hash0 := trie.Hash()
		key1 := kv.Key([]byte{1, 2, 3, 4})
		key2 := kv.Key([]byte{1, 2})

		trie.Set(key1, []byte("hello"))
		hash1 := trie.Hash()

		trie.Set(key2, []byte("world"))
		hash2 := trie.Hash()

		trie.Set(key2, []byte("trie"))
		hash3 := trie.Hash()

		require.NotEqual(t, hash0, hash1)
		require.NotEqual(t, hash1, hash2)
		require.NotEqual(t, hash2, hash3)
	})

	t.Run("should get the same hash if two tries have the identicial key-value pairs", func(t *testing.T) {
		trie1 := NewTrie()
		key1 := kv.Key([]byte{1, 2, 3, 4})
		key2 := kv.Key([]byte{1, 2})
		trie1.Set(key1, []byte("hello"))
		trie1.Set(key2, []byte("world"))

		trie2 := NewTrie()
		trie2.Set(key1, []byte("hello"))
		trie2.Set(key2, []byte("world"))

		require.Equal(t, trie1.Hash(), trie2.Hash())
	})
}
