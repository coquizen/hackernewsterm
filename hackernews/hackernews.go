// Adapted from https://github.com/easyCZ/grpc-web-hacker-news/blob/master/server/hackernews/api.go
package hackernews

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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

// NewHAPI provides an interface to HN's FireBase
func NewHAPI(hasHTTPClient bool, client *http.Client) *HAPI {
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

func (api *HAPI) GetItem(id int) (Item, error) {
	ref, err := api.Ref(fmt.Sprintf(version + "/item/%d", id))
	if err != nil {
		log.Fatalf("request story reference failed @ reference: %s", err)
	}
	var value Item
	if err := ref.Value(&value); err != nil {
		log.Fatalf("story #%d retrieval failed %s", id, err)
	}

	return Item{
		ID:         value.ID,
		Deleted:    value.Deleted,
		Type:       value.Type,
		By:         value.By,
		Time:       value.Time,
		Text:       value.Text,
		Title:	   value.Title,
		Dead:       value.Dead,
		Parent:     value.Parent,
		Poll:       value.Poll,
		Kids:       value.Kids,
		URL:		value.URL,
		Score:      value.Score,
		Parts:  	value.Parts,
		Descendants: value.Descendants,
	}, nil
}

// GetItems is a aggregate function for the top-level endpoints as specified
// above.
func (api *HAPI) GetItems(req Request) (itemChan chan Item) {
	itemChan = make(chan Item)

	go func() {
		for {
			ref, err := api.Firebase.Ref(endPoint[req.RequestType])
			if err != nil {
				log.Fatal("error firebase reference")
			}
			var ids []int
			if err := ref.Value(&ids); err != nil {
				log.Printf("%s stories request failed", req.RequestType)
			}
			iterator, _ := strconv.Atoi(req.Payload)
			ids = ids[:iterator]
			for _, id := range ids {
				go func(id int) {
					item, _ := api.GetItem(id)
					itemChan <- item

				}(int(id))
			}
		}
	}()
	return
}



