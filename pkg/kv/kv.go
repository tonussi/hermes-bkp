package kv

type KV struct {
	store map[uint64][]byte
}

func NewKV() *KV {
	return &KV{
		store: map[uint64][]byte{},
	}
}

func (kv *KV) Get(key uint64) []byte {
	return kv.store[key]
}

func (kv *KV) Set(key uint64, value []byte) {
	kv.store[key] = value
}

func (kv *KV) Delete(key uint64) {
	_, ok := kv.store[key]
	if !ok {
		return
	}

	delete(kv.store, key)
}

func (kv *KV) Snapshot() map[uint64][]byte {
	snapshot := map[uint64][]byte{}

	for k, v := range kv.store {
		snapshot[k] = v
	}

	return snapshot
}
