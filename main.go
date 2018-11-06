package main

import (
	"github.com/rivo/tview"

	"github.com/caninodev/hackernewsterm/hackernews"
	. "github.com/caninodev/hackernewsterm/models"
)

var (
	app            *tview.Application
	UI             *mainWindow
	defaultRequest *Request
	ItemList       *HNList
)

type mainWindow struct {
	flex *tview.Flex
}

type HNList struct {
	*tview.List
}

var fb = hackernews.NewHAPI(false, nil)

func main() {
	app = tview.NewApplication()
	defaultRequest = &Request{
		RequestType: "top",
		Payload: "10",
	}

	list := tview.NewList().ShowSecondaryText(true)
	list.SetBorder(true).SetTitle(defaultRequest.RequestType+" stories")
	itemChan := fb.GetItems(*defaultRequest)

	defer close(itemChan)
	go func() {
		for item := range itemChan {
			app.QueueUpdateDraw(func () {
				list.AddItem(item.Title, item.By, 0, nil)
			})
		}
	}()
	if err := app.SetRoot(list, true).Run(); err != nil {
		panic(err)
	}

}
