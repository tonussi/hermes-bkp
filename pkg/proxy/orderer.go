package proxy

type Orderer interface {
	SetOrderedMessageHandler(HandleOrderedMessageFunc)
	Propose(string, []byte) error
}
