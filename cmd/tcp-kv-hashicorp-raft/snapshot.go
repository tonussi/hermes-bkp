package main

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/hashicorp/raft"
)

type HashicorpRaftSnapshot map[string][]byte

////////////////////////////////////////////////////////////////////////////////
//
// raft.FSMSnapshot interface implementation
//
////////////////////////////////////////////////////////////////////////////////

func (snapshot *HashicorpRaftSnapshot) Persist(sink raft.SnapshotSink) error {
	buffer := bytes.NewBuffer([]byte{})
	err := func() error {
		err := gob.NewEncoder(buffer).Encode(snapshot)
		if err != nil {
			return err
		}
		_, err = sink.Write(buffer.Bytes())
		if err != nil {
			return err
		}
		return sink.Close()
	}()
	if err != nil {
		cancelErr := sink.Cancel()
		log.Println(cancelErr.Error())
	}
	return err
}

func (snapshot *HashicorpRaftSnapshot) Release() {}
