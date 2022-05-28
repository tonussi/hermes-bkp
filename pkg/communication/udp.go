package communication

// import (
// 	"net"

// 	"github.com/tonussi/studygo/proxy"
// )

// type UDPCommunicator struct {
// 	deliverConn *net.UDPConn
// }

// func NewUDPCommunicator(
// 	deliverAddr string,
// ) (*UDPCommunicator, error) {
// 	udpAddr, err := net.ResolveUDPAddr("udp4", deliverAddr)
// 	if err != nil {
// 		return nil, err
// 	}

// 	deliverConn, err := net.DialUDP("udp4", nil, udpAddr)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &UDPCommunicator{
// 		deliverConn: deliverConn,
// 	}, nil
// }

// func (comm UDPCommunicator) Listen(
// 	addr string,
// 	handle proxy.HandleIncomingMessageFunc,
// ) error {
// 	udpAddr, err := net.ResolveUDPAddr("udp4", addr)
// 	if err != nil {
// 		return err
// 	}

// 	listener, err := net.ListenUDP("udp4", udpAddr)
// 	if err != nil {
// 		return err
// 	}
// 	defer listener.Close()

// 	buffer := make([]byte, 1024)

// 	for {
// 		n, _, err := listener.ReadFromUDP(buffer)
// 		if err != nil {
// 			return err
// 		}

// 		_, err = handle(buffer[:n])
// 		if err != nil {
// 			return err
// 		}
// 	}
// }

// func (comm UDPCommunicator) Deliver(id string, data []byte) error {
// 	_, err := comm.deliverConn.Write(data)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
