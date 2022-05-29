package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"runtime"
	"sync"
	"time"

	"github.com/r3musketeers/hermes/pkg/kv"
)

var (
	serverAddr   = flag.String("s", ":8000", "server address")
	duration     = flag.Duration("d", time.Second*30, "experiment duration")
	payloadSize  = flag.Int("p", 1024, "payload size")
	keyRange     = flag.Int("k", 100000, "key range")
	readRate     = flag.Int("r", 0, "read percentage proportion")
	nThreads     = flag.Int("n", 1, "number of threads")
	thinkingTime = flag.Duration("t", time.Millisecond*100, "thinking time")
	logFrequency = flag.Uint64("f", 10, "log constant to determine frequency")
	bufferSize   = flag.Int("b", 2048, "response buffer size")
)

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().Unix())

	allSetWg := sync.WaitGroup{}
	clientsWg := sync.WaitGroup{}
	startCh := make(chan struct{})
	stopChan := make(chan struct{})

	tcpAddr, err := net.ResolveTCPAddr("tcp", *serverAddr)
	if err != nil {
		log.Fatal(err)
	}

	allSetWg.Add(*nThreads)
	clientsWg.Add(*nThreads)

	go func(addr net.TCPAddr, bsize int, psize int, rRate int, kRange int, logFreq uint64, tt time.Duration) {
		conn, err := net.DialTCP("tcp", nil, &addr)
		if err != nil {
			log.Fatal(err)
		}

		buffer := make([]byte, bsize)

		payload := make([]byte, psize)
		emptyPayload := make([]byte, 0)

		setReq := kv.Request{
			Op:   kv.SetOp,
			Data: payload,
		}
		getReq := kv.Request{
			Op:   kv.GetOp,
			Data: emptyPayload,
		}

		setReqBytes := setReq.Serialize()
		getReqBytes := getReq.Serialize()

		var reqBytes []byte

		allSetWg.Done()
		<-startCh

		for {
			select {
			case <-stopChan:
				conn.Close()
				clientsWg.Done()
				return
			default:
				randOpNumber := rand.Intn(100)
				if randOpNumber < rRate {
					reqBytes = getReqBytes
				} else {
					reqBytes = setReqBytes
				}

				key := uint64(rand.Intn(kRange))
				binary.PutUvarint(
					reqBytes[kv.OpByteSize:kv.OpByteSize+kv.KeyByteSize],
					key,
				)

				startTime := time.Now()

				conn.Write(reqBytes)

				_, err := conn.Read(buffer)
				if err != nil {
					log.Fatal(err)
				}

				if key%logFreq == 0 {
					fmt.Println(time.Now().UnixNano(), time.Since(startTime).Microseconds())
				}

				time.Sleep(tt)
			}
		}
	}(*tcpAddr, *bufferSize, *payloadSize, *readRate, *keyRange, *logFrequency, *thinkingTime)

	for i := 1; i <= *nThreads-1; i++ {
		go func(clientId int, addr net.TCPAddr, bsize int, psize int, rRate int, kRange int, logFreq uint64, tt time.Duration) {
			conn, err := net.DialTCP("tcp", nil, &addr)
			if err != nil {
				log.Fatal(err)
			}

			buffer := make([]byte, bsize)

			setReq := kv.Request{
				Op:   kv.SetOp,
				Data: make([]byte, psize),
			}
			getReq := kv.Request{
				Op:   kv.GetOp,
				Data: make([]byte, 0),
			}

			setReqBytes := setReq.Serialize()
			getReqBytes := getReq.Serialize()

			var reqBytes []byte

			allSetWg.Done()
			<-startCh

			for {
				select {
				case <-stopChan:
					conn.Close()
					clientsWg.Done()
					return
				default:
					randOpNumber := rand.Intn(100)
					if randOpNumber < rRate {
						reqBytes = getReqBytes
					} else {
						reqBytes = setReqBytes
					}

					key := uint64(rand.Intn(kRange))
					binary.PutUvarint(
						reqBytes[kv.OpByteSize:kv.OpByteSize+kv.KeyByteSize],
						key,
					)

					conn.Write(reqBytes)

					_, err := conn.Read(buffer)
					if err != nil {
						log.Fatal(err)
					}

					time.Sleep(tt)
				}
			}
		}(i, *tcpAddr, *bufferSize, *payloadSize, *readRate, *keyRange, *logFrequency, *thinkingTime)
	}

	allSetWg.Wait()

	close(startCh)

	timer := time.NewTimer(*duration)
	<-timer.C

	close(stopChan)

	clientsWg.Wait()
}
