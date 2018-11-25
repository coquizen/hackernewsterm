package main

import (
	"fmt"
	"html"
	"log"
	"net/url"
	"os/exec"
	"sync"
	"time"

	"github.com/caninodev/hackernewsterm/hnapi"
	"github.com/gdamore/tcell"
	"github.com/rickb777/date/period"
	"github.com/rivo/tview"
)

var wg = sync.WaitGroup{}
var (
	cache   []hnapi.Item
	numCols int
)

type GUI struct {
	layout *tview.Flex
	//header   *tview.TextView
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

	gui.list = tview.NewList()
	gui.list.ShowSecondaryText(true)
	gui.list.SetBorder(true)
	gui.list.SetChangedFunc(updateDisplay)

	gui.content = tview.NewTextView()
	gui.content.SetDynamicColors(true).
		SetBorder(true)
	gui.content.SetScrollable(true)
	gui.content.SetBorderPadding(0, 0, 2, 2)
	_, _, numCols, _ = gui.content.GetInnerRect()

	placeNode := tview.NewTreeNode(".")
	gui.comments = tview.NewTreeView().
		SetGraphics(true).
		SetTopLevel(0).
		SetRoot(placeNode)

	gui.console = tview.NewTextView()
	gui.console.SetDynamicColors(true)

	var defaultRequest = &hnapi.Request{
		PostType: "top",
		NumPosts: 50,
	}

	go func(req *hnapi.Request) {
		gui.getPosts(defaultRequest)
	}(defaultRequest)

	gui.layout = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
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
	case tcell.KeyRune:
		if key.Rune() == 'j' {
			app.main.SetFocus(app.gui.content)
			x, y := app.gui.content.GetScrollOffset()
			app.gui.content.ScrollTo(x+1, y)
			app.main.SetFocus(app.gui.list)
		}
		if key.Rune() == 'k' {
			app.main.SetFocus(app.gui.content)
			x, y := app.gui.content.GetScrollOffset()
			app.gui.content.ScrollTo(x-1, y)
			app.main.SetFocus(app.gui.list)
		}

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

	var topLevelNode *tview.TreeNode
	cmt := fmt.Sprintf("ID: %d", cache[index].ID)
	topLevelNode = tview.NewTreeNode(cmt)

	go func() {
		parseHTML(cache[index])
		germinate(topLevelNode, cache[index])

	}()
	app.main.Draw()
	app.gui.console.Clear()
}

func (gui *GUI) renderListItem(item hnapi.Item, idx int) {
	//mainString := []string{"[yellow:-]", strconv.Itoa(idx), "[-:-:-] ", formatMainText(&item)}
	//mainText := strings.Join(mainString, "")
	//secondaryText := formatSubText(&item)
	m := formatMainText(&item)
	n := formatSubText(&item)
	gui.list.AddItem(*m, *n, rune(idx+1), nil)
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
	fmt.Fprint(app.gui.console, " Loading Page... ")
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
	fmt.Fprint(app.gui.console, " Page done.")
}

// Adapted from github.com/johnshiver/plankton/terminal/treeview.s
func createAllChildNodes(parentItem *hnapi.Item) *tview.TreeNode {
	var getChildrenNodes func(item *hnapi.Item) *tview.TreeNode
	getChildrenNodes = func(item *hnapi.Item) *tview.TreeNode {
		timeSince := time.Since(time.Unix(item.Time, 0))
		p, _ := period.NewOf(timeSince)
		timeStr := p.String()
		commentText := fmt.Sprintf("[::b]%s[::d] (%s) [-:-:-]-- %s", item.By, timeStr, html.UnescapeString(item.Text))
		app.gui.console.Clear()
		commentNode := *tview.NewTreeNode("").SetReference(item.ID).SetText(commentText)
		for _, childID := range item.Kids {
			childItem, _ := app.api.GetItem(childID)
			childNode := getChildrenNodes(childItem)
			commentNode.AddChild(childNode)
		}
		return &commentNode
	}
	return getChildrenNodes(parentItem)
}

func germinate(topNode *tview.TreeNode, item hnapi.Item) {
	fmt.Fprint(app.gui.console, " Loading comments...")
	for childID := range item.Kids {
		rootComment, _ := app.api.GetItem(childID)
		rootNode := createAllChildNodes(rootComment)
		rootNode.SetExpanded(true)
		rootNode.SetReference(childID)
		//commentText := fmt.Sprintf("comment: %d", childID)
		//commentText := fmt.Sprintf("[gre.en::b]%s[-::d] -- %s %d", rootComment.By, string(rootComment.Text), rootComment.ID)
		//log.Print(rootComment.Text)
		//rootNode.SetText(commentText)
		topNode.AddChild(rootNode)
	}
	app.gui.comments.SetRoot(topNode)
	// 	SetSelectedFunc(func(n *tview.TreeNode) {
	// 		original := n.GetReference().(*tview.TreeNode)
	// 		if original.IsExpanded() {
	// 			n.SetExpanded(!n.IsExpanded())
	// 		}
	app.main.Draw()
	fmt.Fprint(app.gui.console, "comments done!")
}
