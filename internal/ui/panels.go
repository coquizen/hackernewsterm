package ui

import (
	"code.rocketnine.space/tslocum/cview"
	"fmt"
	"github.com/CaninoDev/hackernewsterm/internal/hackernews"
	"time"
)

type Panel struct {
	Name     string
	EndPoint hackernews.RequestType
	MaxCount int
}


var (
	panels = []Panel{
		{Name: "New Stories", EndPoint: hackernews.NewStories, MaxCount: 500},
		{Name: "Top Stories", EndPoint: hackernews.TopStories, MaxCount: 500},
		{Name: "Best Stories", EndPoint: hackernews.BestStories, MaxCount: 500},
		{Name: "Ask", EndPoint: hackernews.AskStories, MaxCount: 200},
		{Name: "Show", EndPoint: hackernews.ShowStories, MaxCount: 200},
		{Name: "JobStories", EndPoint: hackernews.JobStories, MaxCount: 200},
	}
)

func (u *ui) Panels() {
	tabbedPanels := cview.NewTabbedPanels()
	tabbedPanels.SetTabTextColorFocused(HNOrange.TrueColor())
	for index, panel := range panels {
		list := cview.NewList()
		tabbedPanels.AddTab(panel.Name, panel.Name, list)
		func(index int) {
			reset := func() {
				list.Clear()
                                var idx = 0
				go func() {
                                  for  {
					item := <-u.firebase.Subscribe(panel.EndPoint)
					tm := time.Unix(int64(item.Time()), 0).UTC().Format(time.RFC850)
					listItem := cview.NewListItem(item.Title())
					listItem.SetSecondaryText(fmt.Sprintf("(%s) by: %s on: %s score: %d", item.URL(), item.By(), tm,
						item.Score()))
                                        list.InsertItem(idx, listItem)
                                        idx++
				}
                              }()
			}
			list.AddContextItem("Refresh", 'r', func(index int) {
				reset()
				list.ContextMenuList().SetItemEnabled(0, true)
			})
		}(index)
	}
	u.QueueUpdateDraw(func() {
		u.SetRoot(tabbedPanels, true)
	})

}

//
//postings, err := u.fb.Stories(u.ctx, panel.EndPoint)
//if err != nil {
//log.Fatalf("error retrieving stories from endpoint %s: %v", panel.Name, err)
//}
//for _, post := range postings {
//tm := time.Unix(int64(post.Time()), 0).UTC().Format(time.RFC850)
//listItem := cview.NewListItem(panel.Name)
//listItem.SetSecondaryText(fmt.Sprintf("(%s) by: %s on: %s score: %d", post.URL(), post.By(), tm,
//post.Score()))
//// listItem.SetSelectedFunc(showStory)
//list.AddItem(listItem)
//}
//list.ContextMenuList().SetItemEnabled(0, true)
