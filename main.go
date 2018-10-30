package main

import (
	"log"
	"net/http"


	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"./hackernewsAPI"
)

// http ws to websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var hn = hackernews.NewHackerNewsAPI(http.DefaultClient)

// Define the server-client JSON format
type Message struct {
	RequestType string `json:"type"`
	Payload     string `json:"payload"`
}

// Create a websocket server
func createServer() {
	router := mux.NewRouter()
	router.HandleFunc("/", wsNewsHandler)
	log.Fatal(http.ListenAndServe(":8000", router))
}

func wsNewsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatalf("upgrade:", err)
		return
	}

	defer ws.Close()

	go handleData(ws, <-hn.GetTopStories())
}

func handleData(ws *websocket.Conn, story *hackernews.Story) error {
	for {
		// Read client request as JSON and map it to a Message object
		log.Print(story)
		err := ws.ReadJSON(&story)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	createServer()
}