package main

import (
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/caninodev/hackernewsterm/hnapi"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var wg = sync.WaitGroup{}
var (
	cache   []hnapi.Item
	numCols int
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
	gui.list = tview.NewList()
	gui.list.ShowSecondaryText(true)
	gui.list.SetBorder(true)
	gui.list.SetChangedFunc(updateDisplay)

	gui.content = tview.NewTextView()
	gui.content.SetDynamicColors(true).
		SetBorder(true)
	gui.content.SetBorderPadding(0, 0, 2, 2)
	_, _, numCols, _ = gui.content.GetInnerRect()

	gui.comments = tview.NewTreeView()
	placeNode := tview.NewTreeNode(".")
	gui.comments.SetRoot(placeNode)

	gui.console = tview.NewTextView()
	gui.console.SetDynamicColors(true)


	var defaultRequest = &hnapi.Request{
		PostType: "top",
		NumPosts: 50,
	}

	go func(req *hnapi.Request) {
		gui.getPosts(defaultRequest)
	}(defaultRequest)

	gui.layout = tview.NewFlex()
	gui.layout.SetDirection(tview.FlexRow)
	gui.layout.AddItem(tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(gui.list, 0, 1, true).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(gui.content, 0, 1, false).
			AddItem(gui.comments, 0, 1, false), 0, 1, false), 0, 1, true)
	gui.layout.AddItem(gui.console, 1, 1, false)
}

func (gui *GUI) KeyHandler(key *tcell.EventKey) *tcell.EventKey {
	switch key.Key() {
	case tcell.KeyEsc:
		app.main.Stop()
	}
	return key
}

func (gui *GUI) getPosts(request *hnapi.Request) {

	gui.list.SetTitle(request.PostType + " stories")

	idx := 0

	stream := app.api.GetPosts(request)
	cache = make([]hnapi.Item, request.NumPosts)
	for item := range stream {
		cache[idx] = *item
		gui.renderListItem(cache[idx], idx)
		idx++
	}
}

func updateDisplay(index int, _ string, _ string, _ rune) {
	go func() {
		parseHTML(cache[index])
		germinate(cache[index])
	}()
	app.main.Draw()
}

func (gui *GUI) renderListItem(item hnapi.Item, idx int) {
	mainString := []string{"[yellow:-]", strconv.Itoa(idx), "[-:-:-] ", formatMainText(&item)}
	mainText := strings.Join(mainString, "")
	secondaryText := formatSubText(&item)
	_, _ = gui.console.Write([]byte(mainText))
	gui.list.AddItem(mainText, secondaryText, rune(idx+1), nil)
}

func formatMainText(item *hnapi.Item) (mText string) {
	addr, _ := url.Parse(item.URL)
	mainText := []string{"[::b]", item.Title, "[::d]", " (", string(addr.Host), ")"}
	mText = strings.Join(mainText, "")
	return mText
}

func formatSubText(item *hnapi.Item) string {
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

func parseHTML(item hnapi.Item) {
	app.gui.content.Clear()
	app.gui.content.SetTitle(item.URL)
	webCMD := exec.Command("w3m", "-dump", "-graph", "-X", "-cols", string(numCols), item.URL)

	app.main.QueueUpdateDraw(func() {
		stdOutput, _ := webCMD.CombinedOutput()
		_, err := app.gui.content.Write(stdOutput)
		if err != nil {
			log.Print(err)
		}
	})
	app.main.Draw()
}

// Adapted from github.com/johnshiver/plankton/terminal/treeview.go
func createAllChildNodes(parentItem *hnapi.Item) *tview.TreeNode {
	var addNode func(item *hnapi.Item) *tview.TreeNode
	addNode = func(item *hnapi.Item) *tview.TreeNode {
		commentText := fmt.Sprintf("[::b]%s[::d] (%d) [-:-:-]-- %s", item.By, item.Time, item.Text)
		commentNode := *tview.NewTreeNode(commentText).SetReference(item.ID)
		for _, childID := range item.Kids {
			childItem, _ := app.api.GetItem(childID)
			childNode := addNode(childItem)
			commentNode.AddChild(childNode)
		}
		return &commentNode
	}
	return addNode(parentItem)
}

func germinate(item hnapi.Item) {
	app.gui.layout.SetTitle("Comments")
	var topLevelNode *tview.TreeNode
	topLevelNode = tview.NewTreeNode("Comments")

	for childID := range item.Kids {
		rootComment, _ := app.api.GetItem(childID)
		rootNode := createAllChildNodes(rootComment)
		rootNode.SetExpanded(false)
		rootNode.SetReference(childID)
		commentText := fmt.Sprintf("[blue::b]%s[-::d] -- %s", rootComment.By, rootComment.Text)
		rootNode.SetText(commentText)
		topLevelNode.AddChild(rootNode)
	}
	app.gui.comments.SetRoot(topLevelNode)
	// 	SetSelectedFunc(func(n *tview.TreeNode) {
	// 		original := n.GetReference().(*tview.TreeNode)
	// 		if original.IsExpanded() {
	// 			n.SetExpanded(!n.IsExpanded())
	// 		}
	app.main.Draw()
}
