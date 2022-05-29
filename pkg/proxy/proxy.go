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

type HermesProxy struct {
	communicator Communicator
	orderer      Orderer
}

func NewHermesProxy(
	communicator Communicator,
	orderer Orderer,
) *HermesProxy {
	return &HermesProxy{
		communicator: communicator,
		orderer:      orderer,
	}
}

func (proxy *HermesProxy) Run() error {
	proxy.orderer.SetOrderedMessageHandler(proxy.handleOrderedMessage)

	return proxy.communicator.Listen(proxy.handleIncomingMessage)
}

// Unexported functions

func (proxy *HermesProxy) handleIncomingMessage(data []byte) ([]byte, error) {
	return proxy.orderer.Process(data)
}

func (proxy *HermesProxy) handleOrderedMessage(data []byte) ([]byte, error) {
	return proxy.communicator.Deliver(data)
}
