package ui

import (
	"code.rocketnine.space/tslocum/cview"
	"fmt"
	"github.com/CaninoDev/hackernewsterm/internal/hackernews"
	"math"
	"strconv"
	"strings"
	"time"
)

type Panel struct {
	Name     string
	RequestType hackernews.RequestType
	Subscription hackernews.Subscription
}
var (
	panels = &[]Panel{
		{RequestType: hackernews.NewStories, Name: "New Stories"},
		{RequestType: hackernews.TopStories, Name: "Top Stories"},
		{RequestType: hackernews.BestStories, Name: "Best Stories"},
		{RequestType: hackernews.AskStories, Name: "Ask"},
		{RequestType: hackernews.ShowStories, Name: "Show"},
		{RequestType: hackernews.JobStories, Name: "JobStories"},
	}
)

func (u *ui) Panels() {
	tabbedPanels := cview.NewTabbedPanels()
	tabbedPanels.SetTabTextColorFocused(HNOrange.TrueColor())
	for index, panel := range *panels {
		func(index int) {
			var fetcher = hackernews.NewRequestHandler(panel.RequestType)
			list := cview.NewList()
			render := func() {
				var count = 36 / 2
				var tally = int(1)
				var subscription = hackernews.Subscribe(fetcher)
				go func() {
					for item := range subscription.Updates() {
						if tally != count {
							subscription.Command() <- hackernews.Play
							listItem := cview.NewListItem(item.Title())
							listItem.SetSecondaryText(fmt.Sprintf("By %s on %s -- (%d)", item.By(), TimeElapsed(time.Now(), time.Unix(int64(item.Time()), 0), true), item.Score()))
							u.QueueUpdateDraw(func() {
									list.AddItem(listItem)
								})
							tally++
						} else {
							tally = 1
							subscription.Command() <- hackernews.Pause
						}
					}
				}()
			}
			reset := func() {
				list.Clear()
				fetcher = hackernews.NewRequestHandler(panel.RequestType)
				go render()
			}
			list.AddContextItem("Reset", 'r', func(index int) {
				reset()
				list.ContextMenuList().SetItemEnabled(1, true)
			})
			tabbedPanels.AddTab(panel.Name, panel.Name, list)
			reset()
		}(index)
	}
	u.SetRoot(tabbedPanels, true)

}

func s(x float64) string {
	if int(x) == 1 {
		return ""
	}
	return "s"
}

func TimeElapsed(now time.Time, then time.Time, full bool) string {
	var parts []string
	var text string

	year2, month2, day2 := now.Date()
	hour2, minute2, second2 := now.Clock()

	year1, month1, day1 := then.Date()
	hour1, minute1, second1 := then.Clock()

	year := math.Abs(float64(int(year2 - year1)))
	month := math.Abs(float64(int(month2 - month1)))
	day := math.Abs(float64(int(day2 - day1)))
	hour := math.Abs(float64(int(hour2 - hour1)))
	minute := math.Abs(float64(int(minute2 - minute1)))
	second := math.Abs(float64(int(second2 - second1)))

	week := math.Floor(day / 7)

	if year > 0 {
		parts = append(parts, strconv.Itoa(int(year))+" year"+s(year))
	}

	if month > 0 {
		parts = append(parts, strconv.Itoa(int(month))+" month"+s(month))
	}

	if week > 0 {
		parts = append(parts, strconv.Itoa(int(week))+" week"+s(week))
	}

	if day > 0 {
		parts = append(parts, strconv.Itoa(int(day))+" day"+s(day))
	}

	if hour > 0 {
		parts = append(parts, strconv.Itoa(int(hour))+" hour"+s(hour))
	}

	if minute > 0 {
		parts = append(parts, strconv.Itoa(int(minute))+" minute"+s(minute))
	}

	if second > 0 {
		parts = append(parts, strconv.Itoa(int(second))+" second"+s(second))
	}

	if now.After(then) {
		text = " ago"
	} else {
		text = " after"
	}

	if len(parts) == 0 {
		return "just now"
	}

	if full {
		return strings.Join(parts, ", ") + text
	}
	return parts[0] + text
}
