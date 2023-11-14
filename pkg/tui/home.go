package tui

import (
	"database/sql"

	"github.com/Okira-E/patchi/pkg/migration"
	"github.com/Okira-E/patchi/pkg/safego"
	"github.com/Okira-E/patchi/pkg/types"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// DbConnectionInfo is the connection info for a database.
type DbConnectionInfo struct {
	Info          *types.DbConnectionInfo
	SqlConnection *sql.DB
}

type HomeParams struct {
	FirstDb  DbConnectionInfo
	SecondDb DbConnectionInfo
}

// HomeRenderer is a unit that knows about all the widgets that need to be rendered.
// rendering the TUI should always be done by calling the `render` method.
type HomeRenderer struct {
	// width and height are the dimensions state of the terminal.
	width, height int

	// TabPaneWidget is the widget that holds the tabs for the different types of objects.
	TabPaneWidget *widgets.TabPane

	// DiffWidget is the widget that holds the diff of the selected tab.
	DiffWidget *widgets.Paragraph

	// SqlWidget is the widget that holds the sql formulated from the diff.
	SqlWidget *widgets.Paragraph

	// showConfirmation is a flag that indicates whether the confirmation widget should be shown.
	showConfirmation bool

	// params is the params passed to the renderer.
	params *HomeParams
}

// NewHomeRenderer creates a new instance of CompareRootRenderer.
func NewHomeRenderer(params *HomeParams) *HomeRenderer {
	homeRenderer := &HomeRenderer{
		TabPaneWidget:    widgets.NewTabPane("Tables", "Columns", "Views", "Procedures", "Functions", "Triggers"),
		DiffWidget:       widgets.NewParagraph(),
		SqlWidget:        widgets.NewParagraph(),
		showConfirmation: true,
		params:           params,
	}

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

	switch self.TabPaneWidget.ActiveTabIndex {
	case 0: // Tables
		self.DiffWidget.Title = "Tables"
		self.DiffWidget.Text = ""
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
			res := migration.TableDiff(self.params.FirstDb.SqlConnection, self.params.SecondDb.SqlConnection)
			_ = res

			// self.DiffWidget.Text = res
		}

		if confirmationWidget.IsSome() {
			termui.Render(self.TabPaneWidget, self.DiffWidget, self.SqlWidget, confirmationWidget.Unwrap())
		} else {
			termui.Render(self.TabPaneWidget, self.DiffWidget, self.SqlWidget)
		}
	case 1:
		// Columns
		self.DiffWidget.Title = "Columns"
		self.DiffWidget.Text = ""
		termui.Render(self.TabPaneWidget, self.DiffWidget, self.SqlWidget)
	case 2:
		// Views
		self.DiffWidget.Title = "Views"
		self.DiffWidget.Text = ""
		termui.Render(self.TabPaneWidget, self.DiffWidget, self.SqlWidget)
	case 3:
		// Procedures
		self.DiffWidget.Title = "Procedures"
		self.DiffWidget.Text = ""
		termui.Render(self.TabPaneWidget, self.DiffWidget, self.SqlWidget)
	case 4:
		// Functions
		self.DiffWidget.Title = "Functions"
		self.DiffWidget.Text = ""
		termui.Render(self.TabPaneWidget, self.DiffWidget, self.SqlWidget)
	case 5:
		// Triggers
		self.DiffWidget.Title = "Triggers"
		self.DiffWidget.Text = ""
		termui.Render(self.TabPaneWidget, self.DiffWidget, self.SqlWidget)
	}
}
