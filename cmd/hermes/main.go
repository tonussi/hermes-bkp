package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/r3musketeers/hermes/pkg/communication"
	hashicorpraft "github.com/r3musketeers/hermes/pkg/ordering/hashicorp-raft"
	"github.com/r3musketeers/hermes/pkg/proxy"
)

var (
	listenAddr     = flag.String("l", ":8000", "listen requests address")
	deliveryAddr   = flag.String("d", ":8001", "delivery server address")
	listenJoinAddr = flag.String("k", ":9000", "listen join requests address")
	// bufferSize     = flag.Int("b", 2048, "requests buffer size")
	joinAddr = flag.String("j", "", "join address")
)

func main() {
	flag.Parse()

	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		log.Fatal("node id must be set")
	}

	raftAddr := os.Getenv("PROTOCOL_IP") + ":" + os.Getenv("PROTOCOL_PORT")

	httpCommunicator, err := communication.NewHTTPCommunicator(
		*listenAddr,
		*deliveryAddr,
		5,
		2*time.Second,
		// *bufferSize,
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	hashicoprRaftOrderer, err := hashicorpraft.NewHashicorpRaftOrderer(
		nodeID,
		raftAddr,
		10*time.Second,
		"data/hermes/hashicor-raft/"+nodeID,
		2,
		10*time.Second,
		*listenJoinAddr,
		*joinAddr,
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	hermes := proxy.NewHermesProxy(httpCommunicator, hashicoprRaftOrderer)

	err = hermes.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}
