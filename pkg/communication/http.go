package communication

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/tonussi/studygo/pkg/proxy"
)

type HTTPCommunicator struct {
	fromAddr string
	toAddr   string
	urlPath  string
}

func NewHTTPCommunicator(
	fromAddr string,
	toAddr string,
	connAttempts int,
	connAttemptPeriod time.Duration,
) (*HTTPCommunicator, error) {
	// check readliness
	// to see if its responding
	resp, err := http.Get("http://" + toAddr)
	for connAttempts > 0 {
		if resp.StatusCode == 200 {
			break
		}
		resp, err := http.Get("http://" + toAddr)
		log.Println("connection attempts left:", connAttempts)
		log.Println("connection attempts left:", resp)
		if resp.StatusCode != 200 || err != nil {
			connAttempts--
			time.Sleep(connAttemptPeriod)
		}
		if connAttempts == 0 {
			return nil, err
		}
	}
	if resp.StatusCode != 200 || err != nil {
		return nil, err
	}

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
	// build url to post
	deliveryFullUrlString := comm.buildHttpUrlPath()

	// payload bytes as buffered reader
	bufferedPayload := payloadBytesAsBufferedReader(data)

	// delivery to a server
	resp, err := http.Post(deliveryFullUrlString, "application/json", bufferedPayload)
	if err != nil {
		log.Fatalln(err)
	}

	// close response body
	defer resp.Body.Close()

	// read body response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return nil, err
	}

	// see data that has been returned to the client
	bodyString := string(body)
	log.Print(bodyString)

	return body, err
}

// Extra functions to clean code a little bit

func (comm *HTTPCommunicator) requestHandler(w http.ResponseWriter, r *http.Request, handle proxy.HandleIncomingMessageFunc) {
	comm.urlPath = r.URL.Path
	bodyBytes, _ := ioutil.ReadAll(r.Body)

	log.Println("handling connection reading bytes and sending to handler")

	resp, err := handle(bodyBytes)

	if err != nil {
		log.Fatalln(err)
	}

	bodyString := string(resp)
	log.Print(bodyString)
	fmt.Fprintf(w, "%+v", resp)
	if err != nil {
		log.Fatalln(err)
	}
}

func (comm *HTTPCommunicator) buildHttpUrlPath() string {
	return "http://" + comm.toAddr + comm.urlPath
}

func payloadBytesAsBufferedReader(data []byte) (ioBufferedValues *bytes.Buffer) {
	return bytes.NewBuffer(data)
}
