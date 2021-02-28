package communication

// import (
// 	"bufio"
// 	"bytes"
// 	"crypto/sha256"
// 	"net/http"
// 	"net/http/httputil"
// 	"net/url"

// 	"github.com/r3musketeers/hermes/proxy"
// )

// type HTTPCommunicator struct {
// 	deliverAddr   string
// 	clientWriters map[string]http.ResponseWriter
// 	clientRelease map[string]chan interface{}
// }

// func NewHTTPCommunicator(
// 	deliverAddr string,
// ) *HTTPCommunicator {
// 	return &HTTPCommunicator{
// 		deliverAddr:   deliverAddr,
// 		clientWriters: map[string]http.ResponseWriter{},
// 		clientRelease: map[string]chan interface{}{},
// 	}
// }

// func (comm HTTPCommunicator) Listen(
// 	addr string,
// 	handle proxy.HandleIncomingMessageFunc,
// ) error {
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		buffer := bytes.NewBuffer([]byte{})

// 		err := r.Write(buffer)
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		hash := sha256.New()
// 		_, err = hash.Write(buffer.Bytes())
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		id, err := handle(buffer.Bytes())
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		comm.clientWriters[id] = w
// 		comm.clientRelease[id] = make(chan interface{})

// 		<-comm.clientRelease[id]

// 		delete(comm.clientWriters, id)
// 		delete(comm.clientRelease, id)
// 	})

// 	return http.ListenAndServe(addr, nil)
// }

// func (comm HTTPCommunicator) Deliver(id string, data []byte) error {
// 	defer close(comm.clientRelease[id])

// 	req, _ := http.ReadRequest(bufio.NewReader(bytes.NewReader(data)))

// 	url := &url.URL{
// 		Scheme:     "http",
// 		Host:       comm.deliverAddr,
// 		Path:       req.URL.Path,
// 		ForceQuery: false,
// 	}

// 	req.URL = url
// 	req.Host = url.Host
// 	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))

// 	proxy := httputil.NewSingleHostReverseProxy(url)
// 	proxy.ServeHTTP(comm.clientWriters[id], req)

// 	return nil
// }
