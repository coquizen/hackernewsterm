/* HNdb is intended to provide a stripped down interface customized to use specifically with Hacker News firbase.
It has been adapted from https://github.com/easyCZ/grpc-web-hacker-news/blob/master/server/hackernews/api.go.
*/
package hnapi

import (
	"fmt"
	"log"
	"sync"

	"gopkg.in/zabawaba99/firego.v1"
)

const baseURL = "https://hacker-news.firebaseio.com"

const version = "/v0"

var wg sync.WaitGroup

var endPoint = map[string]string{
	"top":  "/v0/topstories",
	"new":  "/v0/newstories",
	"best": "/v0/beststories",
	"ask":  "/v0/askstories",
	"jobs": "/v0/jobstories",
	"show": "/v0/showstories",
}

// HNdb has an embedded struct for the firebase interface
type HNdb struct {
	*firego.Firebase
}

// New establishes an API to Hacker New's Firebase.
func New() *HNdb {
	return &HNdb{
		Firebase: firego.New(baseURL, nil),
	}
}

// GetItem retrieves the specified item and parses it.
func (db *HNdb) GetItem(id int) (*Item, error) {
	ref, err := db.Ref(fmt.Sprintf(version+"/item/%d", id))
	if err != nil {
		log.Fatalf("request story reference failed @ reference: %s", err)
	}
	var value Item
	if err := ref.Value(&value); err != nil {
		log.Fatalf("story #%d retrieval failed %s", id, err)
	}

	return &Item{
		ID:          value.ID,
		Deleted:     value.Deleted,
		Type:        value.Type,
		By:          value.By,
		Time:        value.Time,
		Text:        value.Text,
		Title:       value.Title,
		Dead:        value.Dead,
		Parent:      value.Parent,
		Poll:        value.Poll,
		Kids:        value.Kids,
		URL:         value.URL,
		Score:       value.Score,
		Parts:       value.Parts,
		Descendants: value.Descendants,
	}, nil
}

// GetPosts retrieves the specified type and number of posts.
func (db *HNdb) GetPosts(req *Request) (contentChan chan *Item) {
	contentChan = make(chan *Item, req.NumPosts)

	for {
		ref, err := db.Firebase.Ref(endPoint[req.PostType])
		if err != nil {
			log.Fatal("error firebase reference")
		}
		var ids []int
		if err := ref.Value(&ids); err != nil {
			log.Printf("%s stories request failed", req.PostType)
		}

		ids = ids[:req.NumPosts]
		for _, id := range ids {
			item, _ := db.GetItem(id)
			contentChan <- item
		}
		close(contentChan)
		return contentChan
	}
}
