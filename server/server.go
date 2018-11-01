package server

import (
	"github.com/caninodev/hackernewsterm/models"
	"github.com/caninodev/hackernewsterm/server/hackernews"
	"log"
	"net/http"

	_ "github.com/caninodev/hackernewsterm/models"
	"github.com/gorilla/websocket"
)

// http ws to websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var hn = hackernews.NewHackerNewsAPI(http.DefaultClient)

func wsNewsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatalf("upgrade: %s", err)
		return
	}

	defer ws.Close()

	go handleData(ws)
}

func handleData(ws *websocket.Conn) error {
	for {
		msg := Message{}
		req := ws.ReadJSON(&msg)
		// Read client request as JSON and map it to a Story object
		stories := <-hn.getStories(&req)
		err := ws.WriteJSON(&stories)
		if err != nil {
			log.Fatalf("%s is an invalid request", reqType)
		}
	}
	return nil
}

func createServer() {
	http.HandleFunc("/", wsNewsHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
