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

	// showConfirmation is a flag that indicates whether the confirmation widget should be shown.
	showConfirmation bool

	params *GlobalRendererParams
}

// NewGlobalRenderer creates a new instance of CompareRootRenderer.
func NewGlobalRenderer(params *GlobalRendererParams) *GlobalRenderer {
	globalRenderer := &GlobalRenderer{
		TabPaneWidget:    widgets.NewTabPane("Tables", "Columns", "Views", "Procedures", "Functions", "Triggers"),
		DiffWidget:       widgets.NewList(),
		SqlWidget:        widgets.NewList(),
		MessageBarWidget: widgets.NewParagraph(),
		showConfirmation: true,
		params:           params,
	}

	globalRenderer.FocusedWidget = globalRenderer.DiffWidget
	globalRenderer.FocusedWidget.BorderStyle = focusedWidgetBorderStyle

	globalRenderer.DiffWidget.TextStyle = termui.NewStyle(termui.ColorWhite)
	globalRenderer.DiffWidget.SelectedRowStyle = termui.NewStyle(termui.ColorBlack, termui.ColorWhite)
	globalRenderer.DiffWidget.WrapText = false
	globalRenderer.DiffWidget.PaddingLeft = 2

	globalRenderer.SqlWidget.Title = "SQL"
	globalRenderer.SqlWidget.TextStyle = termui.NewStyle(termui.ColorWhite)
	globalRenderer.SqlWidget.SelectedRowStyle = termui.NewStyle(termui.ColorGreen)
	globalRenderer.SqlWidget.WrapText = false
	globalRenderer.SqlWidget.PaddingLeft = 2

	globalRenderer.MessageBarWidget.Border = false
	globalRenderer.MessageBarWidget.Text = "Press Enter to fetch changes."

	return globalRenderer
}

// ResizeWidgets resizes the widgets based on the new width and height of the terminal.
func (self *GlobalRenderer) ResizeWidgets(width int, height int) {
	const tabPaneHeight = 14

	self.TabPaneWidget.SetRect(0, 0, width/2, height/tabPaneHeight)
	self.DiffWidget.SetRect(0, height/tabPaneHeight, width/2, height-(height/35))
	self.SqlWidget.SetRect(width/2, 0, width, height-(height/35))
	self.MessageBarWidget.SetRect(0, height-(height/35), width, height)

	self.width = width
	self.height = height
}

// SetShowConfirmation sets the showConfirmation flag.
func (self *GlobalRenderer) SetShowConfirmation(showConfirmation bool) {
	self.showConfirmation = showConfirmation
}

// GetShowConfirmation returns the showConfirmation flag.
func (self *GlobalRenderer) GetShowConfirmation() bool {
	return self.showConfirmation
}

// ClearBorderStyles clears the border styles of all the widgets.
func (self *GlobalRenderer) ClearBorderStyles() {
	self.DiffWidget.BorderStyle = termui.NewStyle(termui.ColorClear)
	self.SqlWidget.BorderStyle = termui.NewStyle(termui.ColorClear)
}

// Render is responsible for rendering the TUI.
func (self *GlobalRenderer) Render() {
	termui.Clear()

	if self.TabPaneWidget.ActiveTabIndex == 0 { // Tables

		self.DiffWidget.Title = "Tables"
		self.DiffWidget.Rows = []string{}
		diffWidgetRec := self.DiffWidget.GetRect()

		// Check if the user has started the comparing process or not yet. Confirmation widget is shown instead
		// of the diff widget if not.
		confirmationWidget := safego.None[*widgets.Paragraph]()
		if self.showConfirmation {
			confirmationWidget = safego.Some(widgets.NewParagraph())
			confirmationWidget.Unwrap().SetRect(diffWidgetRec.Min.X, diffWidgetRec.Min.Y+1, diffWidgetRec.Max.X, diffWidgetRec.Max.Y)
			confirmationWidget.Unwrap().BorderTop = false

			confirmationWidget.Unwrap().Text = "Press Enter to fetch changes."
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

		// Render the widgets.
		if confirmationWidget.IsSome() {
			termui.Render(self.TabPaneWidget, self.DiffWidget, self.SqlWidget, self.MessageBarWidget, confirmationWidget.Unwrap())
		} else {
			termui.Render(self.TabPaneWidget, self.DiffWidget, self.SqlWidget, self.MessageBarWidget)
		}

	} else if self.TabPaneWidget.ActiveTabIndex == 1 { // Columns

		self.DiffWidget.Title = "Columns"
		self.DiffWidget.Rows = []string{}
		termui.Render(self.TabPaneWidget, self.DiffWidget, self.SqlWidget, self.MessageBarWidget)

	} else if self.TabPaneWidget.ActiveTabIndex == 2 { // Views

		self.DiffWidget.Title = "Views"
		self.DiffWidget.Rows = []string{}
		termui.Render(self.TabPaneWidget, self.DiffWidget, self.SqlWidget, self.MessageBarWidget)

	} else if self.TabPaneWidget.ActiveTabIndex == 3 { // Procedures

		self.DiffWidget.Title = "Procedures"
		self.DiffWidget.Rows = []string{}
		termui.Render(self.TabPaneWidget, self.DiffWidget, self.SqlWidget, self.MessageBarWidget)

	} else if self.TabPaneWidget.ActiveTabIndex == 4 { // Functions

		self.DiffWidget.Title = "Functions"
		self.DiffWidget.Rows = []string{}
		termui.Render(self.TabPaneWidget, self.DiffWidget, self.SqlWidget, self.MessageBarWidget)

	} else if self.TabPaneWidget.ActiveTabIndex == 5 { // Triggers

		self.DiffWidget.Title = "Triggers"
		self.DiffWidget.Rows = []string{}
		termui.Render(self.TabPaneWidget, self.DiffWidget, self.SqlWidget, self.MessageBarWidget)

	}
}
