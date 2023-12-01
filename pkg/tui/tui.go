package tui

import (
	"fmt"

	"github.com/Okira-E/patchi/pkg/sequelizer"
	"github.com/Okira-E/patchi/pkg/utils"
	"github.com/Okira-E/patchi/safego"
	"github.com/atotto/clipboard"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// FEAT: Make <Ctrl + r> key event that reloads the diff.
// FEAT: Make <Ctrl + s> key event that saves the generated SQL to a file.
// FEAT: Make <Ctrl + e> key event that opens the generated SQL in an editor.

// RenderTui is the entry point for rendering the TUI for Patchi.
// It represents the TUI library and event loop of the application. Actual UI related to the application and its
// functionality can be found in `globalRenderer.go` which gives this function the widgets to render.
// Things related to how widgets look and behave can be found in the `globalRenderer.go` file. But things related to
// handling events, keyboard events, and the event loop can be found in this file.
func RenderTui(params *GlobalRendererParams) {
	if err := termui.Init(); err != nil {
		utils.Abort(fmt.Sprintf("failed to initialize termui: %v", err))
	}
	defer termui.Close()

	width, height := termui.TerminalDimensions()

	globalRenderer := NewGlobalRenderer(params)

	globalRenderer.ResizeWidgets(width, height)

	globalRenderer.Render(safego.None[string]())

	for event := range termui.PollEvents() {

		if event.Type == termui.KeyboardEvent && (event.ID == "<Escape>") {
			globalRenderer.ToggleHelpWidget()

			globalRenderer.Render(safego.None[string]())
		}

		if event.Type == termui.KeyboardEvent && (event.ID == "q" || event.ID == "<C-c>") {
			if globalRenderer.ShowHelpWidget {
				globalRenderer.ToggleHelpWidget()

				globalRenderer.Render(safego.None[string]())
			} else {
				break // Exit Patchi.
			}
		}

		// Help widget
		if event.Type == termui.KeyboardEvent && (event.ID == "h" || event.ID == "?") {
			globalRenderer.ToggleHelpWidget()

			globalRenderer.Render(safego.None[string]())
		}

		// Event handling:

		// Change widgets size & location based on resizing the terminal
		if event.Type == termui.ResizeEvent {
			width, height = termui.TerminalDimensions()

			globalRenderer.ResizeWidgets(width, height)

			globalRenderer.Render(safego.None[string]())
		}

		// Setup moving from a tab to another.
		if event.Type == termui.KeyboardEvent && (event.ID == "]" || event.ID == "<Right>" || event.ID == "l") {
			// FEAT: Going to the right of the last option should bring you back. The opposite is true.
			globalRenderer.TabPaneWidget.FocusRight()
			globalRenderer.Render(safego.None[string]())
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "[" || event.ID == "<Left>" || event.ID == "h") {
			globalRenderer.TabPaneWidget.FocusLeft()
			globalRenderer.Render(safego.None[string]())
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<Tab>") {
			if globalRenderer.FocusedWidget == globalRenderer.DiffWidget {
				globalRenderer.FocusedWidget = globalRenderer.SqlWidget
			} else {
				globalRenderer.FocusedWidget = globalRenderer.DiffWidget
			}

			globalRenderer.Render(safego.None[string]())
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<Enter>") {
			// Responsible for:
			// - Diff widget -> Generates the unsynced entrie names (tables, columns, views, ..etc)
			// - SQL widget -> Copies the generated SQL.
			optErrPrompt := safego.None[string]()

			if globalRenderer.ShowConfirmation {
				globalRenderer.ShowConfirmation = false
			} else if globalRenderer.FocusedWidget == globalRenderer.DiffWidget {
				if len(globalRenderer.DiffWidget.Rows) == 0 { // List of entities is empty.
					continue
				}

				// Handle the case where the user presses enter on an empty row if at all possible.
				if len(utils.ExtractExpressions(globalRenderer.DiffWidget.Rows[globalRenderer.DiffWidget.SelectedRow], "\\[(.*?)\\]")) == 0 {
					continue
				}

				currentlySelectedTab := getCurrentTabName(globalRenderer.TabPaneWidget.ActiveTabIndex)
				// FIX: Don't generate SQL for entities twice. Use a map[entityName]bool to not generate the same SQL twice.
				currentlySelectedEntityName := utils.ExtractExpressions(globalRenderer.DiffWidget.Rows[globalRenderer.DiffWidget.SelectedRow], "\\[(.*?)\\]")[0]

				// This is the status of the entity (created, deleted, modified)
				entityStatus := utils.ExtractExpressions(globalRenderer.DiffWidget.Rows[globalRenderer.DiffWidget.SelectedRow], "fg:(.*?)\\)")[0]
				if entityStatus == "green" {
					entityStatus = "created"
				} else if entityStatus == "red" {
					entityStatus = "deleted"
				} else {
					entityStatus = "modified"
				}

				sql := sequelizer.PatchSqlForEntity(params.FirstDb, params.SecondDb, currentlySelectedTab, currentlySelectedEntityName, entityStatus)

				globalRenderer.SqlWidget.Text += sql + "\n\n"
			} else if globalRenderer.FocusedWidget == globalRenderer.SqlWidget {
				if globalRenderer.SqlWidget.Text == "" { // List of entities is empty.
					continue
				}

				sql := globalRenderer.SqlWidget.Text

				err := clipboard.WriteAll(sql)
				if err != nil {
					optErrPrompt = safego.Some("[Error copying to clipboard: " + err.Error() + "](fg:red)")
				} else {
					optErrPrompt = safego.Some("[Copied entire generated SQL to clipboard.](fg:green)")
				}
			}

			globalRenderer.Render(optErrPrompt)

		}
		// FEAT: Add "y" and make it yank the generated SQL also
		if event.Type == termui.KeyboardEvent && ((event.ID == "j") || (event.ID == "<Down>")) {
			switch globalRenderer.FocusedWidget.(type) {
			case *widgets.List:
				if len(globalRenderer.FocusedWidget.(*widgets.List).Rows) != 0 {
					globalRenderer.FocusedWidget.(*widgets.List).ScrollDown()
					globalRenderer.Render(safego.None[string]())
				}
			}
		}
		if event.Type == termui.KeyboardEvent && ((event.ID == "k") || (event.ID == "<Up>")) {
			switch globalRenderer.FocusedWidget.(type) {
			case *widgets.List:
				if len(globalRenderer.FocusedWidget.(*widgets.List).Rows) != 0 {
					globalRenderer.FocusedWidget.(*widgets.List).ScrollUp()
					globalRenderer.Render(safego.None[string]())
				}
			}
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<C-d>") {
			switch globalRenderer.FocusedWidget.(type) {
			case *widgets.List:
				if len(globalRenderer.FocusedWidget.(*widgets.List).Rows) != 0 {
					globalRenderer.FocusedWidget.(*widgets.List).ScrollHalfPageDown()
					globalRenderer.Render(safego.None[string]())
				}
			}
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<C-u>") {
			switch globalRenderer.FocusedWidget.(type) {
			case *widgets.List:
				if len(globalRenderer.FocusedWidget.(*widgets.List).Rows) != 0 {
					globalRenderer.FocusedWidget.(*widgets.List).ScrollHalfPageUp()
					globalRenderer.Render(safego.None[string]())
				}
			}
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "g") || (event.ID == "<Home>") {
			switch globalRenderer.FocusedWidget.(type) {
			case *widgets.List:
				if len(globalRenderer.FocusedWidget.(*widgets.List).Rows) != 0 {
					globalRenderer.FocusedWidget.(*widgets.List).ScrollTop()
					globalRenderer.Render(safego.None[string]())
				}
			}
		}
		if event.Type == termui.KeyboardEvent && ((event.ID == "G") || (event.ID == "<End>")) {
			if len(globalRenderer.DiffWidget.Rows) != 0 {
				switch globalRenderer.FocusedWidget.(type) {
				case *widgets.List:
					if len(globalRenderer.FocusedWidget.(*widgets.List).Rows) != 0 {
						globalRenderer.FocusedWidget.(*widgets.List).ScrollBottom()
						globalRenderer.Render(safego.None[string]())
					}
				}
			}
		}
	}
}

func getCurrentTabName(index int) string {
	if index == 0 {
		return "tables"
	} else if index == 1 {
		return "columns"
	} else if index == 2 {
		return "views"
	} else if index == 3 {
		return "procedures"
	} else if index == 4 {
		return "functions"
	} else if index == 5 {
		return "triggers"
	}

	return ""
}
