package proxy

type HandleIncomingMessageFunc func([]byte) ([]byte, error)

type HandleOrderedMessageFunc func([]byte) ([]byte, error)

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

func handleIncomingMessage(data []byte) ([]byte, error) {
	return _orderer.Process(data)
}

func handleOrderedMessage(data []byte) ([]byte, error) {
	return _communicator.Deliver(data)
}
