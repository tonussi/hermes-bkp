package proxy

import (
	"log"
)

type HandleIncomingMessageFunc func(string, []byte) error

type HandleOrderedMessageFunc func(string, []byte) error

var (
	_communicator Communicator
	_orderer      Orderer
)

func Init(
	communicator Communicator,
	orderer Orderer,
) {
	_communicator = communicator
	_orderer = orderer
}

func Run() error {
	_orderer.SetOrderedMessageHandler(handleOrderedMessage)

	return _communicator.Listen(handleIncomingMessage)
}

// Unexported functions

func handleIncomingMessage(id string, data []byte) error {
	log.Println("handling incoming message")

	err := _orderer.Propose(id, data)
	if err != nil {
		return err
	}

	return nil
}

func handleOrderedMessage(id string, data []byte) error {
	log.Println("handling ordered message")

	return _communicator.Deliver(id, data)
}
