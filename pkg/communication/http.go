package communication

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/r3musketeers/hermes/pkg/proxy"
)

type HTTPCommunicator struct {
	fromAddr      string
	toAddr        string
	urlPath       string
	method        string
	requestURI    string
	httpTextBytes []byte
	bodyBytes     []byte
	r             *http.Request
}

func NewHTTPCommunicator(
	fromAddr string,
	toAddr string,
	connAttempts int,
	connAttemptPeriod time.Duration,
) (*HTTPCommunicator, error) {

	http.Get("http://" + toAddr + "/pulse")

	return &HTTPCommunicator{
		fromAddr: fromAddr,
		toAddr:   toAddr,
		urlPath:  "/",
	}, nil
}

func (comm *HTTPCommunicator) Listen(handle proxy.HandleIncomingMessageFunc) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { comm.requestHandler(w, r, handle) })

	err := http.ListenAndServe(comm.fromAddr, nil)

	return err
}

func (comm *HTTPCommunicator) Deliver(data []byte) ([]byte, error) {
	if comm.httpTextBytes == nil {
		return nil, nil
	}

	client := &http.Client{}
	var buf bytes.Buffer
	var res *http.Response
	var req *http.Request
	var err error

	switch comm.method {
	case "GET":
		req, err = http.NewRequest("GET", "http://"+comm.toAddr+comm.r.RequestURI, nil)
		if err != nil {
			panic(err)
		}

		res, err = client.Do(req)
		if err != nil {
			panic(err)
		}

		defer res.Body.Close()

		if err = res.Write(&buf); err != nil {
			panic(err)
		}

		fmt.Println(buf.String())
	case "POST":
		bodyIoReader := payloadBytesAsBufferedReader(comm.bodyBytes)

		if err != nil {
			panic(err)
		}

		req, err = http.NewRequest("POST", "http://"+comm.toAddr+comm.r.RequestURI, bodyIoReader)
		if err != nil {
			panic(err)
		}

		res, err = client.Do(req)
		if err != nil {
			panic(err)
		}

		defer res.Body.Close()

		if err = res.Write(&buf); err != nil {
			panic(err)
		}

		fmt.Println(buf.String())
	}

	return buf.Bytes(), err
}

func (comm *HTTPCommunicator) requestHandler(w http.ResponseWriter, r *http.Request, handle proxy.HandleIncomingMessageFunc) {
	comm.r = r
	comm.method = r.Method
	comm.urlPath = r.URL.Path
	comm.requestURI = r.RequestURI
	var httpTextBytes bytes.Buffer

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	comm.bodyBytes = bodyBytes

	// https://stackoverflow.com/a/69055473
	r.Write(&httpTextBytes)

	comm.httpTextBytes = httpTextBytes.Bytes()
	fmt.Println(httpTextBytes.String())

	resp, err := handle(comm.httpTextBytes)

	if err != nil {
		panic(err)
	}

	bodyResponseFromAppServer := string(resp)
	log.Print(bodyResponseFromAppServer)
	fmt.Fprintf(w, "%+v", resp)

	if err != nil {
		panic(err)
	}
}

func payloadBytesAsBufferedReader(data []byte) (ioBufferedValues *bytes.Buffer) {
	return bytes.NewBuffer(data)
}
