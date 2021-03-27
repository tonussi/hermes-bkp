package proxy

type Communicator interface {
	Listen(HandleIncomingMessageFunc) error
	Deliver([]byte) ([]byte, error)
}
