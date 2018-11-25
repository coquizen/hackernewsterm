package main

import (
	"fmt"
	"html"
	"log"
	"net/url"
	"os/exec"
	"time"

	"github.com/caninodev/hackernewsterm/hnapi"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var (
	cache         []hnapi.Item
	numCols       int
	hnColorOrange tcell.Color
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

	gui.list = tview.NewList()
	gui.list.ShowSecondaryText(true)
	gui.list.SetChangedFunc(updateDisplay)

	gui.content = tview.NewTextView()
	gui.content.SetDynamicColors(true)
	gui.content.SetScrollable(true)

	placeNode := tview.NewTreeNode("")
	gui.comments = tview.NewTreeView().
		SetGraphics(true).
		SetTopLevel(0).
		SetRoot(placeNode)
	gui.comments.SetSelectedFunc(func(node *tview.TreeNode) {
		currentRefNode := node.GetReference().(*tview.TreeNode)
		if currentRefNode.IsExpanded() {
			node.SetExpanded(!node.IsExpanded())
		}
	})

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
			currentFocus := app.main.GetFocus()
			app.main.SetFocus(app.gui.content)
			x, y := app.gui.content.GetScrollOffset()
			app.gui.content.ScrollTo(x-1, y)
			app.main.SetFocus(currentFocus)
		}

	}
	return key
}

func (gui *GUI) getPosts(request *hnapi.Request) {
	gui.list.SetBorderColor(tcell.ColorSalmon)
	gui.list.SetTitle(request.PostType + " stories")

	idx := 0

	stream := app.api.GetPosts(request)
	cache = make([]hnapi.Item, request.NumPosts)
	itrString := []rune("abcdefghijklmnopqrstuvwxyz1234567890-=_+[]<>?!`~$%^@()")
	for item := range stream {
		cache[idx] = *item
		gui.renderListItem(cache[idx], itrString[idx])
		idx++
	}
	gui.list.SetBorder(false)

}

func updateDisplay(index int, _ string, _ string, _ rune) {

	var topLevelNode *tview.TreeNode
	cmt := fmt.Sprintf("ID: %d", cache[index].ID)
	topLevelNode = tview.NewTreeNode(cmt)
	_, _, numCols, _ = app.gui.content.GetInnerRect()

	go func() {
		parseHTML(cache[index])
		germinate(topLevelNode, cache[index])

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
	app.gui.content.SetBorder(true)
	app.gui.content.SetBorderColor(hnColorOrange)
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
	app.gui.content.SetBorder(false)
	fmt.Fprint(app.gui.console, " Page done.")
}

// Adapted from github.com/johnshiver/plankton/terminal/treeview.s
func createAllChildNodes(parentItem *hnapi.Item) *tview.TreeNode {
	var getChildrenNodes func(item *hnapi.Item) *tview.TreeNode
	getChildrenNodes = func(item *hnapi.Item) *tview.TreeNode {
		var commentNode tview.TreeNode
		for _, childID := range item.Kids {
			childItem, _ := app.api.GetItem(childID)
			commentText := fmt.Sprintf("[-:-:-]%s[::d] (%s) [-:-:-] %d %d %s", item.By, time.Unix(item.Time, 0), item.ID, item.Descendants, html.UnescapeString(item.Text))
			commentNode = *tview.NewTreeNode("").
				SetReference(item).
				SetText(commentText)

			if childItem.Descendants > 0 {
				childNode := getChildrenNodes(childItem)
				commentNode.AddChild(childNode)
			}
		}
		return &commentNode
	}
	return getChildrenNodes(parentItem)
}

func germinate(topNode *tview.TreeNode, item hnapi.Item) {
	app.gui.list.SetBorder(true)
	app.gui.list.SetBorderColor(hnColorOrange)
	fmt.Fprint(app.gui.console, " Loading comments...")
	for childID := range item.Kids {
		rootComment, _ := app.api.GetItem(childID)
		rootNode := createAllChildNodes(rootComment)
		rootNode.SetExpanded(true)
		rootNode.SetReference(childID)
		topNode.AddChild(rootNode)
	}
	app.gui.comments.SetRoot(topNode)
	app.main.Draw()
	app.gui.list.SetBorder(false)
	app.gui.console.Clear()
}
