package main

import (
	"log"
	"os/exec"

	"github.com/caninodev/hackernewsterm/hnapi"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// WebContent is a page that will render the selected item's url
func WebContent(nextSlide func()) (title string, content tview.Primitive) {
	app.gui.content = tview.NewTextView()
	app.gui.content.
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				nextSlide()
				return
			}
		}).
		SetScrollable(true).
		SetDynamicColors(true).
		SetBorderPadding(1, 1, 5, 5)
	return "WebContent", app.gui.content
}

func (gui *GUI) parseHTML(item hnapi.Item) {
	gui.content.Clear()
	gui.console.SetText("Loading page...")
	webCMD := exec.Command("w3m", "-dump", "-graph", "-X", "-cols", string(numCols), item.URL)

	app.main.QueueUpdateDraw(func() {
		stdOutput, err := webCMD.CombinedOutput()
		if err != nil {
			log.Print(err)
		}

		if _, err = gui.content.Write(stdOutput); err != nil {
			log.Print(err)
		}
	})

	gui.console.SetText("Page loaded.")
}
