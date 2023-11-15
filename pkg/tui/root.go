package tui

import (
	"log"

	"github.com/gizak/termui/v3"
)

// GlobalRenderer starts the TUI for comparing 2 databases.
func GlobalRenderer(params *HomeParams) {
	if err := termui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer termui.Close()

	width, height := termui.TerminalDimensions()

	homeRenderer := NewHomeRenderer(params)

	homeRenderer.ResizeWidgets(width, height)

	homeRenderer.Render()

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

			homeRenderer.ResizeWidgets(width, height)

			homeRenderer.Render()
		}

		// - Setup moving from a tab to another.
		if event.Type == termui.KeyboardEvent && (event.ID == "]" || event.ID == "<Right>" || event.ID == "l") {
			homeRenderer.TabPaneWidget.FocusRight()
			homeRenderer.Render()
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "[" || event.ID == "<Left>" || event.ID == "h") {
			homeRenderer.TabPaneWidget.FocusLeft()
			homeRenderer.Render()
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<Enter>") {
			homeRenderer.SetShowConfirmation(false)
			homeRenderer.Render()
		}
		if event.Type == termui.KeyboardEvent && ((event.ID == "j") || (event.ID == "<Down>")) {
			homeRenderer.DiffWidget.ScrollDown()
			homeRenderer.Render()
		}
		if event.Type == termui.KeyboardEvent && ((event.ID == "k") || (event.ID == "<Up>")) {
			homeRenderer.DiffWidget.ScrollUp()
			homeRenderer.Render()
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<C-d>") {
			homeRenderer.DiffWidget.ScrollHalfPageDown()
			homeRenderer.Render()
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<C-u>") {
			homeRenderer.DiffWidget.ScrollHalfPageUp()
			homeRenderer.Render()
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<C-f>") {
			homeRenderer.DiffWidget.ScrollPageDown()
			homeRenderer.Render()
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<C-b>") {
			homeRenderer.DiffWidget.ScrollPageUp()
			homeRenderer.Render()
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "g") {
			if event.Type == termui.KeyboardEvent && (event.ID == "g") {
				homeRenderer.DiffWidget.ScrollTop()
				homeRenderer.Render()
			}
		}
		if event.Type == termui.KeyboardEvent && (event.ID == "<Home>") {
			homeRenderer.DiffWidget.ScrollTop()
			homeRenderer.Render()
		}
		if event.Type == termui.KeyboardEvent && ((event.ID == "G") || (event.ID == "<End>")) {
			homeRenderer.DiffWidget.ScrollBottom()
			homeRenderer.Render()
		}
	}
}