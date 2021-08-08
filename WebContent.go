package main

//// WebContent is a page that will render the selected item's url
//func WebContent(nextSlide func()) (title string, content tview.Primitive) {
//	app.gui.content = tview.NewTextView()
//	app.gui.content.
//		SetScrollable(true).
//		SetDynamicColors(true).
//		SetBorderPadding(1, 1, 5, 5)
//	return "WebContent", app.gui.content
//}
//
//func (gui *GUI) parseWebContent(item hnapi.Item) {
//	gui.console.SetText("Loading page...")
//	app.main.QueueUpdateDraw(func() {
//		if parsedURL, err := nurl.Parse(item.URL); err != nil {
//			log.Print(err)
//			gui.console.SetText("URL parsing error!")
//		} else {
//			article, err := readability.FromURL(parsedURL.String(), 7*time.Second)
//			if err != nil {
//				log.Print(err)
//				consoleStr := fmt.Sprintf("Content parsing error: %s", err)
//				gui.errorStatus(consoleStr)
//			}
//			gui.content.Write([]byte(article.TextContent))
//			gui.infoStatus("Page successfully loaded.")
//		}
//	})
//}
