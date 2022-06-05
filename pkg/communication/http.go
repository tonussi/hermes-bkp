package communication

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/r3musketeers/hermes/pkg/proxy"
)

type HTTPCommunicator struct {
	fromAddr   string
	toAddr     string
	RequestURI string
	Method     string
}

func NewHTTPCommunicator(
	fromAddr string,
	toAddr string,
) (*HTTPCommunicator, error) {

	http.Get("http://" + toAddr + "/pulse")

	return &HTTPCommunicator{
		fromAddr:   fromAddr,
		toAddr:     toAddr,
		RequestURI: "/",
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
	var resp *http.Response
	var err error
	if comm.Method == "GET" {
		resp, err = http.Get(deliveryFullUrlString)
	} else {
		resp, err = http.Post(deliveryFullUrlString, "application/json", bufferedPayload)
	}
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
	log.Println(bodyString)
	log.Println(comm.RequestURI)

	return body, err
}

// Extra functions to clean code a little bit

func (comm *HTTPCommunicator) requestHandler(w http.ResponseWriter, r *http.Request, handle proxy.HandleIncomingMessageFunc) {
	comm.RequestURI = r.RequestURI
	comm.Method = r.Method
	bodyBytes, _ := ioutil.ReadAll(r.Body)

	log.Println("handling connection reading bytes and sending to handler")

	resp, err := handle(bodyBytes)

	if err != nil {
		log.Fatalln(err)
	}

	bodyString := string(resp)
	log.Print(bodyString)
	fmt.Fprintf(w, "%+v", bodyString)
	if err != nil {
		log.Fatalln(err)
	}
}

func (comm *HTTPCommunicator) buildHttpUrlPath() string {
	return "http://" + comm.toAddr + comm.RequestURI
}

func payloadBytesAsBufferedReader(data []byte) (ioBufferedValues *bytes.Buffer) {
	return bytes.NewBuffer(data)
}
