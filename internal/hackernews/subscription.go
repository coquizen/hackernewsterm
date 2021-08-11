package hackernews

import (
	"context"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

type RequestType int

const (
	NewStories RequestType = iota
	BestStories
	TopStories
	AskStories
	ShowStories
	JobStories
	MaxID
)

type Command int

const (
	Pause Command = iota
	Play
	Stop
)

type Fetcher interface {
	Fetch() ([]uint, error)
}
type Subscription interface {
	Updates() <-chan Item
	Command() chan<- Command
	Close() error
}

type subscription struct {
	fetcher   Fetcher
	updates   chan Item
	command   chan Command
	closing   chan chan error
	itemCount int
	request   RequestType
	store     Store
}

// Subscribe returns the Subscriber interface which provides for a channel of items and a way to close out said
// channel
func Subscribe(fetcher Fetcher) Subscription {
	ctx := context.Background()
	expiryDuration := 5 * time.Minute
	firebaseClient := NewFirebaseClientWithDefaultConfig(ctx)
	hackerNewsCache := cache.New(expiryDuration, expiryDuration)
	store := NewCachedStore(hackerNewsCache, firebaseClient)
	s := &subscription{
		fetcher: fetcher,
		command: make(chan Command),
		updates: make(chan Item),
		store:   store,
	}
	go s.loop()
	s.Command() <- Play
	return s
}

func (s *subscription) loop() {
	defer close(s.updates)
	var err error
	var wg sync.WaitGroup
	var item Item
	var IDs []uint

	reset := func() {
		IDs, err = s.fetcher.Fetch()
		if err != nil {
			<-s.closing
			return
		}
	}
	reset()

	for _, ID := range IDs {
		wg.Add(1)
		item, err = s.store.Item(ID)
		if err != nil {
			<-s.closing
			return
		}
		select {
		case cmd := <-s.command:
			switch cmd {
				case Stop:
					<-s.closing
					return
				case Play:
					s.updates <- item
				case Pause:
					continue
				default:
					continue
				}
		case errCh := <-s.closing:
				errCh <- err
				close(s.updates)
				return
			}
			wg.Done()
		}

}

func (s *subscription) Updates() <-chan Item {
	return s.updates
}

func (s *subscription) Command() chan<- Command {
	return s.command
}
func (s *subscription) Close() error {
	errCh := make(chan error)
	s.closing <- errCh
	return <-errCh
}
