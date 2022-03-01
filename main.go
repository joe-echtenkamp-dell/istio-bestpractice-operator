package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type Result struct {
	ServerTime time.Time
	ServerTZ   string

	ClientTime time.Time
	ClientTZ   string

	Pass bool
}

func main() {
	// re-used vars
	reqIdHeaderKey := http.CanonicalHeaderKey("x-request-id")
	client := &http.Client{}

	////////// TEST 1 - Ensure i cant talk directly to the time-server //////////////////////
	serverUuidWithHyphen := uuid.New()

	serverUrl := os.Getenv("SERVERURL")

	// setup request to proxy server
	req, err := http.NewRequest("GET", serverUrl, nil)
	if err != nil {
		log.Print("Failed to create new request.")
		log.Fatal(err.Error())
	}

	// add requestid
	req.Header.Set(reqIdHeaderKey, serverUuidWithHyphen.String())

	// send request to proxy
	resp, err := client.Do(req)
	if err != nil {
		// this is for actual failures, not non-2** codes
		log.Print("Failed to make request to server")
		log.Fatal(err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		// this is a problem, as we shouldnt be able to access the server directly
		log.Print("Can access the time-server, which should violate allow-nothing policy")
		log.Fatal(err.Error())
	}

	// defer resp.Body.Close()
	// var proxyResult Result
	// json.NewDecoder(resp.Body).Decode(&proxyResult)

	/////////////////////////////////////////////////////////////////////////////////////////
	////////// TEST 2 - Ensure i can talk to the proxy server ///////////////////////////////
	// Create request ID for this test
	uuidWithHyphen := uuid.New()

	proxyUrl := os.Getenv("PROXYURL")

	// setup request to proxy server
	req, err = http.NewRequest("GET", proxyUrl, nil)
	if err != nil {
		log.Print("Failed to create new request.")
		log.Fatal(err.Error())
	}

	// add requestid
	req.Header.Set(reqIdHeaderKey, uuidWithHyphen.String())

	// send request to proxy
	resp, err = client.Do(req)
	if err != nil {
		log.Print("Failed to make request to proxy")
		log.Fatal(err.Error())
	}

	defer resp.Body.Close()
	var proxyResult Result
	json.NewDecoder(resp.Body).Decode(&proxyResult)

	// check that the server returned the correct request id
	returnVal, proxyOk := resp.Header[reqIdHeaderKey]

	if !proxyOk {
		log.Print("proxy didnt return x-request-id")
		log.Fatal(err.Error())
	}

	if returnVal[0] != uuidWithHyphen.String() {
		log.Fatal("proxy returned different x-request-id")
	}
	/////////////////////////////////////////////////////////////////////////////////////////

	log.Print("Test complete.")
	os.Exit(0)
}
