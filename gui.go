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
	//cache   [50]*hnapi.Item
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

	var defaultRequest = &hnapi.Request{
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

func (gui *GUI) getPosts(request *hnapi.Request) {

	gui.list.SetTitle(request.PostType + " stories")

	idx := 0

	stream := app.api.GetPosts(request)
	cache := make([]*hnapi.Item, request.NumPosts)

	for item := range stream {
		cache[idx] = item
		gui.render(*cache[idx], idx)
		if idx == 0 {
			gui.parseHTML(*cache[idx])
			gui.comments = gui.germinate(cache[idx])
		}
		idx++
	}
}

func (gui *GUI) render(item hnapi.Item, idx int) {
	mainString := []string{"[[yellow:b]", strconv.Itoa(idx), "[-:-:-]] ", formatMainText(&item)}
	mainText := strings.Join(mainString, "")
	secondaryText := formatSubText(&item)
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

func (gui *GUI) parseHTML(item hnapi.Item) {

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

// Adapted from github.com/johnshiver/plankton/terminal/treeview.go
func createAllChildNodes(nodeID int) *tview.TreeNode {
	var addNode func(id int) *tview.TreeNode
	addNode = func(nodeID int) *tview.TreeNode {
		comment, _ := app.api.GetItem(nodeID)
		commentText := fmt.Sprintf("[::b]%s[::d] (%d) [-:-:-]--[-:-:-]%s", comment.By, comment.Time, comment.Text)
		newNode := *tview.NewTreeNode(commentText)
		for _, childID := range comment.Kids {
			childComment, _ := app.api.GetItem(childID)
			cNode := addNode(childComment.ID)
			newNode.AddChild(cNode)
		}
		return &newNode
	}
	return addNode(nodeID)
}

func (gui *GUI) germinate(item *hnapi.Item) (tree *tview.TreeView) {
	for id := range item.Kids {
		rootComment, _ := app.api.GetItem(id)
		tree := tview.NewTreeView()
		rootNode := createAllChildNodes(id)
		rootNode.SetExpanded(false)

		tree.SetBorder(true).
			SetTitle(rootComment.By)
		tree.SetAlign(false).
			SetTopLevel(0).
			SetGraphics(true).
			SetPrefixes(nil)

		var add func(targetNode *tview.TreeNode) *tview.TreeNode
		add = func(targetNode *tview.TreeNode) *tview.TreeNode {
			commentText := fmt.Sprintf("[::b]%s[::d] (%d) [-:-:-]--[-:-:-]%s", rootComment.By, rootComment.Time, rootComment.Text)
			commentNode := tview.NewTreeNode(commentText).
				SetSelectable(true).
				SetExpanded(targetNode == rootNode).
				SetReference(targetNode)
			if targetNode.IsExpanded() {
				targetNode.SetColor(tcell.ColorGreen)
			}

			for _, childNode := range rootNode.GetChildren() {
				commentNode.AddChild(add(childNode))
			}
			return commentNode
		}

		root := add(rootNode)
		tree.SetRoot(root).
			SetCurrentNode(root).
			SetSelectedFunc(func(n *tview.TreeNode) {
				original := n.GetReference().(*tview.TreeNode)
				if original.IsExpanded() {
					n.SetExpanded(!n.IsExpanded())
				}
			})
		tree.GetRoot().ExpandAll()

	}
	return tree
}
