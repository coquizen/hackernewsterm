package ui
//
//import (
//  "fmt"
//  "github.com/CaninoDev/hackernewsterm/internal/hackernews"
//  "log"
//	"sync"
//  "time"
//
//  "code.rocketnine.space/tslocum/cview"
//	"github.com/gdamore/tcell/v2"
//  "context"
//
//)
//
//var tabs []*tab
//var curTab = -1
//
//var termW, termH int
//
//var bottomBar = cview.NewInputField()
//
//var panes = cview.NewPanels()
//
//var tabbedPanes = cview.NewTabbedPanels()
//
//var layout = cview.NewFlex()
//
//var reformatMu = sync.Mutex{}
//
//var App = cview.NewApplication()
//
//func Init() {
//  App.EnableMouse(true)
//  App.SetRoot(layout, true)
//  App.SetAfterResizeFunc(func(width int, height int) {
//    termW = width
//    termH = height
//
//    //go func(t *tab) {
//    //  reformatMu.Lock()
//    //  for i := range tabs {
//    //    log.Print(i)
//    //  }
//    //  App.Draw()
//    //  reformatMu.Unlock()
//    //}(tabs[curTab])
//  })
//
//  TabbedPanels()
//
//  layout.SetDirection(cview.FlexRow)
//  layout.AddItem(panes, 0, 1, true)
//  layout.AddItem(bottomBar, 1, 1, false)
//
//  bottomBar.SetDoneFunc(func(key tcell.Key) {
//    tab := curTab
//
//    reset := func() {
//      bottomBar.SetLabel("")
//      App.SetFocus(tabs[tab].list)
//    }
//
//    if key == tcell.KeyEsc {
//      reset()
//      return
//    }
//
//  })
//
//  App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
//    _, ok := App.GetFocus().(*cview.Button)
//    if ok {
//      return event
//    }
//    _, ok = App.GetFocus().(*cview.InputField)
//    if ok {
//      return event
//    }
//    return event
//  })
//}
//
//func TabbedPanels(h hackernews.Handler) *cview.TabbedPanels{
//  tabbedPanes.SetTabTextColorFocused(HNOrange.TrueColor())
//  handler := hackernews.NewHandlerWithDefaultConfig(context.Context())
//
//  for index, panel := range panels {
//    table := cview.NewList()
//    func(index int) {
//      reset := func() {
//        table.Clear()
//        go func() {
//          for item := range handler.Subscribe(panel.EndPoint) {
//            tm := time.Unix(int64(item.Time()), 0).UTC().Format(time.RFC850)
//            listItem := cview.NewListItem(item.Title())
//            listItem.SetSecondaryText(fmt.Sprintf("(%s) by: %s on: %s score: %d", item.URL(), item.By(), tm, item.Score()))
//            table.AddItem(listItem)
//          }
//        }()
//      }
//      table.AddContextItem("Refresh", 'r', func(index int) {
//        reset()
//        table.ContextMenuList().SetItemEnabled(0, true)
//      })
//      tabbedPanes.AddTab(panel.Name, panel.Name, table)
//    }(index)
//  }
//  return tabbedPanes
//}
