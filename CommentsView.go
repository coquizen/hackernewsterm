package main

import (
	"fmt"
	"html"
	"log"
	"time"

	"github.com/caninodev/hackernewsterm/hnapi"
	"github.com/dustin/go-humanize"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type node struct {
	text     string
	expand   bool
	selected func()
	children []*node
}

// Adapted from github.com/johnshiver/plankton/terminal/treeview.go
func createChildrenCommentNodes(rootCommentItem *hnapi.Item) *tview.TreeNode {
	var addChildNode func(commentItem *hnapi.Item) *tview.TreeNode
	addChildNode = func(commentItem *hnapi.Item) *tview.TreeNode {
		tm := time.Unix(commentItem.Time, 0)
		commentText := fmt.Sprintf("[-:-:-]%s[::d] (%s) %d --[::-]%s", commentItem.By, humanize.Time(tm), len(commentItem.Kids), html.UnescapeString(commentItem.Text))
		commentNode := tview.NewTreeNode(commentText)
		commentNode.SetReference(commentItem)
		for _, kidID := range commentItem.Kids {
			kid, _ := app.api.GetItem(kidID)
			cNode := addChildNode(kid)
			if len(kid.Kids) > 0 {
				cNode.SetColor(tcell.ColorMediumSeaGreen)
			}
			commentNode.AddChild(cNode)
		}
		return commentNode
	}

	return addChildNode(rootCommentItem)
}

func germinate(storyItem hnapi.Item) {
	app.gui.comments.SetChangedFunc(func(node *tview.TreeNode) {
		item := node.GetReference().(*hnapi.Item)
		app.gui.commentsContent.SetText(html.UnescapeString(item.Text))
	})
	if _, err := fmt.Fprint(app.gui.console, "Loading comments..."); err != nil {
		log.Print(err)
	}

	var add func(targets *tview.TreeNode) *tview.TreeNode
	add = func(target *tview.TreeNode) *tview.TreeNode {
		for _, rootCommentID := range storyItem.Kids {
			rootComment, _ := app.api.GetItem(rootCommentID)
			rootCommentNode := createChildrenCommentNodes(rootComment)
			target.AddChild(rootCommentNode)
		}
		return target
	}

	storyNode := *tview.NewTreeNode(storyItem.Text)
	root := add(&storyNode)
	app.gui.comments.SetRoot(root).
		SetCurrentNode(root)
	app.gui.console.Clear()
	app.main.Draw()
}