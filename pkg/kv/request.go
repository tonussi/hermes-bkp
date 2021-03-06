package kv

import (
	"encoding/binary"
	"errors"
)

type Request struct {
	Op   Op
	Key  uint64
	Data []byte
}

func (req *Request) Parse(reqBytes []byte) error {
	if len(reqBytes) < OpByteSize+KeyByteSize {
		return errors.New("bad request")
	}

	opBytes := reqBytes[:OpByteSize]
	keyBytes := reqBytes[OpByteSize : OpByteSize+KeyByteSize]
	dataBytes := reqBytes[OpByteSize+KeyByteSize:]

	op, _ := binary.Uvarint(opBytes)
	key, _ := binary.Uvarint(keyBytes)

	req.Op = Op(op)
	req.Key = key
	req.Data = dataBytes

	return nil
}

func (req Request) Serialize() []byte {
	reqBytes := make([]byte, OpByteSize+KeyByteSize+len(req.Data))

	binary.PutUvarint(reqBytes[:OpByteSize], uint64(req.Op))
	binary.PutUvarint(reqBytes[OpByteSize:OpByteSize+KeyByteSize], req.Key)
	copy(reqBytes[OpByteSize+KeyByteSize:], req.Data)

	return reqBytes
}
