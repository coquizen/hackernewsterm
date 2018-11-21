package main

import (
	"log"

	"github.com/caninodev/hackernewsterm/hackernews"
	"github.com/rivo/tview"
)

type myApp struct {
	main *tview.Application
	api  *hackernews.HAPI
	//pages *tview.Pages
	gui *GUI
}

var app *myApp

func (a *myApp) initialize() {
	a.main = tview.NewApplication()
	a.api = hackernews.NewHAPI(false, nil)
	a.gui = &GUI{}

	a.gui.Create()

	a.main.SetRoot(a.gui.layout, true)
	a.main.SetInputCapture(a.gui.KeyHandler)
}

// and finally, putting it all together
func main() {
	app = new(myApp)
	app.initialize()

	err := app.main.Run()
	if err != nil {
		log.Panicln(err)
	}
}
