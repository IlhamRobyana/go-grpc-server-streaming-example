package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/websocket"
)

type Message struct {
	Id      int    `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
	Amount  int    `json:"amount,omitempty"`
	Price   int    `json:"price,omitempty"`
}

// Heavily based on Kubernetes' (https://github.com/GoogleCloudPlatform/kubernetes) detection code.
var connectionUpgradeRegex = regexp.MustCompile("(^|.*,\\s*)upgrade($|\\s*,)")

func isWebsocketRequest(req *http.Request) bool {
	return connectionUpgradeRegex.MatchString(strings.ToLower(req.Header.Get("Connection"))) && strings.ToLower(req.Header.Get("Upgrade")) == "websocket"
}

func Handle(w http.ResponseWriter, r *http.Request) {
	// Handle websockets if specified.
	if isWebsocketRequest(r) {
		websocket.Handler(HandleWebSockets).ServeHTTP(w, r)
	} else {
		HandleHttp(w, r)
	}
	log.Print("Finished sending response...")
}

func HandleWebSockets(ws *websocket.Conn) {
	for i := 0; i < 1000000; i++ {
		log.Printf("Sending some data: %d", i)
		m := Message{
			Id:      i,
			Message: fmt.Sprintf("Sending you \"%d\"", i),
		}
		err := websocket.JSON.Send(ws, &m)
		if err != nil {
			log.Printf("Client stopped listening...")
			return
		}

		// Artificially induce a 1s pause
		// time.Sleep(time.Second)
	}
}

func HandleHttp(w http.ResponseWriter, r *http.Request) {
	cn, ok := w.(http.CloseNotifier)
	if !ok {
		http.NotFound(w, r)
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.NotFound(w, r)
		return
	}

	// Send the initial headers saying we're gonna stream the response.
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	enc := json.NewEncoder(w)

	for i := 0; i < 1000000; i++ {
		select {
		case <-cn.CloseNotify():
			log.Printf("Client stopped listening")
			return
		default:
			// Artificially wait a second between reponses.
			// time.Sleep(time.Second)

			log.Printf("Sending some data: %d", i)
			m := Message{
				Id:      i,
				Message: fmt.Sprintf("Sending you \"%d\"", i),
				Name:    "John Doe",
				Address: "Bandung",
				Amount:  100,
				Price:   90290000,
			}

			// Send some data.
			err := enc.Encode(m)
			if err != nil {
				log.Fatal(err)
			}
			flusher.Flush()
		}
	}
}

// Server.
func main() {
	flag.Parse()

	http.HandleFunc("/", Handle)

	log.Printf("Serving...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
