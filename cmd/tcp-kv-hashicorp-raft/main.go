package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"net"
	"time"

	"github.com/r3musketeers/hermes/pkg/kv"
)

var (
	nodeID         = flag.String("i", "", "node id")
	addr           = flag.String("a", ":8000", "server address")
	logPath        = flag.String("l", "throughput.log", "path to log the throughput")
	bufferSize     = flag.Int("b", 2048, "requests buffer size")
	listenJoinAddr = flag.String("k", ":10000", "address to listen join requests")
	raftAddr       = flag.String("r", "localhost:11000", "ordering protocol address bind")
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

	fsm, err := NewFSM(
		*nodeID,
		*raftAddr,
		10*time.Second,
		"data/tcp-kv-hashicor-raft/"+*nodeID,
		2,
		10*time.Second,
		*listenJoinAddr,
		*joinAddr,
	)
	if err != nil {
		log.Fatal(err)
	}

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
					value := fsm.Get(req.Key)
					conn.Write(value)
					conn.Write([]byte("\n"))
				case kv.SetOp:
					err := fsm.Set(req.Key, req.Data)
					if err != nil {
						conn.Write([]byte(err.Error()))
					}
					conn.Write([]byte("\n"))
				case kv.DelOp:
					err := fsm.Delete(req.Key)
					if err != nil {
						conn.Write([]byte(err.Error()))
					}
					conn.Write([]byte("\n"))
				default:
					conn.Write([]byte("unsupported operation\n"))
				}
			}
		}()
	}
}
