package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/r3musketeers/hermes/pkg/kv"
)

var (
	counter    uint64
	last       uint64
	counterMux sync.RWMutex

	addr           = flag.String("a", ":8000", "server address")
	logPath        = flag.String("l", "throughput.log", "path to log the throughput")
	bufferSize     = flag.Int("b", 2048, "requests buffer size")
	listenJoinAddr = flag.String("k", ":9000", "address to listen join requests")
	joinAddr       = flag.String("j", "", "join listener address")
)

func main() {
	flag.Parse()

	tcpAddr, err := net.ResolveTCPAddr("tcp", *addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer listener.Close()

	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		log.Fatal("node id must be set")
	}

	raftAddr := os.Getenv("PROTOCOL_IP") + ":" + os.Getenv("PROTOCOL_PORT")

	fsm, err := NewFSM(
		nodeID,
		raftAddr,
		10*time.Second,
		"data/tcp-kv-hashicor-raft/"+nodeID,
		2,
		10*time.Second,
		*listenJoinAddr,
		*joinAddr,
	)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		logFile, err := os.Create(*logPath)
		if err != nil {
			log.Fatal(err)
		}

		logger := log.New(logFile, "", log.LstdFlags)

		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			delta := atomic.LoadUint64(&counter) - atomic.LoadUint64(&last)
			logger.Println(delta)
			atomic.StoreUint64(&last, atomic.LoadUint64(&counter))
		}
	}()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err.Error())
			continue
		}

		go func() {
			log.Println("connection started")

			buffer := make([]byte, *bufferSize)
			respCh := make(chan []byte, 1)

			for {
				n, err := conn.Read(buffer)
				if err != nil {
					if errors.Is(err, io.EOF) {
						log.Println("connection closed")
						break
					}
					conn.Write([]byte("bad message\n"))
					continue
				}

				req := kv.Request{}
				req.Parse(buffer[:n])

				err, resp := fsm.Process(req, respCh)
				if err != nil {
					conn.Write([]byte(err.Error()))
					conn.Write([]byte("\n"))
					continue
				}

				conn.Write(resp)
				conn.Write([]byte("\n"))

				atomic.AddUint64(&counter, 1)
			}
		}()
	}
}
