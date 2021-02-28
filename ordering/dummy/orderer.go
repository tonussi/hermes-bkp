package dummy

import (
	"time"

	"github.com/r3musketeers/hermes/proxy"
)

type DummyMessage struct {
	ID   string
	Data []byte
}

type DummyOrderer struct {
	orderedCh chan DummyMessage
}

func NewDummyOrderer() *DummyOrderer {
	return &DummyOrderer{
		orderedCh: make(chan DummyMessage),
	}
}

func (orderer DummyOrderer) Run(handle proxy.HandleOrderedMessageFunc) error {
	for message := range orderer.orderedCh {
		err := handle(message.ID, message.Data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (orderer DummyOrderer) Propose(id string, data []byte) error {
	time.Sleep(2 * time.Second)

	orderer.orderedCh <- DummyMessage{ID: id, Data: data}

	return nil
}
