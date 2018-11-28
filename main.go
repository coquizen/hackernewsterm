package main

import (
	"log"

	"github.com/caninodev/hackernewsterm/hnapi"
	"github.com/rivo/tview"
)

type myApp struct {
	main *tview.Application
	api  *hnapi.HNdb
	//pages *tview.Pages
	gui *GUI
}

var app *myApp

func (a *myApp) initialize() {
	a.main = tview.NewApplication()
	a.api = hnapi.New()
	a.gui = &GUI{}

	a.gui.Create()

	a.main.SetRoot(a.gui.layout, true)
	a.main.SetInputCapture(a.gui.keyHandler)
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
