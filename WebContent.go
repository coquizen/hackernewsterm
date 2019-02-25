package main

import (
	"log"

	"github.com/go-shiori/go-readability"

	nurl "net/url"
	"time"

	"github.com/caninodev/hackernewsterm/hnapi"
	"github.com/rivo/tview"
)

// WebContent is a page that will render the selected item's url
func WebContent(nextSlide func()) (title string, content tview.Primitive) {
	app.gui.content = tview.NewTextView()
	app.gui.content.
		SetScrollable(true).
		SetDynamicColors(true).
		SetBorderPadding(1, 1, 5, 5)
	return "WebContent", app.gui.content
}

func (gui *GUI) parseHTML(item hnapi.Item) {
	gui.console.SetText("Loading page...")
	app.main.QueueUpdateDraw(func() {
		if parsedURL, err := nurl.Parse(item.URL); err != nil {
			log.Print(err)
			gui.console.SetText("URL parsing error!")
		} else {
			article, _ := readability.FromURL(parsedURL, 7*time.Second)
			gui.content.Write([]byte(article.Content))
			gui.console.SetText("Page successfully loaded.")
		}
	})
}
