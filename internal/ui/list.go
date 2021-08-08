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
                u.Lock()
                for i:= 0; i <20; i++{
			select {
			case item := <-u.firebase.Subscribe(hackernews.NewStories):
				if item.Type() == "story" {
					tm := time.Unix(int64(item.Time()), 0).UTC().Format(time.RFC1123)
					listItem := cview.NewListItem(item.Title())
					listItem.SetSecondaryText(fmt.Sprintf("(%s) by: %s, on %s", item.URL(), item.By(), tm))
					lv.AddItem(listItem)
				}

			}
		}
                u.Unlock()
		quitItem := cview.NewListItem("Quit")
		quitItem.SetSecondaryText("Press to exit")
		quitItem.SetShortcut('q')
		quitItem.SetSelectedFunc(func() {
			u.Stop()
		})
		lv.AddItem(quitItem)
		lv.ContextMenuList().SetItemEnabled(1, false)
	}
	lv.AddContextItem("Delete item", 'd', func(index int) {
		lv.RemoveItem(index)

		if lv.GetItemCount() == 0 {
			lv.ContextMenuList().SetItemEnabled(0, false)
			lv.ContextMenuList().SetItemEnabled(1, false)
		}
		lv.ContextMenuList().SetItemEnabled(3, true)
	})
	lv.AddContextItem("Reset", 'r', func(index int) {
		reset()
	})
	app.SetRoot(lv, true)
	reset()
	// app.SetRoot(lv, true)
	// if err := app.Run(); err != nil {
	// 	panic(err)
	// }

}
