package app

import (
  "github.com/hashicorp/golang-lru/v2/expirable"
  "github.com/henomis/lingoose/thread"
  "time"
)

func NewMemoryThreadStore(size int, ttl time.Duration) ThreadStore {
  // make cache with 10ms TTL and 5 max keys
  cache := expirable.NewLRU[string, *thread.Thread](size, nil, ttl)

  return &memThreadStore{
    cache: cache,
  }
}

type memThreadStore struct {
  cache *expirable.LRU[string, *thread.Thread]
}

func (m *memThreadStore) GetThread(threadID string) (*thread.Thread, error) {
  res, fnd := m.cache.Get(threadID)
  if !fnd {
    return nil, nil
  }

  return res, nil
}

func (m *memThreadStore) StoreThread(threadID string, thread *thread.Thread) error {
  m.cache.Add(threadID, thread)
  return nil
}
