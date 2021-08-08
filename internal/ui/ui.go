package ui

import (
	"context"
	"log"

	"code.rocketnine.space/tslocum/cview"
	"github.com/CaninoDev/hackernewsterm/internal/hackernews"
)

type uiConfig struct {
	Input       *map[string][]string
	Theme       *cview.Theme
	MouseEnable bool
	Browser     string
	TermMux     string
	TermRow     uint
	TermHeight  uint
}

type Service interface {
}

type ui struct {
	*cview.Application
	firebase hackernews.Firebase
	cfg uiConfig
	ctx context.Context
}

var (
	app    *cview.Application
	config = &uiConfig{
		MouseEnable: true,
		Browser:     "/usr/bin/lynx",
		TermMux:     "/usr/bin/tmux",
	}
)

// Slide is a function the returns the slide's title, any pertinent information,
// and its main primitive. It recieves a "nextSlide" function which can be
// called to advance the main panel to the next view.
type Slide func(nextSlide func()) (title string, info string, content cview.Primitive)

func InitUI(ctx context.Context, handler hackernews.Firebase) error {

	app = cview.NewApplication()

	app.EnableMouse(config.MouseEnable)

	// app.SetAfterResizeFunc(handleResize)
	// app.SetMouseCapture(handleMouse)
	// app.SetInputCapture(handleInput)
	app.SetBeforeFocusFunc(handleBeforeFocus)

	tui := &ui{
		app, handler, *config, ctx,
	}
	tui.ListView()
	if err := tui.Run(); err != nil {
		log.Fatalf("error initializing tui")
	}
	return nil
}

func handleBeforeFocus(primitive cview.Primitive) bool {
	return false
}
