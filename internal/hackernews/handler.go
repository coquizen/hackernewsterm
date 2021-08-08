package hackernews

import (
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

type RequestType int

const (
	NewStories RequestType = iota
	BestStories
	TopStories
	AskStories
	ShowStories
	JobStories
	AStory
	MaxID
)

type Endpoints map[RequestType]string

const (
	BaseURI = "https://hacker-news.firebaseio.com/"
	Version = "v0"
)

type Handler interface {
	Subscribe(request RequestType) chan Item
	Close() <-chan error
}

//goland:noinspection ALL
var (
	endpoint = Endpoints{
		NewStories:  "newstories",
		BestStories: "beststories",
		TopStories:  "topstories",
		AskStories:  "askstories",
		ShowStories: "showstories",
		JobStories:  "jobstories",
		AStory:      "item",
		MaxID:       "maxitem",
	}
)

type Firebase struct {
	ctx context.Context
	fb  *db.Client
}

func NewHandlerWithDefaultConfig(ctx context.Context) *Firebase {
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
	svc := &Firebase{
		fb:  fb,
		ctx: ctx,
	}
	return svc
}

func (f *Firebase) Subscribe(request RequestType) chan Item {
	var err error
	var IDs []uint
	IDs, err = f.fetch(request)
	if err != nil {
		log.Fatalf("error retrieving IDs: %v", err)
	}

	var item Item
	items := make(chan Item, 5)
        go func() {
		for ID := range IDs {
			item, err = f.post(uint(ID))

			if err != nil {
				log.Fatalf("error:  %v", err)
			}
			log.Print(item.Title())
			items <- item
		}
              }()
	return items
}

func (f *Firebase) Close() <-chan error {
	errChan := make(chan error)
	return errChan
}

func (f *Firebase) fetch(request RequestType) ([]uint, error) {
	ref := f.fb.NewRef(fmt.Sprintf("/%s/%s", Version, endpoint[request]))
	if request == MaxID {
		var maxID uint
		err := ref.Get(f.ctx, &maxID)
		if err != nil {
			return nil, err
		}
		return []uint{maxID}, nil
	}
	var posts []uint
	err := ref.Get(context.Background(), &posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (f *Firebase) post(ID uint) (Item, error) {
	ref := f.fb.NewRef(fmt.Sprintf("%s/item/%d", Version, ID))
	var item item
	err := ref.Get(f.ctx, &item)
	if err != nil {
		log.Fatalf("ref get error: %v", err)
	}
	return item, nil
}
