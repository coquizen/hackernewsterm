package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/caninodev/hackernewsterm/server/hackernews"

	. "github.com/caninodev/hackernewsterm/models"
	"github.com/gorilla/websocket"
)

// http ws to websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var hn = hackernews.NewHackerNewsAPI(http.DefaultClient)

var msg Request

func wsNewsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatalf("upgrade: %s", err)
		return
	}

	defer ws.Close()

	// Read client request as JSON and map it to a Story object
	for {
		if err := ws.ReadJSON(&msg); err != nil {
			log.Printf("read error: %s", err)
		}
		fmt.Printf("Got message: %#v\n", msg)

		items := <-hn.GetItems(&msg.RequestType)
		log.Printf("items: %#v\n", items)
		if err := ws.WriteJSON(&items); err != nil {
			log.Printf("%s is what went wrong", err)
		}
	}
}
func main() {
	quit := make(chan os.Signal) // HL
	signal.Notify(quit, os.Interrupt)
	srv := &http.Server{Addr: ":8000", Handler: http.DefaultServeMux}

	go func() { // HL
		<-quit // HL
		log.Println("Shutting down server...")
		if err := srv.Shutdown(context.Background()); err != nil { // HL
			log.Fatalf("could not shutdown: %v", err)
		}
	}()

	http.HandleFunc("/", wsNewsHandler)
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed { // HL
		log.Fatalf("listen: %s\n", err)
	}
}
