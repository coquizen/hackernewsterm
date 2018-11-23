package main

import (
	"log"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	models "github.com/caninodev/hackernewsterm/hnapi"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var wg = sync.WaitGroup{}
var (
	//cache   [50]*models.Item
	numCols int
)

type GUI struct {
	layout   *tview.Flex
	list     *tview.List
	content  *tview.TextView
	comments *tview.TreeView
}

// Create establishes the ui and widget parameters
func (gui *GUI) Create() {
	gui.list = tview.NewList()
	gui.list.ShowSecondaryText(true)
	gui.list.SetBorder(true)

	gui.content = tview.NewTextView()
	gui.content.SetDynamicColors(true).
		SetBorder(true)
	gui.content.SetBorderPadding(0, 0, 2, 2)
	_, _, numCols, _ = gui.content.GetInnerRect()

	gui.comments = tview.NewTreeView()
	gui.comments.SetBorder(true)

	var defaultRequest = &models.Request{
		PostType: "top",
		NumPosts: 50,
	}

	go gui.getPosts(defaultRequest)

	gui.layout = tview.NewFlex()
	gui.layout.SetDirection(tview.FlexColumn)
	gui.layout.AddItem(gui.list, 0, 1, true)
	gui.layout.AddItem(tview.NewFlex().
		SetDirection(tview.FlexRow).AddItem(gui.content, 0, 1, false).
		AddItem(gui.comments, 0, 1, false), 0, 1, false)

}

func (gui *GUI) KeyHandler(key *tcell.EventKey) *tcell.EventKey {
	switch key.Key() {
	case tcell.KeyEsc:
		app.main.Stop()
	}
	return key
}

func (gui *GUI) getPosts(request *models.Request) {

	gui.list.SetTitle(request.PostType + " stories")

	idx := 0

	stream := app.api.GetPosts(request)
	cache := make([]*models.Item, request.NumPosts)

	for item := range stream {
		cache[idx] = item
		gui.render(*cache[idx], idx)
		if idx == 0 {
			gui.parseHTML(*cache[idx])
		}
		idx++
	}
}

func (gui *GUI) render(item models.Item, idx int) {
	mainString := []string{"[[yellow:b]", strconv.Itoa(idx), "[-:-:-]] ", formatMainText(&item)}
	mainText := strings.Join(mainString, "")
	secondaryText := formatSubText(&item)
	gui.list.AddItem(mainText, secondaryText, rune(idx+1), nil)
}

func formatMainText(item *models.Item) (mText string) {
	addr, _ := url.Parse(item.URL)
	mainText := []string{"[::b]", item.Title, "[::d]", " (", string(addr.Host), ")"}
	mText = strings.Join(mainText, "")
	return mText
}

func formatSubText(item *models.Item) string {
	score := strconv.Itoa(item.Score)

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

	str := []string{"[-::d]score: ", scoreColor, score, "[-::-][::d] comments:[::-] ", strconv.Itoa(item.Descendants), "[::d] by:[blue::-]", item.By, "[-:-:-]"}
	return strings.Join(str, "")
}

func (gui *GUI) parseHTML(item models.Item) {

	gui.content.SetTitle(item.URL)

	webCMD := exec.Command("w3m", "-dump", "-graph", "-X", "-cols", string(numCols), item.URL)

	app.main.QueueUpdateDraw(func() {
		stdOutput, _ := webCMD.CombinedOutput()
		_, err := gui.content.Write(stdOutput)
		if err != nil {
			log.Print(err)
		}
	})
	app.main.Draw()
}

// func (gui *GUI) germinate(item *Item) {
// 	gui.comments.SetBorder(false)
// 	rootNode := tview.NewTreeNode("root").
// 		SetText(item.Text)

// 	gui.comments.SetRoot(rootNode)
// }
