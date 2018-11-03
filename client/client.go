package main

import (
	"log"
	"net/url"

	. "github.com/caninodev/hackernewsterm/models"
	"github.com/gorilla/websocket"
)

var item Item
var items []Item
func main() {
	uAddr := url.URL{Scheme: "ws", Host: "localhost:8000", Path: "/"}
	log.Printf("connecting to %s:", uAddr.String())

	ws, _, err := websocket.DefaultDialer.Dial(uAddr.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	top := Request{RequestType: "top", Message: nil}

	if err := ws.WriteJSON(&top); err != nil {
		log.Printf("error writing to socket: %s", err)
	}

	defer ws.Close()

	for {
		if err := ws.ReadJSON(&item); err != nil {
			log.Printf("error reading JSON: %s", err)
		}
		items = append(items, item)
		if err := ws.ReadJSON(&item); err != nil {
			log.Printf("error reading JSON: %s", err)
		}
		items = append(items, item)
		//_, p, err := ws.ReadMessage()
		//if err != nil {
		//	log.Panic(err)
		//}
		//if err := json.Unmarshal(p, &item); err != nil {
		//	log.Panic(err)
		//}
		log.Printf("client receieved %#v\n", item)
	}
}
