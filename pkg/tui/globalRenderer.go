package tui

import (
	"strconv"
	"strings"

	"github.com/Okira-E/patchi/pkg/diff"
	"github.com/Okira-E/patchi/pkg/types"
	"github.com/Okira-E/patchi/safego"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var focusedWidgetBorderStyle = termui.NewStyle(termui.ColorGreen)

type GlobalRendererParams struct {
	FirstDb  types.DbConnection
	SecondDb types.DbConnection
}

// GlobalRenderer is a unit that knows about all the widgets that need to be rendered.
// rendering the TUI should always be done by calling the `render` method.
type GlobalRenderer struct {
	// width and height are the dimensions state of the terminal.
	width, height int

	// TabPaneWidget is the widget that holds the tabs for the different types of objects.
	TabPaneWidget *widgets.TabPane

	// DiffWidget is the widget that holds the diff of the selected tab.
	DiffWidget *widgets.List

	// SqlWidget is the widget that holds the sql formulated from the diff.
	SqlWidget *widgets.List

	// MessageBarWidget is the widget that holds the messages that are shown to the user.
	MessageBarWidget *widgets.Paragraph

	// focusedWidget is the widget that is currently being controlled by the user.
	FocusedWidget *widgets.List

	// HelpWidget shows all the keyboard shortcuts for the app.
	HelpWidget *widgets.List

	// confirmationWidget is the widget that is shown when the user has not yet started the comparing process.
	confirmationWidget *widgets.Paragraph

	// ShowHelpWidget determines if `Render()` needs to render the default helping widget.
	ShowHelpWidget bool

	// ShowConfirmation is a flag that indicates whether the confirmation widget should be shown.
	ShowConfirmation bool

	params *GlobalRendererParams
}

// NewGlobalRenderer creates a new instance of CompareRootRenderer.
func NewGlobalRenderer(params *GlobalRendererParams) *GlobalRenderer {
	globalRenderer := &GlobalRenderer{
		TabPaneWidget:      widgets.NewTabPane("Tables", "Columns", "Views", "Procedures", "Functions", "Triggers"),
		DiffWidget:         widgets.NewList(),
		SqlWidget:          widgets.NewList(),
		MessageBarWidget:   widgets.NewParagraph(),
		HelpWidget:         widgets.NewList(),
		confirmationWidget: widgets.NewParagraph(),
		ShowConfirmation:   true,
		params:             params,
	}

	globalRenderer.FocusedWidget = globalRenderer.DiffWidget

	globalRenderer.DiffWidget.TextStyle = termui.NewStyle(termui.ColorWhite)
	globalRenderer.DiffWidget.SelectedRowStyle = termui.NewStyle(termui.ColorBlack, termui.ColorWhite)
	globalRenderer.DiffWidget.WrapText = false
	globalRenderer.DiffWidget.PaddingLeft = 2
	globalRenderer.DiffWidget.Rows = []string{}

	globalRenderer.SqlWidget.Title = "SQL"
	globalRenderer.SqlWidget.TextStyle = termui.NewStyle(termui.ColorWhite)
	globalRenderer.SqlWidget.SelectedRowStyle = termui.NewStyle(termui.ColorGreen)
	globalRenderer.SqlWidget.WrapText = false
	globalRenderer.SqlWidget.PaddingLeft = 2

	globalRenderer.MessageBarWidget.Border = false
	globalRenderer.MessageBarWidget.Text = `Press <h> or <?> for help.`

	globalRenderer.HelpWidget.Title = "Help"
	globalRenderer.HelpWidget.SelectedRowStyle = termui.NewStyle(termui.ColorBlack, termui.ColorWhite)
	globalRenderer.HelpWidget.Rows = []string{
		`[<h | ?>](fg:green)` + "\t \t \t \t \t \t \t \t to show this help widget.",
		`[<Escape>](fg:green)` + "\t \t \t \t \t \t \t to exit help.",
		`[<q | Ctrl+c>](fg:green)` + "\t \t \t to quit.",
		`<[ | Left>` + "\t \t \t \t \t to move to the previous tab.",
		`<] | Right>` + "\t \t \t \t to move to the next tab.",
		`[<Tab>](fg:green)` + "\t \t \t \t \t \t \t \t \t \t to move between the diff and sql widgets.",
		`[<Enter>](fg:green)` + "\t \t \t \t \t \t \t \t on the SQL widget to copy the selected SQL to the clipboard.",
		`[<Ctrl+a>](fg:green)` + "\t \t \t \t \t \t \t on the SQL widget to copy all SQL to the clipboard.",
	}

	globalRenderer.confirmationWidget.BorderTop = false

	globalRenderer.confirmationWidget.Text = "Press Enter to fetch changes."

	return globalRenderer
}

// ResizeWidgets resizes the widgets based on the new width and height of the terminal.
// This method should be called whenever the terminal is resized.
// It is also what sets the original sizing for most of the widgets. I say almost because widgets with conditional
// rendering like the help widget only get their sizing set when they are rendered. They can all be found at the end
// of the `Render(..)` method.
func (self *GlobalRenderer) ResizeWidgets(width int, height int) {
	const tabPaneHeight = 14

	self.TabPaneWidget.SetRect(0, 0, width/2, height/tabPaneHeight) // BUG: Put a limit on the height of the tab pane.
	self.DiffWidget.SetRect(0, height/tabPaneHeight, width/2, height-(height/35))
	self.SqlWidget.SetRect(width/2, 0, width, height-(height/35))
	self.MessageBarWidget.SetRect(0, height-(height/35), width, height)

	self.width = width
	self.height = height
}

// ClearBorderStyles clears the border styles of all the widgets.
func (self *GlobalRenderer) ClearBorderStyles() {
	self.DiffWidget.BorderStyle = termui.NewStyle(termui.ColorClear)
	self.SqlWidget.BorderStyle = termui.NewStyle(termui.ColorClear)
}

func (self *GlobalRenderer) ToggleHelpWidget() {
	if self.ShowHelpWidget {
		self.ShowHelpWidget = false
		self.FocusedWidget = self.DiffWidget // BUG: This should be the widget that was focused before the help widget was shown.

		self.MessageBarWidget.Text = `Press <h> or <?> for help.`
	} else {
		self.ShowHelpWidget = true
		self.FocusedWidget = self.HelpWidget

		self.MessageBarWidget.Text = `Press <Escape> to exit help.`
	}
}

// Render is responsible for rendering the TUI.
//
// - param `userPrompt` is an optional string that is shown in the message bar.
func (self *GlobalRenderer) Render(userPrompt safego.Option[string]) {
	termui.Clear()

	// BUG: DiffWidget.Rows persists between tabs. I don't want to clear it on each tab switch because
	// I want the user to go back to the previous tab and see the diff that was there before. So we should instead
	// have a map for each tab and store the diff there. The map should be keyed by the tab index.

	if self.TabPaneWidget.ActiveTabIndex == 0 { // Tables

		self.DiffWidget.Title = "Tables"
		if len(self.DiffWidget.Rows) == 0 {

			// Check if the user has started the comparing process or not yet. Confirmation widget is shown instead
			// of the diff widget if not.
			if self.ShowConfirmation { // If `ShowConfirmation` is false then `SetRect()` will not be called on it and it will not be rendered.
			} else {
				diffResult := diff.GetDiffInTablesBetweenSchemas(&self.params.FirstDb, &self.params.SecondDb)

				self.MessageBarWidget.Text = "Found " + strconv.Itoa(len(diffResult)) + " changes in " + strings.ToLower(self.DiffWidget.Title) + "."

				for _, tableDiff := range diffResult {
					text := "[" + tableDiff.TableName + "]"
					if tableDiff.DiffType == "created" {
						text += "(fg:green)"
					} else if tableDiff.DiffType == "deleted" {
						text += "(fg:red)"
					}

					self.DiffWidget.Rows = append(self.DiffWidget.Rows, text)
				}
			}
		}

	} else if self.TabPaneWidget.ActiveTabIndex == 1 { // Columns

		self.DiffWidget.Title = "Columns"

	} else if self.TabPaneWidget.ActiveTabIndex == 2 { // Views

		self.DiffWidget.Title = "Views"

	} else if self.TabPaneWidget.ActiveTabIndex == 3 { // Procedures

		self.DiffWidget.Title = "Procedures"

	} else if self.TabPaneWidget.ActiveTabIndex == 4 { // Functions

		self.DiffWidget.Title = "Functions"

	} else if self.TabPaneWidget.ActiveTabIndex == 5 { // Triggers

		self.DiffWidget.Title = "Triggers"

	}

	// Set the currently focused widget's border style to be green.
	self.ClearBorderStyles()
	self.FocusedWidget.BorderStyle = focusedWidgetBorderStyle

	if userPrompt.IsSome() {
		self.MessageBarWidget.Text = userPrompt.Unwrap()
	}

	// You can't pass nil to termui.Render() so we have to set the rect of any conditional widget to 0, 0, 0, 0 if we
	// wish to not render/hide a widget.
	if self.ShowHelpWidget {
		self.HelpWidget.SetRect(
			self.width/4,
			self.height/3,
			self.width-self.width/4,
			self.height-self.height/4,
		)
	} else {
		self.HelpWidget.SetRect(0, 0, 0, 0)
	}
	if self.ShowConfirmation {
		diffWidgetRec := self.DiffWidget.GetRect()
		self.confirmationWidget.SetRect(diffWidgetRec.Min.X, diffWidgetRec.Min.Y+1, diffWidgetRec.Max.X, diffWidgetRec.Max.Y)
	} else {
		self.confirmationWidget.SetRect(0, 0, 0, 0)
	}

	// Render the widgets.
	termui.Render(
		self.TabPaneWidget,
		self.DiffWidget,
		self.SqlWidget,
		self.MessageBarWidget,
		self.confirmationWidget,
		self.HelpWidget,
	)
}
