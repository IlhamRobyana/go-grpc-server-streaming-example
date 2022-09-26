package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/pramonow/go-grpc-server-streaming-example/utils"
	"golang.org/x/net/websocket"
)

var useWebsockets = flag.Bool("websockets", false, "Whether to use websockets")

type Message struct {
	Id      int    `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
	Amount  int    `json:"amount,omitempty"`
	Price   int    `json:"price,omitempty"`
}

// Client.
func main() {
	flag.Parse()

	if *useWebsockets {
		defer utils.TimeTrack(time.Now(), "websocket streaming")
		ws, err := websocket.Dial("ws://localhost:8080/", "", "http://localhost:8080")
		if err != nil {
			log.Fatal(err)
		}
		for {
			var m Message
			err = websocket.JSON.Receive(ws, &m)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}
			log.Printf("Received: %+v", m)
		}
	} else {
		defer utils.TimeTrack(time.Now(), "http streaming")
		log.Print("Sending request...")
		req, err := http.NewRequest("GET", "http://localhost:8080", nil)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Status code is not OK: %v (%s)", resp.StatusCode, resp.Status)
		}

		dec := json.NewDecoder(resp.Body)
		done := make(chan bool)

		go func() {
			for {
				var m Message
				err := dec.Decode(&m)
				if err != nil {
					if err == io.EOF {
						done <- true //close(done)
						break
					}
					log.Fatal(err)
				}
				log.Printf("Got response: %+v", m)
			}
		}()

		<-done
	}

	log.Printf("Server finished request...")
}
