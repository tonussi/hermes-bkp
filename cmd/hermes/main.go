package main

import (
	"flag"
	"log"
	"time"

	"github.com/r3musketeers/hermes/pkg/communication"
	hashicorpraft "github.com/r3musketeers/hermes/pkg/ordering/hashicorp-raft"
	"github.com/r3musketeers/hermes/pkg/proxy"
)

var (
	nodeID       = flag.String("i", "", "node id")
	clientsAddr  = flag.String("c", ":8001", "client requests address")
	deliveryAddr = flag.String("d", ":8000", "delivery server address")
	ordererAddr  = flag.String("o", "localhost:10012", "ordering protocol address bind")
	joinAddr     = flag.String("j", "", "join listener address")
)

func main() {
	flag.Parse()

	if *nodeID == "" {
		log.Fatal("node id cannot be empty")
	}

	tcpCommunicator, err := communication.NewTCPCommunicator(
		*clientsAddr,
		*deliveryAddr,
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	hashicoprRaftOrderer, err := hashicorpraft.NewHashicorpRaftOrderer(
		*nodeID,
		*ordererAddr,
		10*time.Second,
		"data/hashicor-raft/"+*nodeID,
		2,
		10*time.Second,
		*joinAddr == "",
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	proxy.Init(tcpCommunicator, hashicoprRaftOrderer)
	if err := proxy.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
