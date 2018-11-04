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
	"top":  version + "/topstories",
	"new":  version + "/newstories",
	"best": version + "/beststories",
	"ask":  version + "/askstories",
	"jobs": version + "/jobstories",
	"show": version + "/showstories",
}

type HAPI struct {
	*firego.Firebase
}

// NewHackerNewsAPI provides an interface to HN's FireBase
func NewHackerNewsAPI(hasHTTPClient bool, client *http.Client) *HAPI {
	if hasHTTPClient == true {
		if client == nil {
			client = http.DefaultClient
		}
	} else {
		client = nil
	}
	fb := firego.New(baseURL, client)
	return &HAPI{
		Firebase: fb,
	}
}

func (api *HAPI) GetItem(id int) (*Item, error) {
	ref, err := api.Ref(fmt.Sprintf("/item/%d", id))
	if err != nil {
		log.Fatalf("request story reference failed @ reference: %s", err)
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

func (api *HAPI) GetItems() (requestChan chan *Request, itemChan chan *Item) {
	requestChan = make(chan *Request)
	itemChan = make(chan *Item)

	go func() {
		for {
			requestType := <-requestChan
			ref, err := api.Firebase.Ref(endPoint[requestType])
			if err != nil {
				log.Fatal("error firebase reference")
			}
			var ids []int
			if err := ref.Value(&ids); err != nil {
				log.Printf("%s stories request failed", reqType)
			}
			for _, id := range ids {
				go func(id int) {
					item, err := api.GetItem(id)
					if err != nil {
						log.Printf("#%d error: %s", id, err)
					}
					itemChan <- item
				}(int(id))
			}
		}
	}()
	return requestChan, itemChan
}
