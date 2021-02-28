package main

import (
	"errors"
	"io"
	"log"
	"net"
)

func main2() {
	// tcpCommunicator, err := communication.NewTCPCommunicator(":8001", ":8000")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// // udpCommunicator, err := communication.NewUDPCommunicator(":8001")
	// // if err != nil {
	// // 	log.Fatal(err.Error())
	// // }

	// //httpCommunicator := communication.NewHTTPCommunicator(":8000")

	// dummyOrderer := dummy.NewDummyOrderer()

	// proxy.Init(tcpCommunicator, dummyOrderer)
	// if err := proxy.Run(":8001"); err != nil {
	// 	log.Fatal(err.Error())
	// }

	upstreamChan := make(chan []byte)
	downstreamChan := make(chan []byte)

	go func() {
		addr, err := net.ResolveTCPAddr("tcp", ":8000")
		if err != nil {
			log.Fatal(err)
		}

		listener, err := net.ListenTCP("tcp", addr)
		if err != nil {
			log.Fatal(err)
		}

		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				log.Fatal(err)
			}

			log.Println("client connected", conn.RemoteAddr().String())

			go handleClientConnection(conn, upstreamChan, downstreamChan)
		}
	}()

	addr, err := net.ResolveTCPAddr("tcp", ":8001")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	go handleUpstreamConnection(conn, downstreamChan, upstreamChan)

	for {
	}
}

func handleClientConnection(
	conn *net.TCPConn,
	upstreamChan chan<- []byte,
	downstreamChan <-chan []byte,
) {
	buffer := make([]byte, 1024)
	clientChan := make(chan []byte)
	doneClientChan := make(chan interface{})

	go func() {
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					log.Println("client disconnected", conn.RemoteAddr().String())
				} else {
					log.Println(err)
				}
				doneClientChan <- nil
				return
			}

			clientChan <- buffer[:n]
		}
	}()

	for {
		select {
		case message := <-clientChan:
			upstreamChan <- message
		case message, ok := <-downstreamChan:
			if !ok {
				log.Println("proxy won't send messages to client anymore")
				doneClientChan <- nil
				return
			}

			conn.Write(message)
		case <-doneClientChan:
			close(clientChan)
		}
	}
}

func handleUpstreamConnection(
	conn *net.TCPConn,
	downstreamChan chan<- []byte,
	upstreamChan <-chan []byte,
) {
	buffer := make([]byte, 1024)
	serverChan := make(chan []byte)
	doneServerChan := make(chan interface{})

	go func() {
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					log.Println("upstream server disconnected", conn.RemoteAddr().String())
				} else {
					log.Println(err)
				}
				doneServerChan <- nil
				return
			}

			serverChan <- buffer[:n]
		}
	}()

	for {
		select {
		case message := <-serverChan:
			downstreamChan <- message
		case message, ok := <-upstreamChan:
			if !ok {
				log.Println("proxy won't send messages to upstream anymore")
				doneServerChan <- nil
				return
			}

			conn.Write(message)
		case <-doneServerChan:
			close(serverChan)
		}
	}
}
