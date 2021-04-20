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

	payload := make([]byte, *payloadSize)
	emptyPayload := make([]byte, 0)

	allSetWg.Add(*nThreads)
	clientsWg.Add(*nThreads)

	go func() {
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			log.Fatal(err)
		}

		buffer := make([]byte, *bufferSize)

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
				if randOpNumber < *readRate {
					reqBytes = getReqBytes
				} else {
					reqBytes = setReqBytes
				}

				key := uint64(rand.Intn(*keyRange))
				binary.PutUvarint(
					reqBytes[kv.OpByteSize:kv.OpByteSize+kv.KeyByteSize],
					key,
				)

				startTime := time.Now()

				conn.Write(reqBytes)
				conn.Write([]byte("\n"))

				_, err := conn.Read(buffer)
				if err != nil {
					log.Fatal(err)
				}

				if key%*logFrequency == 0 {
					fmt.Println(time.Now().UnixNano(), time.Since(startTime).Microseconds())
				}

				time.Sleep(*thinkingTime)
			}
		}
	}()

	for i := 1; i <= *nThreads-1; i++ {
		go func(clientId int) {
			conn, err := net.DialTCP("tcp", nil, tcpAddr)
			if err != nil {
				log.Fatal(err)
			}

			buffer := make([]byte, *bufferSize)

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
					if randOpNumber < *readRate {
						reqBytes = getReqBytes
					} else {
						reqBytes = setReqBytes
					}

					key := uint64(rand.Intn(*keyRange))
					binary.PutUvarint(
						reqBytes[kv.OpByteSize:kv.OpByteSize+kv.KeyByteSize],
						key,
					)

					conn.Write(reqBytes)
					conn.Write([]byte("\n"))

					_, err := conn.Read(buffer)
					if err != nil {
						log.Fatal(err)
					}

					time.Sleep(*thinkingTime)
				}
			}
		}(i)
	}

	allSetWg.Wait()

	close(startCh)

	timer := time.NewTimer(*duration)
	<-timer.C

	close(stopChan)

	clientsWg.Wait()
}
