package tui

import (
	"fmt"
	"github.com/Okira-E/patchi/pkg/migration"
	"github.com/Okira-E/patchi/pkg/safego"
	"github.com/Okira-E/patchi/pkg/types"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type HomeParams struct {
	FirstDb  types.DbConnection
	SecondDb types.DbConnection
}

// HomeRenderer is a unit that knows about all the widgets that need to be rendered.
// rendering the TUI should always be done by calling the `render` method.
type HomeRenderer struct {
	// width and height are the dimensions state of the terminal.
	width, height int

	// TabPaneWidget is the widget that holds the tabs for the different types of objects.
	TabPaneWidget *widgets.TabPane

	// DiffWidget is the widget that holds the diff of the selected tab.
	DiffWidget *widgets.List

	// SqlWidget is the widget that holds the sql formulated from the diff.
	SqlWidget *widgets.Paragraph

	// showConfirmation is a flag that indicates whether the confirmation widget should be shown.
	showConfirmation bool

	params *HomeParams
}

// NewHomeRenderer creates a new instance of CompareRootRenderer.
func NewHomeRenderer(params *HomeParams) *HomeRenderer {
	homeRenderer := &HomeRenderer{
		TabPaneWidget:    widgets.NewTabPane("Tables", "Columns", "Views", "Procedures", "Functions", "Triggers"),
		DiffWidget:       widgets.NewList(),
		SqlWidget:        widgets.NewParagraph(),
		showConfirmation: true,
		params:           params,
	}

	homeRenderer.DiffWidget.TextStyle = termui.NewStyle(termui.ColorYellow)
	homeRenderer.DiffWidget.WrapText = false

	homeRenderer.SqlWidget.Title = "SQL"

	return homeRenderer
}

// ResizeWidgets resizes the widgets based on the new width and height of the terminal.
func (self *HomeRenderer) ResizeWidgets(width int, height int) {
	const tabPaneHeight = 14

	self.TabPaneWidget.SetRect(0, 0, width/2, height/tabPaneHeight)
	self.DiffWidget.SetRect(0, height/tabPaneHeight, width/2, height)
	self.SqlWidget.SetRect(width/2, 0, width, height)

	self.width = width
	self.height = height
}

// SetShowConfirmation sets the showConfirmation flag.
func (self *HomeRenderer) SetShowConfirmation(showConfirmation bool) {
	self.showConfirmation = showConfirmation
}

// GetShowConfirmation returns the showConfirmation flag.
func (self *HomeRenderer) GetShowConfirmation() bool {
	return self.showConfirmation
}

// Render is responsible for rendering the TUI.
func (self *HomeRenderer) Render() {
	termui.Clear()

	if self.TabPaneWidget.ActiveTabIndex == 0 { // Tables

		self.DiffWidget.Title = "Tables"
		self.DiffWidget.Rows = []string{"default"}
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
			diff := migration.TableDiff(&self.params.FirstDb, &self.params.SecondDb)

			for _, tableDiff := range diff {
				self.DiffWidget.Rows = append(self.DiffWidget.Rows, fmt.Sprintf("%s (%s)", tableDiff.TableName, tableDiff.DiffType))
			}
		}

		// Render the widgets.
		if confirmationWidget.IsSome() {
			termui.Render(self.TabPaneWidget, self.DiffWidget, self.SqlWidget, confirmationWidget.Unwrap())
		} else {
			termui.Render(self.TabPaneWidget, self.DiffWidget, self.SqlWidget)
		}

	} else if self.TabPaneWidget.ActiveTabIndex == 1 { // Columns

		self.DiffWidget.Title = "Columns"
		self.DiffWidget.Rows = []string{}
		termui.Render(self.TabPaneWidget, self.DiffWidget, self.SqlWidget)

	} else if self.TabPaneWidget.ActiveTabIndex == 2 { // Views

		self.DiffWidget.Title = "Views"
		self.DiffWidget.Rows = []string{}
		termui.Render(self.TabPaneWidget, self.DiffWidget, self.SqlWidget)

	} else if self.TabPaneWidget.ActiveTabIndex == 3 { // Procedures

		self.DiffWidget.Title = "Procedures"
		self.DiffWidget.Rows = []string{}
		termui.Render(self.TabPaneWidget, self.DiffWidget, self.SqlWidget)

	} else if self.TabPaneWidget.ActiveTabIndex == 4 { // Functions

		self.DiffWidget.Title = "Functions"
		self.DiffWidget.Rows = []string{}
		termui.Render(self.TabPaneWidget, self.DiffWidget, self.SqlWidget)

	} else if self.TabPaneWidget.ActiveTabIndex == 5 { // Triggers

		self.DiffWidget.Title = "Triggers"
		self.DiffWidget.Rows = []string{}
		termui.Render(self.TabPaneWidget, self.DiffWidget, self.SqlWidget)

	}
}
