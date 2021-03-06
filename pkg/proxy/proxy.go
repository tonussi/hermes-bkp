package proxy

import (
	"log"
)

type HandleIncomingMessageFunc func(string, []byte) error

type HandleOrderedMessageFunc func(string, []byte) error

var (
	_communicator Communicator
	_orderer      Orderer

	_orderedCh chan Message
)

func Init(
	communicator Communicator,
	orderer Orderer,
) {
	_communicator = communicator
	_orderer = orderer
}

func Run() error {
	errCh := make(chan error)

	go func() {
		errCh <- _communicator.Listen(handleIncomingMessage)
	}()

	go func() {
		errCh <- _orderer.Run(handleOrderedMessage)
	}()

	go func() {
		for message := range _orderedCh {
			_communicator.Deliver(message.ID, message.Data)
		}
	}()

	err := <-errCh
	if err != nil {
		close(_orderedCh)
	}

	return err
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

	_orderedCh <- Message{ID: id, Data: data}

	return nil
}
