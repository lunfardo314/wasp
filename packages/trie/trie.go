package trie

import (
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/dict"
)

type Trie interface {
	kv.KVStore
	Hash() hashing.HashValue
}

type SimpleTrie struct {
	d dict.Dict
}

func NewTrie() Trie {
	return &SimpleTrie{
		d: dict.New(),
	}
}

func (s *SimpleTrie) Set(key kv.Key, value []byte) {
	s.d.Set(key, value)
}

func (s *SimpleTrie) Del(key kv.Key) {
	s.d.Del(key)
}

func (s *SimpleTrie) Get(key kv.Key) ([]byte, error) {
	return s.d.Get(key)
}

func (s *SimpleTrie) Has(key kv.Key) (bool, error) {
	return s.d.Has(key)
}

func (s *SimpleTrie) MustGet(key kv.Key) []byte {
	panic("implement me")
}

func (s *SimpleTrie) MustHas(key kv.Key) bool {
	panic("implement me")
}

func (s *SimpleTrie) Iterate(prefix kv.Key, f func(key kv.Key, value []byte) bool) error {
	panic("implement me")
}

func (s *SimpleTrie) IterateKeys(prefix kv.Key, f func(key kv.Key) bool) error {
	panic("implement me")
}

func (s *SimpleTrie) MustIterate(prefix kv.Key, f func(key kv.Key, value []byte) bool) {
	panic("implement me")
}

func (s *SimpleTrie) MustIterateKeys(prefix kv.Key, f func(key kv.Key) bool) {
	panic("implement me")
}

func (s *SimpleTrie) Hash() hashing.HashValue {
	return s.d.Hash()
}
