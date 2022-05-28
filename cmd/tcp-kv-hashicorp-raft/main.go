package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tonussi/studygo/pkg/kv"
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
	keyRange       = flag.Int("r", 100000, "key range")
	valueSize      = flag.Int("v", 1024, "base value size for pre-population")
)

func main() {
	flag.Parse()

	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		log.Fatal("node id must be set")
	}

	raftAddr := os.Getenv("PROTOCOL_IP") + ":" + os.Getenv("PROTOCOL_PORT")

	tcpAddr, err := net.ResolveTCPAddr("tcp", *addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer listener.Close()

	fsm, err := NewFSM(
		nodeID,
		raftAddr,
		10*time.Second,
		"data/tcp-kv-hashicor-raft/"+nodeID,
		2,
		10*time.Second,
		*listenJoinAddr,
		*joinAddr,
		*valueSize,
		*keyRange,
	)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		logFile, err := os.Create(*logPath)
		if err != nil {
			log.Fatal(err)
		}

		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			delta := atomic.LoadUint64(&counter) - atomic.LoadUint64(&last)
			fmt.Fprintln(logFile, time.Now().UnixNano(), delta)
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

				resp, err := fsm.Process(req)
				if err != nil {
					conn.Write([]byte(err.Error()))
					continue
				}

				conn.Write(resp)

				atomic.AddUint64(&counter, 1)
			}
		}()
	}
}
