package proxy

type Communicator interface {
	Listen(HandleIncomingMessageFunc) error
	Deliver(string, []byte) error
}
