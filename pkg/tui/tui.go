package tui

import (
	"fmt"
	"github.com/Okira-E/patchi/pkg/sequelizer"
	"github.com/Okira-E/patchi/pkg/utils"
	"github.com/Okira-E/patchi/safego"
	"github.com/atotto/clipboard"
	"github.com/gizak/termui/v3"
)

// TODO: Make a refresh key event that refreshes the diff.

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
		// Keys mappings:

		// - Exiting
		if event.Type == termui.KeyboardEvent && (event.ID == "q" || event.ID == "<C-c>") {
			break
		}

		// - Help widget
		if event.Type == termui.KeyboardEvent && (event.ID == "h" || event.ID == "?") {
			// TODO: Show help widget.
		}

		// Event handling:

		// - Change widgets size & location based on resizing the terminal
		if event.Type == termui.ResizeEvent {
			width, height = termui.TerminalDimensions()

			globalRenderer.ResizeWidgets(width, height)

			globalRenderer.Render(safego.None[string]())
		}

		// - Setup moving from a tab to another.
		if event.Type == termui.KeyboardEvent && (event.ID == "]" || event.ID == "<Right>" || event.ID == "l") {
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

			globalRenderer.ClearBorderStyles()
			globalRenderer.FocusedWidget.BorderStyle = focusedWidgetBorderStyle

			globalRenderer.Render(safego.None[string]())
		}
		// Copy all generated SQL to clipboard.
		if event.Type == termui.KeyboardEvent && (event.ID == "<C-a>") && (globalRenderer.FocusedWidget == globalRenderer.SqlWidget) {
			userPrompt := safego.None[string]()

			if len(globalRenderer.SqlWidget.Rows) == 0 { // List of entities is empty.
				continue
			}

			netSql := ""
			for _, sql := range globalRenderer.SqlWidget.Rows {
				netSql += sql
			}

			err := clipboard.WriteAll(netSql)
			if err != nil {
				userPrompt = safego.Some("[Error copying to clipboard: " + err.Error() + "](fg:red)")
			} else {
				userPrompt = safego.Some("[Copied entire generated SQL to clipboard.](fg:green)")
			}

			globalRenderer.Render(userPrompt)
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<Enter>") {
			errPrompt := safego.None[string]()

			if globalRenderer.showConfirmation {
				globalRenderer.SetShowConfirmation(false)
			} else if globalRenderer.FocusedWidget == globalRenderer.DiffWidget {
				if len(globalRenderer.DiffWidget.Rows) == 0 { // List of entities is empty.
					continue
				}

				// Handle the case where the user presses enter on an empty row if at all possible.
				if len(utils.ExtractExpressions(globalRenderer.DiffWidget.Rows[globalRenderer.DiffWidget.SelectedRow], "\\[(.*?)\\]")) == 0 {
					continue
				}

				currentlySelectedTab := getCurrentTabName(globalRenderer.TabPaneWidget.ActiveTabIndex)
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

				globalRenderer.SqlWidget.Rows = append(globalRenderer.SqlWidget.Rows, sql+"\n")
			} else if globalRenderer.FocusedWidget == globalRenderer.SqlWidget {
				if len(globalRenderer.SqlWidget.Rows) == 0 { // List of entities is empty.
					continue
				}

				selectedSqlForCopy := globalRenderer.SqlWidget.Rows[globalRenderer.SqlWidget.SelectedRow]

				err := clipboard.WriteAll(selectedSqlForCopy)
				if err != nil {
					errPrompt = safego.Some("[Error copying to clipboard: " + err.Error() + "](fg:red)")
				} else {
					errPrompt = safego.Some("[Copied selected SQL to clipboard.](fg:green)")
				}
			}

			globalRenderer.Render(errPrompt)

		}
		if event.Type == termui.KeyboardEvent && ((event.ID == "j") || (event.ID == "<Down>")) {
			focusedWidget := globalRenderer.FocusedWidget

			if len(focusedWidget.Rows) != 0 {
				focusedWidget.ScrollDown()
				globalRenderer.Render(safego.None[string]())
			}
		}
		if event.Type == termui.KeyboardEvent && ((event.ID == "k") || (event.ID == "<Up>")) {
			focusedWidget := globalRenderer.FocusedWidget

			if len(focusedWidget.Rows) != 0 {
				focusedWidget.ScrollUp()
				globalRenderer.Render(safego.None[string]())
			}
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<C-d>") {
			if len(globalRenderer.DiffWidget.Rows) != 0 {
				globalRenderer.DiffWidget.ScrollHalfPageDown()
				globalRenderer.Render(safego.None[string]())
			}
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<C-u>") {
			if len(globalRenderer.DiffWidget.Rows) != 0 {
				globalRenderer.DiffWidget.ScrollHalfPageUp()
				globalRenderer.Render(safego.None[string]())
			}
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<C-f>") {
			if len(globalRenderer.DiffWidget.Rows) != 0 {
				globalRenderer.DiffWidget.ScrollPageDown()
				globalRenderer.Render(safego.None[string]())
			}
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<C-b>") {
			if len(globalRenderer.DiffWidget.Rows) != 0 {
				globalRenderer.DiffWidget.ScrollPageUp()
				globalRenderer.Render(safego.None[string]())
			}
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "g") {
			if event.Type == termui.KeyboardEvent && (event.ID == "g") {
				globalRenderer.DiffWidget.ScrollTop()
				globalRenderer.Render(safego.None[string]())
			}
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<Home>") {
			if len(globalRenderer.DiffWidget.Rows) != 0 {
				globalRenderer.DiffWidget.ScrollTop()
				globalRenderer.Render(safego.None[string]())
			}
		}
		if event.Type == termui.KeyboardEvent && ((event.ID == "G") || (event.ID == "<End>")) {
			if len(globalRenderer.DiffWidget.Rows) != 0 {
				globalRenderer.DiffWidget.ScrollBottom()
				globalRenderer.Render(safego.None[string]())
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