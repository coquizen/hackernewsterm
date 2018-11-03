// Adapted from https://github.com/easyCZ/grpc-web-hacker-news/blob/master/server/hackernews/api.go
package hackernews

import (
	"fmt"
	"log"
	"net/http"

	. "github.com/caninodev/hackernewsterm/models"
	"gopkg.in/zabawaba99/firego.v1"
)

const baseURL = "https://hacker-news.firebaseio.com"

const version = "/v0"

var endPoint = map[string]string{
	"top":  "/topstories",
	"new":  "/newstories",
	"best": "/beststories",
	"ask":  "/askstories",
	"jobs": "/jobstories",
	"show": "/showstories",
}

type hackerNewsApi struct {
	*firego.Firebase
}

// Service provides an interface to HN's FireBase
func NewHackerNewsAPI(client *http.Client) *hackerNewsApi {
	if client == nil {
		client = http.DefaultClient
	}

	fb := firego.New(baseURL, client)
	return &hackerNewsApi{
		Firebase: fb,
	}
}

func (api *hackerNewsApi) GetItem(id int) (*Item, error) {
	ref, err := api.Ref(fmt.Sprintf("/v0/item/%d", id))
	if err != nil {
		log.Fatalf("request story reference failed @ reference: %")
	}
	var value Item
	if err := ref.Value(&value); err != nil {
		log.Fatalf("story #%d retrieval failed %s", id, err)
	}

	return &Item{
		value.ID,
		value.Deleted,
		value.Type,
		value.By,
		value.Time,
		value.Text,
		value.Dead,
		value.Parent,
		value.Poll,
		value.Kids,
		value.URL,
		value.Score,
		value.Parts,
		value.Descendants,
	}, nil
}

func (api *hackerNewsApi) GetItems(reqType *string) chan *Item {

	items := make(chan *Item)
	ref, err := api.Firebase.Ref(version + endPoint[*reqType])
	if err != nil {
		log.Fatal("error firebase reference")
	}

	var ids []uint32
	if err := ref.Value(&ids); err != nil {
		log.Printf("%s stories request failed", reqType)
	}
	ids = ids[:5]
	log.Printf("IDs: %v", ids)
	for _, id := range ids {
		go func(id int) {
			item, _ := api.GetItem(id)
			items <- item
		}(int(id))
	}
	return items
}
