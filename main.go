package main

import (
	"log"

	"github.com/rivo/tview"

	"github.com/caninodev/hackernewsterm/hackernews"
	. "github.com/caninodev/hackernewsterm/models"
)

/*
Views holds reference for all views and renders them
*/
type Views struct {
	*Layout
	*AppState
}

// AppState connects the backend with the frontend via ContentChannel
type AppState struct {
	app    *tview.Application
	api    *hackernews.HAPI
	stream chan *Item
}

// Layout contains the root layout for the app
type Layout struct {
	*tview.Grid
	*tview.Flex
}

var (
	state          *AppState
	defaultRequest *Request
)

// ListView represents the item selector
type ListView struct {
	*tview.List
}

// ContentView is where user's selection is rendered
type ContentView struct {
	*tview.TextView
}

func (ui Views) createLayout() tview.Primitive {
	layout := new(Layout)
	(layout.Flex) = tview.NewFlex()

	list := &ListView{tview.NewList()}
	content := &ContentView{tview.NewTextView()}

	defaultRequest = &Request{
		RequestType: "top",
		Payload:     "20",
	}

	layout.Flex.AddItem(tview.Primitive(list), 0, 1, true)
	layout.Flex.AddItem(tview.Primitive(content), 0, 3, false)

	list.populate(defaultRequest)

	(ui.Layout) = layout

	return tview.Primitive(layout.Flex)
}

func (l ListView) populate(reqType *Request) tview.Primitive {
	l.ShowSecondaryText(true)
	l.SetBorder(true).SetTitle(reqType.RequestType + " stories")

	stream := state.api.GetItems(defaultRequest)

	go func() {
		go func() {
			for item := range stream {
				state.app.QueueUpdateDraw(func() {
					l.AddItem(item.Title, item.By, 0, nil)
				})
			}
		}()
	}()

	return l
}

func connectUI(app *tview.Application) error {
	state = &AppState{
		app: app,
		api: hackernews.NewHAPI(false, nil),
	}
	mainView := new(Views)

	app.SetRoot(mainView.createLayout(), true)

	return nil
}

// and finally, putting it all together
func main() {
	app := tview.NewApplication()
	err := connectUI(app)
	if err != nil {
		panic(err)
	}

	if err := app.Run(); err != nil {
		log.Panicln(err)
	}
}
