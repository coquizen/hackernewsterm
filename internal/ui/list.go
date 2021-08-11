package ui

//
import (
	"fmt"
	"time"

	"github.com/CaninoDev/hackernewsterm/internal/hackernews"

	"code.rocketnine.space/tslocum/cview"
)

type Headers [][]string

func (u *ui) ListView() {
	lv := cview.NewList()
	reset := func() {
		lv.Clear()
		var request = hackernews.NewHandler(hackernews.NewStories)
		var count = 37 / 2
		var tally = 1
		var subscription = hackernews.Subscribe(request)
		go func() {
			for item := range subscription.Updates() {
				if tally != count {
					subscription.Command() <- hackernews.Play
					listItem := cview.NewListItem(item.Title())
					listItem.SetSecondaryText(fmt.Sprintf("By %s posted %s ago", item.By(), TimeElapsed(time.Now(), time.Unix(int64(item.Time()), 0).UTC(), true)))
					u.QueueUpdateDraw(func() {
						lv.AddItem(listItem)
					})
					tally++
				} else {
					tally = 1
					subscription.Command() <- hackernews.Pause
					time.Sleep(5 * time.Second)
					subscription.Command() <- hackernews.Play
				}
			}
		}()
	}
	u.SetRoot(lv, true)
	go reset()
}
//	//quitItem := cview.NewListItem("Quit")
//	//quitItem.SetSecondaryText("Press to exit")
//	//quitItem.SetShortcut('q')
//	//quitItem.SetSelectedFunc(func() {
//	//	u.Stop()
//	//})
//	//lv.AddItem(quitItem)
//	//lv.ContextMenuList().SetItemEnabled(1, false)
//	//lv.AddContextItem("Delete item", 'd', func(index int) {
//	//	lv.RemoveItem(index)
//	//
//	//	if lv.GetItemCount() == 0 {
//	//		lv.ContextMenuList().SetItemEnabled(0, false)
//	//		lv.ContextMenuList().SetItemEnabled(1, false)
//	//	}
//	//	lv.ContextMenuList().SetItemEnabled(3, true)
//	//})
//	//lv.AddContextItem("Reset", 'r', func(index int) {
//	//	reset()
//	//})
//	u.SetRoot(lv, true)
//	reset()
//	// app.SetRoot(lv, true)
//	// if err := app.Run(); err != nil {
//	// 	panic(err)
//	// }
//
//}
