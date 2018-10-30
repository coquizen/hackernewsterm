// Adapted from https://github.com/easyCZ/grpc-web-hacker-news/blob/master/server/hackernews/api.go

// Reference to HackerNewsAPI
// const endPoints = {
// 	topNews: "topnews.json",
// 	user: "user/",
// 	maxItem: "maxitem.json",
// 	askStories: "askstories.json",
// 	showStories: "showStories.json",
// 	jobStories: "jobStories.json"
// }

package hackernews

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gopkg.in/zabawaba99/firego.v1"
)

const baseURL = "https://hacker-news.firebaseio.com"
const version = "V0"

type hackerNewsApi struct {
	*firego.Firebase
}

// Story is the posting data type
type Story struct {
	By          string    `json:"by,omitempty"`
	Descendants int       `json:"descendants"`
	ID          int       `json:"id,omitempty"`
	Kids        []int     `json:"kids,omitempty"`
	Score       int       `json:"score,omitempty"`
	Time        time.Time `json:"time,omitempty"`
	Type        string    `json:"type,omitempty"`
	URL         string    `json:"url,omitempty"`
}

// Service provides an interface to HN's FireBase
func NewHackerNewsAPI(client *http.Client) *hackerNewsApi {
	if client == nil {
		client = http.DefaultClient
	}

	fb := firego.New(baseURL+"/"+version, client)
	return &hackerNewsApi{
		Firebase: fb,
	}
}

func (api *hackerNewsApi) GetStory(id int) (*Story, error) {
	ref, err := api.Ref(fmt.Sprintf("/item/%d", id))
	if err != nil {
		log.Fatalf("request story reference failed @ reference: %")
	}
	var value Story
	if err := ref.Value(&value); err != nil {
		log.Fatalf("story #%d retrieval failed %d", id, err)
	}

	return &Story{
		By:          value.By,
		Descendants: value.Descendants,
		ID:          value.ID,
		Kids:        value.Kids,
		Score:       value.Score,
		Time:        value.Time,
		Type:        value.Type,
		URL:         value.URL,
	}, nil
}

func (api *hackerNewsApi) GetTopStories() (chan *Story) {
	stories := make(chan *Story)
	ref, err := api.Firebase.Ref("/topstories")
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

