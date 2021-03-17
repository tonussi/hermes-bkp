package communication

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/r3musketeers/hermes/pkg/proxy"
)

type TCPCommunicator struct {
	listener    *net.TCPListener
	deliverConn *net.TCPConn

	connIDsMux   sync.RWMutex
	messageConns map[string]string

	connsMux    sync.RWMutex
	connsCount  int
	clientConns map[string]*net.TCPConn

	responseBuffer []byte
}

func NewTCPCommunicator(
	fromAddr string,
	toAddr string,
) (*TCPCommunicator, error) {
	listenAddr, err := net.ResolveTCPAddr("tcp", fromAddr)
	if err != nil {
		return nil, err
	}

	listener, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		return nil, err
	}

	deliverAddr, err := net.ResolveTCPAddr("tcp", toAddr)
	if err != nil {
		return nil, err
	}

	deliverConn, err := net.DialTCP("tcp", nil, deliverAddr)

	return &TCPCommunicator{
		listener:    listener,
		deliverConn: deliverConn,

		connIDsMux:   sync.RWMutex{},
		messageConns: map[string]string{},

		connsMux:    sync.RWMutex{},
		connsCount:  0,
		clientConns: map[string]*net.TCPConn{},

		responseBuffer: make([]byte, 1024),
	}, nil
}

func (comm TCPCommunicator) Listen(
	handle proxy.HandleIncomingMessageFunc,
) error {
	defer comm.listener.Close()

	for {
		conn, err := comm.listener.AcceptTCP()
		if err != nil {
			return err
		}

		go comm.handleConnection(conn, handle)
	}
}

func (comm TCPCommunicator) Deliver(id string, data []byte) error {
	defer func() {
		comm.connIDsMux.Lock()
		delete(comm.messageConns, id)
		comm.connIDsMux.Unlock()
	}()

	comm.deliverConn.Write(data)

	log.Println("delivered message")

	n, err := comm.deliverConn.Read(comm.responseBuffer)

	comm.connIDsMux.RLock()
	connID, ok := comm.messageConns[id]
	comm.connIDsMux.RUnlock()

	if !ok {
		log.Println("no need to respond the message")
		return nil
	}

	log.Println("responding for connection", connID)

	comm.connsMux.RLock()
	clientConn, ok := comm.clientConns[connID]
	comm.connsMux.RUnlock()

	if !ok {
		log.Println("client connection not here")
		return nil
	}

	if err != nil {
		clientConn.Write([]byte(err.Error()))
	} else {
		clientConn.Write(comm.responseBuffer[:n])
	}

	return nil
}

func (comm *TCPCommunicator) handleConnection(
	clientConn *net.TCPConn,
	handle proxy.HandleIncomingMessageFunc,
) {
	// saves client connection for future use
	comm.connsMux.Lock()
	connID := fmt.Sprintf("%d", comm.connsCount)
	comm.connsCount++
	log.Println("starting connection", connID)
	comm.clientConns[connID] = clientConn
	comm.connsMux.Unlock()

	errChan := make(chan error)
	quitChan := make(chan struct{})
	wg := sync.WaitGroup{}

	// starts reading from the client
	wg.Add(1)
	go func() {
		defer wg.Done()

		buffer := make([]byte, 1024)
		for {
			select {
			case <-quitChan:
				return
			default:
				err := clientConn.SetReadDeadline(time.Now().Add(2 * time.Second))
				if err != nil {
					log.Println("failed setting read deadline for client connection", connID)
					continue
				}

				n, err := clientConn.Read(buffer)
				if err != nil {
					if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
						continue
					}

					if errors.Is(err, io.EOF) {
						log.Println("client closed connection", connID)
						errChan <- nil
					} else {
						errChan <- fmt.Errorf("error reading from client: %w", err)
					}

					return
				}

				log.Print("message from client for connection", connID)

				// hash := sha256.New()
				// hash.Write(buffer[:n])
				// id := string(hash.Sum([]byte(clientConn.RemoteAddr().String())))
				id := uuid.NewString()

				comm.connIDsMux.Lock()
				comm.messageConns[id] = connID
				comm.connIDsMux.Unlock()

				err = handle(id, buffer[:n])
				if err != nil {
					clientConn.Write([]byte(err.Error()))
					continue
				}
			}
		}
	}()

	err := <-errChan
	if err != nil {
		log.Printf("error for connection %s: %s", connID, err.Error())
	}

	close(quitChan)
	close(errChan)

	wg.Wait()

	comm.connsMux.Lock()
	delete(comm.clientConns, connID)
	comm.connsMux.Unlock()

	log.Println("finished connection", connID)
}
