package proxy

type Orderer interface {
	Run(HandleOrderedMessageFunc) error
	Propose(string, []byte) error
}
