package ui

import (
	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
)

type tabMode int

const (
  tabModeDone tabMode = iota
  tabModeLoading
)

type tab struct {
  list *cview.List,
  mode tabMode
  tabLabel string
  tabText string
}

func (u *ui) makeNewTab(list *cview.List, label,text string) *tab {
  t := tab{
    list: list,
    mode: tabModeDone,
    tabLabel: label,
    tabText: text,
  }
  t.list.SetSelectedTextColor(HNOrange)
  t.list.SetChangedFunc(
     u.Draw() 
  )
  t.list.SetDoneFunc(func(ket tcell.Key) {
    tab := curTab
    if tabs[tab].mode != tabModeDone {
      return
    }

    if key == tcell.KeyEsc {
      // Do something
      return
    }
  })
}
