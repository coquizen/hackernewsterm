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

var (
	treeNextSlide func()
)

// Adapted from github.com/johnshiver/plankton/terminal/treeview.go
func createChildrenCommentNodes(rootCommentItem *hnapi.Item) *tview.TreeNode {
	var addChildNode func(commentItem *hnapi.Item) *tview.TreeNode
	addChildNode = func(commentItem *hnapi.Item) *tview.TreeNode {
		tm := time.Unix(commentItem.Time, 0)
		commentText := fmt.Sprintf("[-:-:-]%s[::d] (%s) %d replies", commentItem.By, humanize.Time(tm), len(commentItem.Kids))
		commentNode := tview.NewTreeNode(commentText)
		commentNode.SetReference(commentItem)
		for _, kidID := range commentItem.Kids {
			kid, _ := app.api.GetItem(kidID)
			cNode := addChildNode(kid)
			if len(kid.Kids) > 0 {
				cNode.SetColor(tcell.ColorMediumSeaGreen)
				cNode.SetSelectable(true)
				cNode.SetSelectedFunc(func() {
					cNode.SetExpanded(!cNode.IsExpanded())
				})
			}
			commentNode.AddChild(cNode)
		}
		return commentNode
	}

	return addChildNode(rootCommentItem)
}

// Comments is a page that will render the comments tree as well as the selected comments
func Comments(nextSlide func()) (title string, content tview.Primitive) {
	treeNextSlide = nextSlide

	placeNode := tview.NewTreeNode("Loading...")

	app.gui.commentsContent = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true)

	app.gui.comments = tview.NewTreeView()
	app.gui.comments.SetRoot(placeNode).
		SetBorder(true).
		SetTitle("Comments")
	app.gui.comments.SetSelectedFunc(func(n *tview.TreeNode) {
		item := n.GetReference().(*hnapi.Item)
		app.gui.commentsContent.
			SetText(html.UnescapeString(item.Text))
	})

	return "Comments", tview.NewFlex().
		AddItem(app.gui.comments, 0, 2, true).
		AddItem(app.gui.commentsContent, 0, 5, false)
}

func germinate(storyItem hnapi.Item) {
	if _, err := fmt.Fprint(app.gui.console, "Loading comments..."); err != nil {
		log.Print(err)
	}

	// var add func(targets *tview.TreeNode) *tview.TreeNode
	add := func(target *tview.TreeNode) *tview.TreeNode {
		for _, rootCommentID := range storyItem.Kids {
			rootComment, err := app.api.GetItem(rootCommentID)
			if err != nil {
				log.Print(err)
			}
			if _, err := fmt.Fprintf(app.gui.console, "Created node for %#v... ", rootComment); err != nil {
				log.Print(err)
			}
			rootCommentNode := createChildrenCommentNodes(rootComment)
			if len(rootComment.Kids) > 0 {
				rootCommentNode.SetSelectable(true).
					SetSelectedFunc(func() {
						rootCommentNode.SetExpanded(!rootCommentNode.IsExpanded())
						rootCommentNode.SetColor(tcell.ColorGreen)
					})
			}
			target.AddChild(rootCommentNode)
		}
		return target
	}

	storyNode := *tview.NewTreeNode(storyItem.Text)
	root := add(&storyNode)
	app.gui.comments = tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)
	app.gui.console.Clear()
}
