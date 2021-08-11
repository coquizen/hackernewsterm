package main

import (
	"context"
	"github.com/CaninoDev/hackernewsterm/internal/hackernews"
	"github.com/CaninoDev/hackernewsterm/internal/ui"
	"log"
)

// type App struct {
// 	ctx  *context.Context
// 	UI   *ui.UI
// 	HNdb *api.HackerNewsFB
// }

func main() {
	ctx := context.Background()
	handler := hackernews.NewFirebaseClientWithDefaultConfig(ctx)
        //var items []hackernews.Item


// 	for i := 0; i < 10; i++ {
// 		item := <-handler.Subscribe(hackernews.NewStories)
//                 items = append(items, item)
// 		log.Print(item.Title())
// 
// 	}
	if err := ui.InitUI(ctx, *handler); err != nil {
		log.Fatalf("error executing UI: %v", err)
	}


}
