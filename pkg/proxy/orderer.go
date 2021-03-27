package proxy

type Orderer interface {
	SetOrderedMessageHandler(HandleOrderedMessageFunc)
	Process([]byte) ([]byte, error)
}
