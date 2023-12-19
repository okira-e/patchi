package tui

import (
	"fmt"

	"github.com/Okira-E/patchi/pkg/tui/patchi_renderer"
	"github.com/Okira-E/patchi/pkg/utils"
	"github.com/Okira-E/patchi/safego"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// FEAT: Make <Ctrl + r> key event that reloads the diff.
// FEAT: Make <Ctrl + s> key event that saves the generated SQL to a file.
// FEAT: Make <Ctrl + e> key event that opens the generated SQL in an editor.

// RenderTui is the entry point for rendering the TUI for Patchi.
// It represents the TUI library and event loop of the application. Actual UI related to the application and its
// functionality can be found in patchiRenderer.go which gives this function the widgets to render.
// Things related to how widgets look and behave can be found in the patchiRenderer.go file. But things related to
// handling events, keyboard events, and the event loop can be found in this file.
func RenderTui(params *patchi_renderer.PatchiRendererParams) {
	if err := termui.Init(); err != nil {
		utils.Abort(fmt.Sprintf("failed to initialize termui: %v", err))
	}
	defer termui.Close()

	width, height := termui.TerminalDimensions()

	patchiRenderer := patchi_renderer.NewPatchiRenderer(params)

	patchiRenderer.ResizeWidgets(width, height)

	patchiRenderer.RenderWidgets(safego.None[string]())

	for event := range termui.PollEvents() {

		if event.Type == termui.KeyboardEvent && (event.ID == "<Escape>") {
			patchiRenderer.ToggleHelpWidget()

			patchiRenderer.RenderWidgets(safego.None[string]())
		}

		if event.Type == termui.KeyboardEvent && (event.ID == "q" || event.ID == "<C-c>") {
			if patchiRenderer.ShowHelpWidget {
				patchiRenderer.ToggleHelpWidget()

				patchiRenderer.RenderWidgets(safego.None[string]())
			} else {
				break // Exit Patchi.
			}
		}

		// Help widget
		if event.Type == termui.KeyboardEvent && (event.ID == "h" || event.ID == "?") {
			patchiRenderer.ToggleHelpWidget()

			patchiRenderer.RenderWidgets(safego.None[string]())
		}

		// Event handling:

		// Change widgets size & location based on resizing the terminal
		if event.Type == termui.ResizeEvent {
			width, height = termui.TerminalDimensions()

			patchiRenderer.ResizeWidgets(width, height)

			patchiRenderer.RenderWidgets(safego.None[string]())
		}

		// Setup moving from a tab to another.
		if event.Type == termui.KeyboardEvent && (event.ID == "]" || event.ID == "<Right>" || event.ID == "l") {
			// FEAT: Going to the right of the last option should bring you back. The opposite is true.
			patchiRenderer.TabPaneWidget.FocusRight()
			patchiRenderer.ResetMsgBar()
			patchiRenderer.RenderWidgets(safego.None[string]())
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "[" || event.ID == "<Left>" || event.ID == "h") {
			patchiRenderer.TabPaneWidget.FocusLeft()
			patchiRenderer.ResetMsgBar()
			patchiRenderer.RenderWidgets(safego.None[string]())
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<Tab>") {
			if patchiRenderer.FocusedWidget == patchiRenderer.DiffWidget {
				patchiRenderer.FocusedWidget = patchiRenderer.SqlWidget
			} else {
				patchiRenderer.FocusedWidget = patchiRenderer.DiffWidget
			}

			patchiRenderer.RenderWidgets(safego.None[string]())
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<Enter>") {
			patchiRenderer.HandleActionOnEnter()
		}
		// FEAT: Add "y" and make it yank the generated SQL also
		if event.Type == termui.KeyboardEvent && ((event.ID == "j") || (event.ID == "<Down>")) {
			switch patchiRenderer.FocusedWidget.(type) {
			case *widgets.List:
				if len(patchiRenderer.FocusedWidget.(*widgets.List).Rows) != 0 {
					patchiRenderer.FocusedWidget.(*widgets.List).ScrollDown()
					patchiRenderer.RenderWidgets(safego.None[string]())
				}
			}
		}
		if event.Type == termui.KeyboardEvent && ((event.ID == "k") || (event.ID == "<Up>")) {
			switch patchiRenderer.FocusedWidget.(type) {
			case *widgets.List:
				if len(patchiRenderer.FocusedWidget.(*widgets.List).Rows) != 0 {
					patchiRenderer.FocusedWidget.(*widgets.List).ScrollUp()
					patchiRenderer.RenderWidgets(safego.None[string]())
				}
			}
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<C-d>") {
			switch patchiRenderer.FocusedWidget.(type) {
			case *widgets.List:
				if len(patchiRenderer.FocusedWidget.(*widgets.List).Rows) != 0 {
					patchiRenderer.FocusedWidget.(*widgets.List).ScrollHalfPageDown()
					patchiRenderer.RenderWidgets(safego.None[string]())
				}
			}
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<C-u>") {
			switch patchiRenderer.FocusedWidget.(type) {
			case *widgets.List:
				if len(patchiRenderer.FocusedWidget.(*widgets.List).Rows) != 0 {
					patchiRenderer.FocusedWidget.(*widgets.List).ScrollHalfPageUp()
					patchiRenderer.RenderWidgets(safego.None[string]())
				}
			}
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "g") || (event.ID == "<Home>") {
			switch patchiRenderer.FocusedWidget.(type) {
			case *widgets.List:
				if len(patchiRenderer.FocusedWidget.(*widgets.List).Rows) != 0 {
					patchiRenderer.FocusedWidget.(*widgets.List).ScrollTop()
					patchiRenderer.RenderWidgets(safego.None[string]())
				}
			}
		}
		if event.Type == termui.KeyboardEvent && ((event.ID == "G") || (event.ID == "<End>")) {
			if len(patchiRenderer.DiffWidget.Rows) != 0 {
				switch patchiRenderer.FocusedWidget.(type) {
				case *widgets.List:
					if len(patchiRenderer.FocusedWidget.(*widgets.List).Rows) != 0 {
						patchiRenderer.FocusedWidget.(*widgets.List).ScrollBottom()
						patchiRenderer.RenderWidgets(safego.None[string]())
					}
				}
			}
		}
	}
}
