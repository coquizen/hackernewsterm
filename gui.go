package main

import (
	"container/ring"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"reflect"

	"github.com/caninodev/hackernewsterm/hnapi"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var (
	cache         []hnapi.Item
	numCols       int
	hnColorOrange tcell.Color
	finder        *ring.Ring
)

type GUI struct {
	layout   *tview.Flex
	list     *tview.List
	content  *tview.TextView
	comments *tview.TreeView
	console  *tview.TextView
}

// Create establishes the ui and widget parameters
func (gui *GUI) Create() {
	//gui.header = tview.NewTextView().
	//	SetDynamicColors(true).
	//	SetRegions(true).
	//	SetWrap(false)
	//
	//gui.header.SetTextAlign(tview.AlignCenter)

	hnColorOrange = tcell.NewRGBColor(238, 111, 45)

	finder = ring.New(3)

	gui.list = tview.NewList()
	gui.list.ShowSecondaryText(true)
	gui.list.SetChangedFunc(updateDisplay)
	finder.Value = gui.list
	finder.Next()

	gui.content = tview.NewTextView()
	gui.content.SetDynamicColors(true)
	gui.content.SetScrollable(true)
	finder.Value = gui.content
	finder.Next()

	placeNode := tview.NewTreeNode("Loading...")
	gui.comments = tview.NewTreeView().
		SetRoot(placeNode)
	finder.Value = gui.comments
	finder.Next()

	gui.console = tview.NewTextView()
	gui.console.SetDynamicColors(true)

	var defaultRequest = &hnapi.Request{
		PostType: "top",
		NumPosts: 50,
	}

	go func(req *hnapi.Request) {
		gui.getPosts(defaultRequest)
	}(defaultRequest)

	// The following produces the Tall layout (one main pane to the left with the other two divided vertically to the right
	gui.layout = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexColumn).
			AddItem(gui.list, 0, 1, true).
			AddItem(tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(gui.content, 0, 1, false).
				AddItem(gui.comments, 0, 1, false), 0, 1, false), 0, 1, true).
		AddItem(gui.console, 1, 1, false)
}

func (gui *GUI) keyHandler(key *tcell.EventKey) *tcell.EventKey {
	switch key.Key() {
	case tcell.KeyEsc:
		app.main.Stop()
	case tcell.KeyTab:
		gui.changeFinderFocus()

	}
	return key
}

func (gui *GUI) changeFinderFocus() {
	currentlyFocused := finder.Value.(tview.Primitive)
	finder.Next()
	newlyFocused := finder.Value.(tview.Primitive)
	app.main.SetFocus(newlyFocused)
	logCFocus := fmt.Sprintf("CurrentlyFocused: %#v, newlyFocused: %#v", reflect.TypeOf(currentlyFocused), reflect.TypeOf(newlyFocused))
	gui.console.SetText(logCFocus)
	// newlyFocused.SetBorder(true).SetBorderColor(hnColorOrange)

}

func (gui *GUI) getPosts(request *hnapi.Request) {
	gui.list.SetTitle(request.PostType + " stories")

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

func updateDisplay(index int, _ string, _ string, _ rune) {
	_, _, numCols, _ = app.gui.content.GetInnerRect()

	go func() {
		parseHTML(cache[index])
		germinate(cache[index])

	}()
	app.main.Draw()
	app.gui.console.Clear()
}

func (gui *GUI) renderListItem(item hnapi.Item, idx rune) {
	//mainString := []string{"[yellow:-]", strconv.Itoa(idx), "[-:-:-] ", formatMainText(&item)}
	//mainText := strings.Join(mainString, "")
	//secondaryText := formatSubText(&item)
	m := formatMainText(&item)
	n := formatSubText(&item)
	gui.list.AddItem(*m, *n, idx, nil)

}

func formatMainText(item *hnapi.Item) *string {
	addr, _ := url.Parse(item.URL)
	mainText := fmt.Sprintf("[::b] %s [::d](%s)[::-]", item.Title, string(addr.Host))
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

	str := fmt.Sprintf("[-::d]score: %s %d [-::-][::d] comments:[::-] %d [::d] by:[green::-] %s [-:-:-]", scoreColor, item.Score, item.Descendants, item.By)
	return &str
}

func parseHTML(item hnapi.Item) {
	app.gui.content.SetTitle(item.URL)
	webCMD := exec.Command("w3m", "-dump", "-graph", "-X", "-cols", string(numCols), item.URL)

	app.main.QueueUpdateDraw(func() {
		stdOutput, _ := webCMD.CombinedOutput()
		_, err := app.gui.content.Write(stdOutput)
		if err != nil {
			log.Print(err)
		}
	})

	fmt.Fprint(app.gui.console, " Page done.")
}
