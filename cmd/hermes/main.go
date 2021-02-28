package main

import (
	"log"

	"github.com/r3musketeers/hermes/communication"
	"github.com/r3musketeers/hermes/ordering/dummy"
	"github.com/r3musketeers/hermes/proxy"
)

func main() {
	tcpCommunicator, err := communication.NewTCPCommunicator(":8001", ":8000")
	if err != nil {
		log.Fatal(err.Error())
	}

	// udpCommunicator, err := communication.NewUDPCommunicator(":8001")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// httpCommunicator := communication.NewHTTPCommunicator(":8000")

	dummyOrderer := dummy.NewDummyOrderer()

	proxy.Init(tcpCommunicator, dummyOrderer)
	if err := proxy.Run(":8001"); err != nil {
		log.Fatal(err.Error())
	}
}
