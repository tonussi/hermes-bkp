package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/r3musketeers/hermes/pkg/kv"
)

type HashicorpRaftMessage struct {
	ID   string
	Data []byte
}

type FSM struct {
	nodeID          string
	address         string
	raft            *raft.Raft
	proposalTimeout time.Duration

	joinAddr string

	store   *kv.KV
	history map[string][]byte

	connCount      int
	messageConnMap map[string]int
	connMap        map[int]*net.TCPConn
	connMux        sync.RWMutex
}

func NewFSM(
	nodeID string,
	addrStr string,
	transportTimeout time.Duration,
	baseDir string,
	snapshotRetain int,
	proposalTimeout time.Duration,
	listenJoinAddr string,
	joinAddr string,
) (*FSM, error) {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeID)

	addr, err := net.ResolveTCPAddr("tcp", addrStr)
	if err != nil {
		return nil, err
	}

	transport, err := raft.NewTCPTransport(
		addrStr,
		addr,
		3,
		transportTimeout,
		os.Stderr,
	)
	if err != nil {
		return nil, err
	}

	snapshotStore, err := raft.NewFileSnapshotStore(
		baseDir,
		snapshotRetain,
		os.Stderr,
	)
	if err != nil {
		return nil, err
	}

	boltDB, err := raftboltdb.NewBoltStore(
		filepath.Join(baseDir, "raft.db"),
	)
	if err != nil {
		return nil, err
	}

	fsm := &FSM{
		nodeID:          nodeID,
		address:         addrStr,
		proposalTimeout: proposalTimeout,

		joinAddr: joinAddr,

		store:   kv.NewKV(),
		history: map[string][]byte{},

		messageConnMap: map[string]int{},
		connMap:        map[int]*net.TCPConn{},
	}

	raftInstance, err := raft.NewRaft(
		config,
		fsm,
		boltDB,
		boltDB,
		snapshotStore,
		transport,
	)
	if err != nil {
		return nil, err
	}

	fsm.raft = raftInstance

	if joinAddr == "" {
		configSingle := raft.Configuration{
			Servers: []raft.Server{
				{ID: config.LocalID, Address: transport.LocalAddr()},
			},
		}
		raftInstance.BootstrapCluster(configSingle)
	} else {
		joinBody, err := json.Marshal(map[string]string{
			"id":   nodeID,
			"addr": addrStr,
		})
		if err != nil {
			return nil, err
		}

		resp, err := http.Post(
			fmt.Sprintf("http://%s/hashicorp-raft/join", joinAddr),
			"application-type/json",
			bytes.NewReader(joinBody),
		)
		if err != nil {
			return nil, err
		}

		resp.Body.Close()
	}

	// listen join thread
	if listenJoinAddr != "" {
		go func() {
			http.HandleFunc("/hashicorp-raft/join", fsm.joinHandler)

			log.Println("listening join requests at", listenJoinAddr)
			http.ListenAndServe(listenJoinAddr, nil)
		}()
	}

	return fsm, nil
}

func (fsm *FSM) Process(req kv.Request) ([]byte, error) {
	if fsm.raft.State() != raft.Leader {
		return nil, errors.New("not a raft leader")
	}

	if req.Op != kv.GetOp && req.Op != kv.SetOp && req.Op != kv.DelOp {
		return nil, errors.New("unsupported operation")
	}

	messageID := uuid.NewString()

	buffer := bytes.NewBuffer([]byte{})
	gob.NewEncoder(buffer).Encode(
		HashicorpRaftMessage{ID: messageID, Data: req.Serialize()},
	)

	raftFuture := fsm.raft.Apply(buffer.Bytes(), fsm.proposalTimeout)

	err := raftFuture.Error()

	return raftFuture.Response().([]byte), err
}

////////////////////////////////////////////////////////////////////////////////
//
// raft.FSM interface implementation
//
////////////////////////////////////////////////////////////////////////////////

func (fsm *FSM) Apply(logEntry *raft.Log) interface{} {
	var message HashicorpRaftMessage

	buffer := bytes.NewReader(logEntry.Data)
	gob.NewDecoder(buffer).Decode(&message)

	req := kv.Request{}
	req.Parse(message.Data)

	resp := []byte{}

	switch req.Op {
	case kv.GetOp:
		resp = fsm.store.Get(req.Key)
	case kv.SetOp:
		fsm.store.Set(req.Key, req.Data)
	case kv.DelOp:
		fsm.store.Delete(req.Key)
	}

	return resp
}

func (fsm *FSM) Snapshot() (raft.FSMSnapshot, error) {
	history := map[string][]byte{}

	for id, message := range fsm.history {
		history[id] = message
	}

	snapshot := HashicorpRaftSnapshot(history)

	return &snapshot, nil
}

func (fsm *FSM) Restore(snapshotReader io.ReadCloser) error {
	snapshot := &HashicorpRaftSnapshot{}
	err := gob.NewDecoder(snapshotReader).Decode(snapshot)
	if err != nil {
		return err
	}

	for id, message := range *snapshot {
		fsm.history[id] = message
	}

	return nil
}
