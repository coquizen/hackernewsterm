package hackernews

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"fmt"
	"google.golang.org/api/option"
	"log"
)
type Request map[RequestType]int

var endpointCount = Request{
	NewStories: 500,
	TopStories: 500,
	BestStories: 500,
	AskStories: 200,
	ShowStories: 200,
	JobStories: 200,
}

type Endpoints map[RequestType]string
//goland:noinspection ALL
var (
	endpoint = Endpoints{
		NewStories:  "newstories",
		BestStories: "beststories",
		TopStories:  "topstories",
		AskStories:  "askstories",
		ShowStories: "showstories",
		JobStories:  "jobstories",
		MaxID:       "maxitem",
	}
)

const (
	BaseURI = "https://hacker-news.firebaseio.com/"
	Version = "v0"
)

type Handler interface {
	FetchRequestedIDs(RequestType) ([]uint, error)
}
type FirebaseClient struct {
	ctx context.Context
	fb  *db.Client
}

// NewFirebaseClientWithDefaultConfig returns an instance of the firebase client without
// user authentication
func NewFirebaseClientWithDefaultConfig(ctx context.Context) *FirebaseClient {
	defaultCfg := &firebase.Config{
		DatabaseURL: BaseURI,
	}
	app, err := firebase.NewApp(ctx, defaultCfg, option.WithoutAuthentication())
	if err != nil {
		log.Fatalf("error intiializing firebase app: %v", err)
	}

	fb, err := app.Database(ctx)
	if err != nil {
		log.Fatalf("error intiializing firebase connection: %v", err)
	}
	fbClient := &FirebaseClient{
		fb:  fb,
		ctx: ctx,
	}
	return fbClient
}

func (f *FirebaseClient) FetchRequestedIDs(request RequestType) ([]uint, error) {
	ref := f.fb.NewRef(fmt.Sprintf("%s/%s", Version, endpoint[request]))
	if request == MaxID {
		var maxID uint
		err := ref.Get(f.ctx, &maxID)
		if err != nil {
			return nil, err
		}
		return []uint{maxID}, nil
	}
	var IDs []uint
	err := ref.Get(f.ctx, &IDs)
	if err != nil {
		return nil, err
	}
	return IDs, nil
}

func (f *FirebaseClient) Item(id uint) (Item, error) {
	ref := f.fb.NewRef(fmt.Sprintf("%s/item/%d", Version, id))
	var post item
	err := ref.Get(f.ctx, &post)
	if err != nil {
		return nil, err
	}
	return post, err
}

