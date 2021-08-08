package ui

//type tabMode int
//
//const (
//  tabModeDone tabMode = iota
//  tabModeLoading
//)
//
//type tab struct {
//  list *cview.List
//  mode tabMode
//  tabLabel string
//  tabText string
//}
//
//func (u *ui) makeNewTab(list *cview.List, label,text string) *tab {
//  t := tab{
//    list: list,
//    mode: tabModeDone,
//    tabLabel: label,
//    tabText: text,
//  }
//  t.list.SetSelectedTextColor(HNOrange)
//  t.list.SetChangedFunc(
//     u.Draw()
//  )
//}
