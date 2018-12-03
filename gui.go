package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/caninodev/hackernewsterm/hnapi"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var (
	cache         []hnapi.Item
	numCols       int
	hnColorOrange tcell.Color
	nextSlide     func()
)

// Slide is a function which returns the slide's main primitive and its title.
// It receives a "nextSlide" function which can be called to advance the
// presentation to the next slide.
type Slide func(nextSlide func()) (title string, content tview.Primitive)

// GUI structure contains all the UI element for the application.
type GUI struct {
	layout          *tview.Flex
	list            *tview.List
	content         *tview.TextView
	commentTitle    *tview.TextView
	comments        *tview.TreeView
	commentsContent *tview.TextView
	console         *tview.TextView
	pages           *tview.Pages
}

// Create establishes the ui and widget parameters
func (gui *GUI) Create() {
	hnColorOrange = tcell.NewRGBColor(238, 111, 45)

	var defaultRequest = &hnapi.Request{
		PostType: "top",
		NumPosts: 50,
	}

	gui.topPane()
	gui.bottomPane()

	go func(req *hnapi.Request) {
		gui.getPosts(defaultRequest)
	}(defaultRequest)

	gui.console = tview.NewTextView()
	gui.console.
		SetDynamicColors(true).
		SetBackgroundColor(hnColorOrange)

	gui.layout = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(gui.list, 0, 2, true).
		AddItem(gui.pages, 0, 5, true).
		AddItem(gui.console, 1, 1, false)
}

func (gui *GUI) topPane() {
	gui.list = tview.NewList()
	gui.list.ShowSecondaryText(true).
		SetChangedFunc(gui.updateDisplay)
}

func (gui *GUI) bottomPane() {
	slides := []Slide{
		WebContent,
		Comments,
	}
	gui.pages = tview.NewPages()

	currentSlide := 0

	// previousSlide := func() {
	// 	currentSlide = (currentSlide - 1 + len(slides)) % len(slides)
	// 	gui.pages.SwitchToPage(strconv.Itoa(currentSlide))
	// }

	nextSlide = func() {
		currentSlide = (currentSlide + 1) % len(slides)
		gui.pages.SwitchToPage(strconv.Itoa(currentSlide))
	}

	for index, slide := range slides {
		_, primitive := slide(nextSlide)
		gui.pages.AddPage(strconv.Itoa(index), primitive, true, index == currentSlide)
	}
}

func (gui *GUI) keyHandler(key *tcell.EventKey) *tcell.EventKey {
	switch key.Key() {
	case tcell.KeyEsc:
		app.main.Stop()
	case tcell.KeyRune:
		if key.Rune() == 'j' {
			app.main.SetFocus(app.gui.content)
			x, y := gui.content.GetScrollOffset()
			gui.content.ScrollTo(x+1, y)
			app.main.SetFocus(app.gui.list)
		}
		if key.Rune() == 'k' {
			currentFocus := app.main.GetFocus()
			app.main.SetFocus(app.gui.content)
			x, y := app.gui.content.GetScrollOffset()
			gui.content.ScrollTo(x-1, y)
			app.main.SetFocus(currentFocus)
		}
		if key.Rune() == 'C' {
			nextSlide()
		}
		if key.Rune() == 'I' {
			n := gui.comments.GetCurrentNode()
			item := n.GetReference().(*hnapi.Item)
			_, err := fmt.Fprint(gui.console, item.Text)
			if err != nil {
				log.Print(err)
			}
		}
		if key.Rune() == 'F' {
			if gui.list.HasFocus() {
				app.main.SetFocus(gui.pages)
				gui.pages.SetBorder(true)
				gui.list.SetBorder(false)
			} else {
				app.main.SetFocus(gui.list)
				gui.list.SetBorder(true)
				gui.pages.SetBorder(false)
			}
		}
	}
	return key
}

func (gui *GUI) getPosts(request *hnapi.Request) {
	gui.list.SetTitle(" " + request.PostType + " stories ")

	idx := 0

	stream := app.api.GetPosts(request)
	cache = make([]hnapi.Item, request.NumPosts)
	itrString := []rune("abcdefghilmnopqrstuvwxyz1234567890-=_+[]<>?!`~$%^@()")

	for item := range stream {
		cache[idx] = *item
		gui.renderListItem(*item, itrString[idx])
		idx++
	}

}

func (gui *GUI) updateDisplay(index int, _ string, _ string, _ rune) {
	_, _, numCols, _ = gui.content.GetInnerRect()

	go func(index int) {
		gui.parseHTML(cache[index])
		gui.germinate(cache[index])
	}(index)

}

func (gui *GUI) renderListItem(item hnapi.Item, idx rune) {
	m := formatMainText(&item)
	n := formatSubText(&item)
	app.main.QueueUpdateDraw(func() {
		gui.list.AddItem(*m, *n, idx, nil)
	})
}

func formatMainText(item *hnapi.Item) *string {
	addr, err := url.Parse(item.URL)
	if err != nil {
		log.Print(err)
	}
	mainText := fmt.Sprintf("[::b] %s [::d](%s)[::-]", item.Title, addr.Host)
	return &mainText
}

func formatSubText(item *hnapi.Item) *string {
	i := item.Score
	var scoreColor string
	switch {
	case i < 25:
		scoreColor = "[red::d]"
	case i < 75:
		scoreColor = "[orange::-]"
	case i < 100:
		scoreColor = "[yellow::b]"
	case i >= 100:
		scoreColor = "[green::b]"
	}

	str := fmt.Sprintf("[-::d] %s %d points,[-:-:-] %d [::d]comments, by:[green::-] %s [-:-:-]", scoreColor, item.Score, item.Descendants, item.By)
	return &str
}
