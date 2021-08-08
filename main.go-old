package main

import (
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

func (app *myApp) initialize() {
	app.main = tview.NewApplication()
	app.api = hnapi.New()
	app.gui = &GUI{}

	app.gui.Create()

	app.main.SetInputCapture(app.gui.keyHandler)

}

// and finally, putting it all together
func main() {
	app = new(myApp)
	app.initialize()
	// Start the application.
	if err := app.main.SetRoot(app.gui.layout, true).Run(); err != nil {
		panic(err)
	}
}
