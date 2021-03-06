package kv

import "sync"

type KV struct {
	store map[uint64][]byte
	mux   sync.RWMutex
}

func NewKV() *KV {
	return &KV{
		store: map[uint64][]byte{},
	}
}

func (kv *KV) Get(key uint64) []byte {
	kv.mux.RLock()
	value, ok := kv.store[key]
	kv.mux.RUnlock()
	if !ok {
		return []byte{}
	}

	return value
}

func (kv *KV) Set(key uint64, value []byte) {
	if len(value) == 0 {
		return
	}

	kv.mux.Lock()
	kv.store[key] = value
	kv.mux.Unlock()
}

func (kv *KV) Delete(key uint64) {
	_, ok := kv.store[key]
	if !ok {
		return
	}

	kv.mux.Lock()
	delete(kv.store, key)
	kv.mux.Unlock()
}

func (kv *KV) Snapshot() map[uint64][]byte {
	snapshot := map[uint64][]byte{}

	kv.mux.RLock()
	for k, v := range kv.store {
		snapshot[k] = v
	}
	kv.mux.RUnlock()

	return snapshot
}
