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
	"sync/atomic"
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

	counter    uint64
	last       uint64
	counterMux sync.RWMutex

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

	// throughput log thread
	go func() {
		logFile, err := os.Create(*logPath)
		if err != nil {
			log.Fatal(err)
		}

		logger := log.New(logFile, "", log.LstdFlags)

		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			delta := atomic.LoadUint64(&fsm.counter) - atomic.LoadUint64(&fsm.last)
			logger.Println(delta)
			atomic.StoreUint64(&fsm.last, atomic.LoadUint64(&fsm.counter))
		}
	}()

	return fsm, nil
}

func (fsm *FSM) Get(key uint64) []byte {
	defer atomic.AddUint64(&fsm.counter, 1)
	return fsm.store.Get(key)
}

func (fsm *FSM) Set(key uint64, value []byte) error {
	if fsm.raft.State() != raft.Leader {
		return errors.New("not a raft leader")
	}

	req := kv.Request{Op: kv.SetOp, Key: key, Data: value}

	buffer := bytes.NewBuffer([]byte{})
	gob.NewEncoder(buffer).Encode(
		HashicorpRaftMessage{ID: uuid.NewString(), Data: req.Serialize()},
	)

	raftFuture := fsm.raft.Apply(buffer.Bytes(), fsm.proposalTimeout)

	return raftFuture.Error()
}

func (fsm *FSM) Delete(key uint64) error {
	if fsm.raft.State() != raft.Leader {
		return errors.New("not a raft leader")
	}

	req := kv.Request{Op: kv.DelOp, Key: key}

	buffer := bytes.NewBuffer([]byte{})
	gob.NewEncoder(buffer).Encode(
		HashicorpRaftMessage{ID: uuid.NewString(), Data: req.Serialize()},
	)

	raftFuture := fsm.raft.Apply(buffer.Bytes(), fsm.proposalTimeout)

	return raftFuture.Error()
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

	switch req.Op {
	case kv.SetOp:
		fsm.store.Set(req.Key, req.Data)
	case kv.DelOp:
		fsm.store.Delete(req.Key)
	}

	atomic.AddUint64(&fsm.counter, 1)

	return nil
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
