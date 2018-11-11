package main

import (
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"

	"github.com/gdamore/tcell"
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
	list    *ListView
	content *ContentView
}

var (
	state          *AppState
	defaultRequest *Request
	layout         *Layout
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
	arrangement := tview.NewFlex()

	list := &ListView{tview.NewList()}
	content := &ContentView{tview.NewTextView()}
	layout = &Layout{
		list:    list,
		content: content,
	}

	defaultRequest = &Request{
		RequestType: "top",
		Payload:     "20",
	}

	arrangement.AddItem(tview.Primitive(list), 0, 1, true)
	arrangement.AddItem(tview.Primitive(content), 0, 3, false)

	list.populate(defaultRequest)

	return tview.Primitive(arrangement)

}

func (l ListView) populate(reqType *Request) tview.Primitive {
	l.ShowSecondaryText(true)
	l.SetBorder(true).SetTitle(reqType.RequestType + " stories")

	stream := state.api.GetItems(defaultRequest)

	go func() {
		for item := range stream {
			state.app.QueueUpdateDraw(func() {
				l.AddItem(item.Title, item.By, 0, func() {
					layout.content.render(item)
				})
			})
		}
	}()
	return l
}

func (c ContentView) render(item *Item) {
	c.Clear()
	c.SetBorder(false).SetTitle(item.Title)
	_, _, numCols, _ := c.GetInnerRect()
	prsedNumCols := strconv.Itoa(numCols)

	if _, err := exec.LookPath("w3m"); err != nil {
		c.SetTextColor(tcell.ColorRed).SetText("Please install w3m for full functionality")
	} else {
		webCMD := exec.Command("w3m", "-dump", "-cols", "-X", prsedNumCols, item.URL.String())
		webOutPipe, webErr := webCMD.StdoutPipe()
		if webErr != nil {
			log.Print(webErr)
		}
		renderedPage, err := ioutil.ReadAll(webOutPipe)
		if err != nil {
			log.Print(err)
		}
		if _, err := c.Write(renderedPage); err != nil {
			log.Print(err)
		}
	}
	return
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
