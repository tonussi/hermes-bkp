package main

import (
	"encoding/json"
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

	addr       = flag.String("a", ":8001", "server address")
	logPath    = flag.String("l", "throughput.log", "path to log the throughput")
	bufferSize = flag.Int("b", 2048, "requests buffer size")
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

	store := kv.NewKV()

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

				switch req.Op {
				case kv.GetOp:
					value := store.Get(req.Key)
					conn.Write(value)
					conn.Write([]byte("\n"))
				case kv.SetOp:
					store.Set(req.Key, req.Data)
					conn.Write([]byte("\n"))
				case kv.DelOp:
					store.Delete(req.Key)
					conn.Write([]byte("\n"))
				case kv.SnapOp:
					snapshot := store.Snapshot()
					json.NewEncoder(conn).Encode(snapshot)
				default:
					conn.Write([]byte("unsupported operation\n"))
				}

				atomic.AddUint64(&counter, 1)
			}
		}()
	}
}
