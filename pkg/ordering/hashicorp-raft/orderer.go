package hashicorpraft

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
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/r3musketeers/hermes/pkg/proxy"
)

type HashicorpRaftMessage struct {
	ID   string
	Data []byte
}

type HashicorpRaftOrderer struct {
	nodeID          string
	address         string
	raft            *raft.Raft
	proposalTimeout time.Duration

	listenJoinAddr string
	joinAddr       string

	orderedCh chan HashicorpRaftMessage

	history map[string][]byte
}

func NewHashicorpRaftOrderer(
	nodeID string,
	addrStr string,
	transportTimeout time.Duration,
	baseDir string,
	snapshotRetain int,
	proposalTimeout time.Duration,
	listenJoinAddr string,
	joinAddr string,
) (*HashicorpRaftOrderer, error) {
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
		10*time.Second,
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

	orderer := &HashicorpRaftOrderer{
		nodeID:          nodeID,
		address:         addrStr,
		proposalTimeout: proposalTimeout,
		listenJoinAddr:  listenJoinAddr,
		joinAddr:        joinAddr,

		orderedCh: make(chan HashicorpRaftMessage),
	}

	raftInstance, err := raft.NewRaft(
		config,
		orderer,
		boltDB,
		boltDB,
		snapshotStore,
		transport,
	)
	if err != nil {
		return nil, err
	}

	orderer.raft = raftInstance

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

	return orderer, nil
}

////////////////////////////////////////////////////////////////////////////////
//
// proxy.Orderer interface implementation
//
////////////////////////////////////////////////////////////////////////////////

func (orderer HashicorpRaftOrderer) Run(handle proxy.HandleOrderedMessageFunc) error {
	listenJoinErrCh := make(chan error)

	if orderer.listenJoinAddr != "" {
		go func() {
			http.HandleFunc("/hashicorp-raft/join", orderer.joinHandler)

			log.Println("listening join requests at", orderer.listenJoinAddr)
			listenJoinErrCh <- http.ListenAndServe(orderer.listenJoinAddr, nil)
		}()
	}

	for {
		select {
		case err := <-listenJoinErrCh:
			log.Println("error on join listener", err.Error())
			return err
		default:
			message := <-orderer.orderedCh

			err := handle(message.ID, message.Data)
			if err != nil {
				log.Println("error handling message", message.ID)
			}
		}
	}
}

func (orderer HashicorpRaftOrderer) Propose(id string, data []byte) error {
	if orderer.raft.State() != raft.Leader {
		return errors.New("not a raft leader")
	}

	buffer := bytes.NewBuffer([]byte{})
	gob.NewEncoder(buffer).Encode(
		HashicorpRaftMessage{ID: id, Data: data},
	)

	raftFuture := orderer.raft.Apply(buffer.Bytes(), orderer.proposalTimeout)

	return raftFuture.Error()
}

////////////////////////////////////////////////////////////////////////////////
//
// raft.FSM interface implementation
//
////////////////////////////////////////////////////////////////////////////////

func (orderer HashicorpRaftOrderer) Apply(logEntry *raft.Log) interface{} {
	// TODO: Add to local cache
	var message HashicorpRaftMessage

	buffer := bytes.NewReader(logEntry.Data)
	gob.NewDecoder(buffer).Decode(&message)

	orderer.orderedCh <- message
	return nil
}

func (orderer HashicorpRaftOrderer) Snapshot() (raft.FSMSnapshot, error) {
	return orderer.snapshot(), nil
}

func (orderer HashicorpRaftOrderer) Restore(snapshotReader io.ReadCloser) error {
	snapshot := &HashicorpRaftSnapshot{}
	err := gob.NewDecoder(snapshotReader).Decode(snapshot)
	if err != nil {
		return err
	}

	orderer.restore(snapshot)

	return nil
}

// Unexported methods

func (orderer HashicorpRaftOrderer) snapshot() *HashicorpRaftSnapshot {
	history := map[string][]byte{}

	for id, message := range orderer.history {
		history[id] = message
	}

	snapshot := HashicorpRaftSnapshot(history)

	return &snapshot
}

func (orderer *HashicorpRaftOrderer) restore(
	snapshot *HashicorpRaftSnapshot,
) {
	for id, message := range *snapshot {
		orderer.history[id] = message
	}
}
