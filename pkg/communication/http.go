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
	fromAddr string
	toAddr   string

	httpTextBytes []byte
	bodyBytes     []byte
	r             *http.Request
	client        *http.Client
}

func NewHTTPCommunicator(
	fromAddr string,
	toAddr string,
) (*HTTPCommunicator, error) {
	// Creates a client to be used inside Deliver
	// This client will act as a http requester
	client := &http.Client{}

	// Constructor of the HTTPCommunicator class
	return &HTTPCommunicator{
		fromAddr: fromAddr,
		toAddr:   toAddr,
		client:   client,
	}, nil
}

func (comm *HTTPCommunicator) Listen(handle proxy.HandleIncomingMessageFunc) error {
	// Set a handler function for the http[s]://www.{ip}:{port}/
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { comm.InterceptClientRequest(w, r, handle) })

	// Start to server the http server to receive the requests
	err := http.ListenAndServe(comm.fromAddr, nil)

	// Return errors of the listen and server
	return err
}

func (comm *HTTPCommunicator) Deliver(data []byte) ([]byte, error) {
	// Data that enters Deliver is garanteed to be ordered already

	// TODO(developer): investigate bug of Hermes when Hermes starts
	if comm.httpTextBytes == nil {
		return nil, nil
	}

	var buf bytes.Buffer
	var res *http.Response
	var req *http.Request
	var err error

	// Restore the request that was previous formatted as bytes by the InterceptClientRequest function
	bodyIoReader := bytes.NewBuffer(comm.bodyBytes)

	// Build the request using the client
	req, err = http.NewRequest(comm.r.Method, "http://"+comm.toAddr+comm.r.RequestURI, bodyIoReader)

	// Add the Header to the request
	req.Header = comm.r.Header
	if err != nil {
		log.Fatal(err.Error())
	}

	// Do the request to the Server
	// If something wrong happens here its important to treat
	// Because the server did not receive the message
	// All the other replicas received
	res, err = comm.client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	// defer statement is executed, and places
	// res.Body.Close() on a list to be executed
	// prior to the function returning
	defer res.Body.Close()

	// Write the response of the Stateful Service to Hermes
	// SS -> H
	res.Write(&buf)

	// Return the bytes
	return buf.Bytes(), err

	// res.Body.Close() is now invoked
}

func (comm *HTTPCommunicator) InterceptClientRequest(w http.ResponseWriter, r *http.Request, handle proxy.HandleIncomingMessageFunc) {
	// The request that enters here
	// Is not ordered yet
	// And is a request pointed to Hermes, not to the Stateful Service
	// Save client request
	comm.r = r

	// Put on a list to be executed after the return
	defer r.Body.Close()

	// Save body in byte array format
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	comm.bodyBytes = bodyBytes

	// Save HTTP text request in bytes format
	var httpTextBytes bytes.Buffer
	r.Write(&httpTextBytes)
	comm.httpTextBytes = httpTextBytes.Bytes()

	// Send bytes to be ordered
	resp, err := handle(comm.httpTextBytes)
	// Receive ordered bytes
	// This error is important to be treated
	// Errors here means that the ordering failed
	if err != nil {
		log.Fatal(err.Error())
	}

	// Transform the bytes that were sent
	bodyResponseFromAppServer := string(resp)

	// Use the http writter to answer to the client
	// Can be omitted if you don't want Hermes answering back
	fmt.Fprintf(w, "%s", bodyResponseFromAppServer)
}
