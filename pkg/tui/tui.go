package tui

import (
	"fmt"
	"github.com/Okira-E/patchi/pkg/sequelizer"
	"github.com/Okira-E/patchi/pkg/utils"
	"github.com/gizak/termui/v3"
)

// RenderTui is the entry point for rendering the TUI for Patchi.
func RenderTui(params *GlobalRendererParams) {
	if err := termui.Init(); err != nil {
		utils.Abort(fmt.Sprintf("failed to initialize termui: %v", err))
	}
	defer termui.Close()

	width, height := termui.TerminalDimensions()

	globalRenderer := NewGlobalRenderer(params)

	globalRenderer.ResizeWidgets(width, height)

	globalRenderer.Render()

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

			globalRenderer.Render()
		}

		// - Setup moving from a tab to another.
		if event.Type == termui.KeyboardEvent && (event.ID == "]" || event.ID == "<Right>" || event.ID == "l") {
			globalRenderer.TabPaneWidget.FocusRight()
			globalRenderer.Render()
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "[" || event.ID == "<Left>" || event.ID == "h") {
			globalRenderer.TabPaneWidget.FocusLeft()
			globalRenderer.Render()
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<Tab>") {
			if globalRenderer.FocusedWidget == globalRenderer.DiffWidget {
				globalRenderer.FocusedWidget = globalRenderer.SqlWidget
			} else {
				globalRenderer.FocusedWidget = globalRenderer.DiffWidget
			}

			globalRenderer.ClearBorderStyles()
			globalRenderer.FocusedWidget.BorderStyle = focusedWidgetBorderStyle

			globalRenderer.Render()
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<Enter>") {
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

			}

			globalRenderer.Render()

		}
		if event.Type == termui.KeyboardEvent && ((event.ID == "j") || (event.ID == "<Down>")) {
			focusedWidget := globalRenderer.FocusedWidget

			if len(focusedWidget.Rows) != 0 {
				focusedWidget.ScrollDown()
				globalRenderer.Render()
			}
		}
		if event.Type == termui.KeyboardEvent && ((event.ID == "k") || (event.ID == "<Up>")) {
			focusedWidget := globalRenderer.FocusedWidget

			if len(focusedWidget.Rows) != 0 {
				focusedWidget.ScrollUp()
				globalRenderer.Render()
			}
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<C-d>") {
			if len(globalRenderer.DiffWidget.Rows) != 0 {
				globalRenderer.DiffWidget.ScrollHalfPageDown()
				globalRenderer.Render()
			}
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<C-u>") {
			if len(globalRenderer.DiffWidget.Rows) != 0 {
				globalRenderer.DiffWidget.ScrollHalfPageUp()
				globalRenderer.Render()
			}
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<C-f>") {
			if len(globalRenderer.DiffWidget.Rows) != 0 {
				globalRenderer.DiffWidget.ScrollPageDown()
				globalRenderer.Render()
			}
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<C-b>") {
			if len(globalRenderer.DiffWidget.Rows) != 0 {
				globalRenderer.DiffWidget.ScrollPageUp()
				globalRenderer.Render()
			}
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "g") {
			if event.Type == termui.KeyboardEvent && (event.ID == "g") {
				globalRenderer.DiffWidget.ScrollTop()
				globalRenderer.Render()
			}
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<Home>") {
			if len(globalRenderer.DiffWidget.Rows) != 0 {
				globalRenderer.DiffWidget.ScrollTop()
				globalRenderer.Render()
			}
		}
		if event.Type == termui.KeyboardEvent && ((event.ID == "G") || (event.ID == "<End>")) {
			if len(globalRenderer.DiffWidget.Rows) != 0 {
				globalRenderer.DiffWidget.ScrollBottom()
				globalRenderer.Render()
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
