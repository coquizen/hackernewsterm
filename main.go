package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"

	"github.com/rivo/tview"

	"github.com/caninodev/hackernewsterm/hackernews"
	. "github.com/caninodev/hackernewsterm/models"
)

/*
UI holds reference for all views and renders them
*/
type UI struct {
	*GUI
	Populate func(request Request)

	layout *tview.Flex
	list *tview.List
	content *tview.TextView
	comments *tview.TreeView

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

//
//type CommentsView struct {
//	*tview.TreeView
//}

func (ui *UI) Create() {
	ui.list = tview.NewList()
	ui.list.ShowSecondaryText(true)

	ui.content = tview.NewTextview()
	ui.content.SetDynamicColors(true).
	 SetBorder(true)

   ui.comments = tview.NewTreeView()

	ui.layout= tview.NewFlex()
	ui.layout.SetDirection(tview.FlexColumn)
	ui.layout AddItem(ui.list, 0, 1, true)
	ui.layout.Additem((tview.NewFlex().
	SetDirection(tview.FlexRow).
	AddItem(ui.content, 0, 1, false).
	AddItem(ui.comments, 0, 1, false)), 0, 3, false)


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

func (ui *UI) Populate(reqType *Request) tview.Primitive {
	l.ShowSecondaryText(true)
	l.SetBorder(true).SetTitle(reqType.RequestType + " stories")

	stream := state.api.GetItems(defaultRequest)

	state.app.QueueUpdateDraw(func() {
		go func() {
			for item := range stream {
				l.AddItem(item.Title, item.By, 0, func() {
					layout.content.render(item)
				})
			}
		}()
	})
	return l
}

func (c ContentView) render(item *Item) {
	c.SetBorder(false).SetTitle(item.URL)
	c.SetDynamicColors(true)
	c.SetChangedFunc(func() {
		state.app.Draw()
	})

	_, _, numCols, _ := c.GetInnerRect()
	parsedNumCols := strconv.Itoa(numCols)

	webCMD := exec.Command("w3m", "-dump", "-graph", "-X", "-cols", parsedNumCols, item.URL)

	outr, _ := webCMD.CombinedOutput()
	fmt.Fprint(c, string(outr))
}

func connectUI(app *tview.Application) error {
	state = &AppState{
		app: app,
		api: hackernews.NewHAPI(false, nil),
	}
	mainView := new(UI)

	app.SetRoot(mainView.initUI(), true)

	return nil
}

// and finally, putting it all together
func main() {
	app := tview.NewApplication()
	err := connectUI(app)
	state.app.Draw()

	if err != nil {
		panic(err)
	}

	if err := app.Run(); err != nil {
		log.Panicln(err)
	}
}
