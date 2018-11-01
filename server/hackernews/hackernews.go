// Adapted from https://github.com/easyCZ/grpc-web-hacker-news/blob/master/server/hackernews/api.go
package hackernews

import (
	"fmt"
	"github.com/caninodev/hackernewsterm/models"
	"log"
	"net/http"

	"gopkg.in/zabawaba99/firego.v1"
	_ "github.com/caninodev/hackernewsterm/models"
)

const baseURL = "https://hacker-news.firebaseio.com"

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

func (api *hackerNewsApi) GetStory(id int) (*Story, error) {
	ref, err := api.Ref(fmt.Sprintf("/v0/item/%d", id))
	if err != nil {
		log.Fatalf("request story reference failed @ reference: %")
	}
	var value models.Story
	if err := ref.Value(&value); err != nil {
		log.Fatalf("story #%d retrieval failed %d", id, err)
	}

	return &Story{
		By:    value.By,
		Id:    value.Id,
		Score: value.Score,
		Time:  value.Time,
		Type:  value.Type,
		Url:   value.Url,
	}, nil
}

func (api *hackerNewsApi) GetTopStories() chan *Story {
	stories := make(chan *Story)
	ref, err := api.Firebase.Ref("/v0/topstories")
	if err != nil {
		log.Fatal("error firebase reference")
	}

	var ids []uint32
	if err := ref.Value(&ids); err != nil {
		log.Fatalf("top stories request failed")
	}

	ids = ids[:10]
	for _, id := range ids {
		go func(id int) {
			story, _ := api.GetStory(id)
			stories <- story
		}(int(id))
	}
	return stories
}
