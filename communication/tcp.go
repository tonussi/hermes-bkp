package communication

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/r3musketeers/hermes/proxy"
)

type TCPCommunicator struct {
	listener    *net.TCPListener
	deliverAddr *net.TCPAddr

	connIDsMux   sync.RWMutex
	messageConns map[string]string

	connsMux    sync.RWMutex
	connsCount  int
	clientConns map[string]*net.TCPConn
	serverConns map[string]*net.TCPConn
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

	return &TCPCommunicator{
		listener:    listener,
		deliverAddr: deliverAddr,

		connIDsMux:   sync.RWMutex{},
		messageConns: map[string]string{},

		connsMux:    sync.RWMutex{},
		connsCount:  0,
		clientConns: map[string]*net.TCPConn{},
		serverConns: map[string]*net.TCPConn{},
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

	comm.connIDsMux.RLock()
	connID := comm.messageConns[id]
	comm.connIDsMux.RUnlock()

	log.Println("delivering for connection", connID)

	comm.connsMux.RLock()
	serverConn := comm.serverConns[connID]
	comm.connsMux.RUnlock()

	serverConn.Write(data)

	log.Println("delivered message")

	return nil
}

func (comm *TCPCommunicator) handleConnection(
	clientConn *net.TCPConn,
	handle proxy.HandleIncomingMessageFunc,
) {
	// starts connection with server
	serverConn, err := net.DialTCP("tcp", nil, comm.deliverAddr)
	if err != nil {
		clientConn.Write([]byte(err.Error()))
		return
	}

	// saves client and server connections for future use
	comm.connsMux.Lock()
	// connID := uuid.New().String()
	connID := fmt.Sprintf("%d", comm.connsCount)
	comm.connsCount++
	log.Println("starting connection", connID)
	comm.clientConns[connID] = clientConn
	comm.serverConns[connID] = serverConn
	comm.connsMux.Unlock()

	errChan := make(chan error)
	quitChan := make(chan struct{})
	wg := sync.WaitGroup{}

	// starts reading from the server
	wg.Add(1)
	go func() {
		defer wg.Done()

		buffer := make([]byte, 1024)
		for {
			select {
			case <-quitChan:
				return
			default:
				err := serverConn.SetReadDeadline(time.Now().Add(2 * time.Second))
				if err != nil {
					log.Println("failed setting read deadline for server connection", connID)
					continue
				}

				n, err := serverConn.Read(buffer)
				if err != nil {
					if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
						continue
					}

					if errors.Is(err, io.EOF) {
						log.Println("server closed connection", connID)
						errChan <- nil
					} else {
						errChan <- fmt.Errorf("error reading from client: %w", err)
					}

					return
				}

				log.Println("message from server for connection", connID)

				clientConn.Write(buffer[:n])

				log.Println("replied message on connection", connID)
			}
		}
	}()

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

				hash := sha256.New()
				hash.Write(buffer[:n])
				id := string(hash.Sum([]byte(clientConn.RemoteAddr().String())))

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

	err = <-errChan
	if err != nil {
		log.Printf("error for connection %s: %s", connID, err.Error())
	}

	close(quitChan)
	close(errChan)

	wg.Wait()

	serverConn.Close()

	comm.connsMux.Lock()
	delete(comm.clientConns, connID)
	delete(comm.serverConns, connID)
	comm.connsMux.Unlock()

	log.Println("finished connection", connID)
}
