package main

import (
	"github.com/rivo/tview"

	"github.com/caninodev/hackernewsterm/hackernews"
	. "github.com/caninodev/hackernewsterm/models"
)

var (
	app *tview.Application
)

var fb = hackernews.NewHAPI(false, nil)

type ui struct {
	list *tview.List
}

func main() {
	catalog := tview.NewList().ShowSecondaryText(true)


	requestChan, itemChan := fb.GetItems()

	req := &Request{
		"top",
		"10",
	}

	requestChan <- req

	defer close(itemChan)

	for item := range itemChan {
		catalog.AddItem(item.Title, item.By, 0, nil)
	}

	catalog.AddItem("Quit", "Select to Exit", 'q', func() {
		app.Stop()
	})
	if err := tview.NewApplication().SetRoot(catalog, true).Run(); err != nil {
		panic(err)
	}
}
