package main

import (
	"flag"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/r3musketeers/hermes/pkg/kv"
)

var (
	serverAddr = flag.String("s", ":8000", "server address")
	duration   = flag.Duration("d", time.Second*10, "experiment duration")
	logPath    = flag.String("l", "latency.log", "path to log the latency")
)

func main() {
	flag.Parse()

	wg := sync.WaitGroup{}
	startCh := make(chan struct{})
	stopChan := make(chan struct{})

	wg.Add(1)
	go func() {
		tcpAddr, err := net.ResolveTCPAddr("tcp", *serverAddr)
		if err != nil {
			log.Fatal(err)
		}

		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			log.Fatal(err)
		}

		logFile, err := os.Create(*logPath)
		if err != nil {
			log.Fatal(err)
		}

		logger := log.New(logFile, "", log.LstdFlags)

		buffer := make([]byte, 2048)

		<-startCh

		for {
			select {
			case <-stopChan:
				wg.Done()
				return
			default:
				randKey := rand.Intn(1000000)
				req := kv.Request{
					Op:   kv.SetOp,
					Key:  uint64(randKey),
					Data: make([]byte, 1008),
				}

				startTime := time.Now()

				conn.Write(req.Serialize())

				_, err := conn.Read(buffer)
				if err != nil {
					log.Fatal(err)
				}

				if randKey%10 == 0 {
					logger.Println(time.Since(startTime).Microseconds())
				}

				time.Sleep(time.Millisecond * 100)
			}
		}
	}()

	close(startCh)

	timer := time.NewTimer(*duration)
	<-timer.C

	close(stopChan)
	wg.Wait()
}
