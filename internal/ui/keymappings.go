package ui
//
// import (
// 	"fmt"
// 	"strings"
//
// 	"code.rocketnine.space/tslocum/cbind"
//         "github.com/gdamore/tcell/v2"
// )
//
// const (
// 	actionNewsStories = "newsstories"
// 	actionTopStories  = "topstories"
// 	actionJobsStories = "jobstories"
// 	actionBestStories = "beststories"
// 	actionComments    = "view-comments"
// 	actionLoadURL     = "view-in-browser"
// 	actionRefresh     = "refresh-items"
// 	actionSortByTime  = "time-sort"
// 	actionSortByScore = "score-sort"
// )
//
// var actionHandlers = map[string]func(){
// 	actionNewsStories: selectNewsTab,
// 	actionTopStories:  selectTopsTab,
// 	actionJobsStories: selectJobsTab,
// 	actionBestStories: selectBestTab,
// 	actionComments:    viewComments,
// 	actionLoadURL:     viewInBrowser,
// 	actionRefresh:     reloadItems,
// 	actionSortByTime:  sortItemsByTime,
// 	actionSortByScore: sortItemssByScore,
// }
//
// var inputConfig = cbind.NewConfiguration()
//
// func wrapEventHandler(f func()) func(_ *tcell.EventKey) *tcell.EventKey {
// 	return func(_ *tcell.EventKey) *tcell.EventKey {
// 		f()
// 		return nil
// 	}
// }
//
// func setKeyBinds() error {
// 	if len(config.Input) == 0 {
// 		setDefaultKeyBinds()
// 	}
//
// 	for a, keys := range config.Input {
// 		a = strings.ToLower(a)
// 		handler := actionHandlers[a]
// 		if handler == nil {
// 			return fmt.Errorf("failed to set keybind for %s: unknown action", a)
// 		}
//
// 		for _, k := range keys {
// 			mod, key, ch, err := cbind.Decode(k)
// 			if err != nil {
// 				return fmt.Errorf("failed to set keybind %s for %s: %s", k, a, err)
// 			}
//
// 			if key == tcell.KeyRune {
// 				inputConfig.SetRune(mod, ch, wrapEventHandler(handler))
// 			} else {
// 				inputConfig.SetKey(mod, key, wrapEventHandler(handler))
// 			}
// 		}
// 	}
//
// 	return nil
// }
//
// func setDefaultKeyBinds() {
// 	config.Input = map[string][]string{
// 		actionNewsStories: {"N"},
// 		actionTopStories:  {"T"},
// 		actionJobsStories: {"J"},
// 		actionBestStories: {"B"},
// 		actionComments:    {"C"},
// 		actionLoadURL:     {"U"},
// 		actionRefresh:     {"Ctrl+r"},
// 		actionSortByTime:  {"Ctrl+t"},
// 		actionSortByScore: {"Ctrl+s"},
// 	}
// }
