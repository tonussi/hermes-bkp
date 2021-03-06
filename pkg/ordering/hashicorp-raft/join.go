package hashicorpraft

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/raft"
)

type JoinRequest struct {
	ID   string `json:"id"`
	Addr string `json:"addr"`
}

func (orderer *HashicorpRaftOrderer) join(id, addr string) error {
	log.Printf("received join request for remote node %s at %s", id, addr)

	configFuture := orderer.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		log.Printf("failed to get raft configuration: %v", err)
		return err
	}

	for _, srv := range configFuture.Configuration().Servers {
		// If a node already exists with either the joining node's ID or address,
		// that node may need to be removed from the config first.
		if srv.ID == raft.ServerID(id) || srv.Address == raft.ServerAddress(addr) {
			// However if *both* the ID and the address are the same, then nothing -- not even
			// a join operation -- is needed.
			if srv.Address == raft.ServerAddress(addr) && srv.ID == raft.ServerID(id) {
				log.Printf("node %s at %s already member of cluster, ignoring join request", id, addr)
				return nil
			}

			future := orderer.raft.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				return fmt.Errorf("error removing existing node %s at %s: %s", id, addr, err)
			}
		}
	}

	f := orderer.raft.AddVoter(raft.ServerID(id), raft.ServerAddress(addr), 0, 0)
	if f.Error() != nil {
		return f.Error()
	}
	log.Printf("node %s at %s joined successfully", id, addr)
	return nil
}

func (orderer *HashicorpRaftOrderer) joinHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("received join request")

	var joinReq JoinRequest
	err := json.NewDecoder(r.Body).Decode(&joinReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = orderer.join(joinReq.ID, joinReq.Addr)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
