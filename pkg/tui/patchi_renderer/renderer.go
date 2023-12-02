package patchi_renderer

import (
	"strconv"
	"strings"

	"github.com/Okira-E/patchi/pkg/difftool"
	"github.com/Okira-E/patchi/pkg/sequelizer"
	"github.com/Okira-E/patchi/pkg/utils"
	"github.com/atotto/clipboard"

	"github.com/Okira-E/patchi/pkg/types"
	"github.com/Okira-E/patchi/safego"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var focusedWidgetBorderStyle = termui.NewStyle(termui.ColorGreen)

type PatchiRendererParams struct {
	FirstDb  types.DbConnection
	SecondDb types.DbConnection
}

// PatchiRenderer is a unit that knows about all the widgets that need to be rendered.
// rendering the TUI should always be done by calling the `render` method.
type PatchiRenderer struct {
	// width and height are the dimensions state of the terminal.
	width, height int

	// TabPaneWidget is the widget that holds the tabs for the different types of objects.
	TabPaneWidget *widgets.TabPane

	// DiffWidget is the widget that holds the diff of the selected tab.
	DiffWidget *widgets.List

	// SqlWidget is the widget that holds the sql formulated from the diff.
	SqlWidget *widgets.Paragraph

	// MessageBarWidget is the widget that holds the messages that are shown to the user.
	MessageBarWidget *widgets.Paragraph

	// focusedWidget is the widget that is currently being controlled by the user.
	FocusedWidget any // FIX: Could not find an interface that is a union of all the widgets in termui so I used `any`.

	// LastFocusedWidget is the widget that was focused before the current one.
	LastFocusedWidget any

	// HelpWidget shows all the keyboard shortcuts for the app.
	HelpWidget *widgets.List

	// confirmationWidget is the widget that is shown when the user has not yet started the comparing process.
	confirmationWidget *widgets.Paragraph

	// ShowHelpWidget determines if Render method needs to render HelpWidget or not.
	ShowHelpWidget bool

	// ShowConfirmation determines if Render method needs to render confirmationWidget or not.
	ShowConfirmation bool

	// tabsData holds the data for each tab.
	tabsData [6][]string

	// params holds the data parameters that are passed to the PatchiRenderer.
	params *PatchiRendererParams

	// alreadyRenderedEntities Makes sure that we don't generate SQL for an entity twice.
	alreadyRenderedEntities map[string]bool
}

// NewPatchiRenderer creates a new instance of CompareRootRenderer.
func NewPatchiRenderer(params *PatchiRendererParams) *PatchiRenderer {
	patchiRenderer := &PatchiRenderer{
		TabPaneWidget:           widgets.NewTabPane("Tables", "Columns", "Views", "Procedures", "Functions", "Triggers"),
		DiffWidget:              widgets.NewList(),
		SqlWidget:               widgets.NewParagraph(),
		MessageBarWidget:        widgets.NewParagraph(),
		HelpWidget:              widgets.NewList(),
		confirmationWidget:      widgets.NewParagraph(),
		ShowConfirmation:        true,
		params:                  params,
		alreadyRenderedEntities: map[string]bool{},
	}

	patchiRenderer.FocusedWidget = patchiRenderer.DiffWidget

	patchiRenderer.DiffWidget.TextStyle = termui.NewStyle(termui.ColorWhite)
	patchiRenderer.DiffWidget.SelectedRowStyle = termui.NewStyle(termui.ColorBlack, termui.ColorWhite)
	patchiRenderer.DiffWidget.WrapText = true
	patchiRenderer.DiffWidget.PaddingLeft = 2
	patchiRenderer.DiffWidget.Rows = []string{}

	patchiRenderer.SqlWidget.Title = "SQL"
	patchiRenderer.SqlWidget.TextStyle = termui.NewStyle(termui.ColorWhite)
	patchiRenderer.SqlWidget.WrapText = true
	patchiRenderer.SqlWidget.PaddingLeft = 2
	patchiRenderer.SqlWidget.Text = ""

	patchiRenderer.MessageBarWidget.Border = false
	patchiRenderer.MessageBarWidget.Text = `Press <h> or <?> for help.`

	patchiRenderer.HelpWidget.Title = "Help"
	patchiRenderer.HelpWidget.SelectedRowStyle = termui.NewStyle(termui.ColorBlack, termui.ColorWhite)
	patchiRenderer.HelpWidget.Rows = []string{
		`[<h | ?>](fg:green)` + "\t \t \t \t \t \t \t \t to show this help widget.",
		`[<Escape>](fg:green)` + "\t \t \t \t \t \t \t to exit help.",
		`[<q | Ctrl+c>](fg:green)` + "\t \t \t to quit.",
		// NOTE: The two rules below are not highlighted like the rest because a never closing '[' breaks termui regex.
		`<[ | Left>` + "\t \t \t \t \t to move to the previous tab.",
		`<] | Right>` + "\t \t \t \t to move to the next tab.",
		`[<Tab>](fg:green)` + "\t \t \t \t \t \t \t \t \t \t to move between the diff and sql widgets.",
		`[<Enter>](fg:green)` + "\t \t \t \t \t \t \t \t on the SQL widget to copy the SQL.",
	}

	patchiRenderer.confirmationWidget.BorderTop = false

	patchiRenderer.confirmationWidget.Text = "Press Enter to fetch changes."

	return patchiRenderer
}

// ResizeWidgets resizes the widgets based on the new width and height of the terminal.
// This method should be called whenever the terminal is resized.
// NOTE: It is also what sets the original sizing for most of the widgets. I say most because widgets with conditional
// rendering like the help widget only get their sizing set when they are rendered. They can all be found at the end
// of the Render method. The Render method uses the width and height set by this widget whenever the window is
// resized, so conditional widgets that are sized in the Render method will also be corrected when the window is resized.
func (self *PatchiRenderer) ResizeWidgets(width int, height int) {
	const tabPaneHeight = 14

	messageBarHeight := height / 35

	self.TabPaneWidget.SetRect(0, 0, width/2, height/tabPaneHeight) // BUG: Put a limit on the height of the tab pane.
	self.DiffWidget.SetRect(0, height/tabPaneHeight, width/2, height-(messageBarHeight))
	self.SqlWidget.SetRect(width/2, 0, width, height-(messageBarHeight))
	self.MessageBarWidget.SetRect(0, height-(messageBarHeight), width, height)

	self.width = width
	self.height = height
}

// ClearBorderStyles clears the border styles of all the widgets.
func (self *PatchiRenderer) ClearBorderStyles() {
	self.DiffWidget.BorderStyle = termui.NewStyle(termui.ColorClear)
	self.SqlWidget.BorderStyle = termui.NewStyle(termui.ColorClear)
}

// ToggleHelpWidget toggles the help widget.
func (self *PatchiRenderer) ToggleHelpWidget() {
	if self.ShowHelpWidget {
		self.ShowHelpWidget = false
		self.FocusedWidget = self.LastFocusedWidget

		self.MessageBarWidget.Text = `Press <h> or <?> for help.`
	} else {
		self.ShowHelpWidget = true
		self.LastFocusedWidget = self.FocusedWidget
		self.FocusedWidget = self.HelpWidget

		self.MessageBarWidget.Text = `Press <Escape> to exit help.`
	}
}

// HandleActionOnEnter Handle every case scenario of pressing the "action button" in any state of the app.
// It may sound obvious but this doesn't render anything. It just changes the state that the Render method relies on.
func (self *PatchiRenderer) HandleActionOnEnter() {
	optErrPrompt := safego.None[string]()

	if self.ShowConfirmation { // The user pressed <Enter> in the "Press Enter to start." stage.
		// Setting this to false means that the Render method will not render the confirmation message (but instead
		// render the actual db diff.)
		self.ShowConfirmation = false
	} else if self.FocusedWidget == self.DiffWidget { // user pressed <Enter> on an entity off the list on the Diff view.
		// May never happen but being paranoid is better than debugging a runtime error.
		if len(self.DiffWidget.Rows) == 0 || len(utils.ExtractExpressions(self.DiffWidget.Rows[self.DiffWidget.SelectedRow], "\\[(.*?)\\]")) == 0 {
			return
		}

		// generate the SQL for this entity (an entity could be a name of table that is created for example) based on its 
		// type (table, column, ..etc).

		currentlySelectedTab := getTabNameBasedOnIndex(self.TabPaneWidget.ActiveTabIndex)
		currentlySelectedEntityName := utils.ExtractExpressions(self.DiffWidget.Rows[self.DiffWidget.SelectedRow], "\\[(.*?)\\]")[0]

		if ok, _ := self.alreadyRenderedEntities[currentlySelectedEntityName]; !ok { // Check if we already generated the SQL for this entity.

			// This is the status of the entity (created, deleted, modified)
			entityStatus := utils.ExtractExpressions(self.DiffWidget.Rows[self.DiffWidget.SelectedRow], "fg:(.*?)\\)")[0]
			if entityStatus == "green" {
				entityStatus = "created"
			} else if entityStatus == "red" {
				entityStatus = "deleted"
			} else {
				entityStatus = "modified"
			}

			dialect := self.params.FirstDb.Info.Dialect

			generatedSql := sequelizer.GenerateSqlForEntity(self.params.FirstDb.SqlConnection, self.params.SecondDb.SqlConnection, dialect, currentlySelectedTab, currentlySelectedEntityName, entityStatus)

			self.SqlWidget.Text += generatedSql + "\n\n"

			self.alreadyRenderedEntities[currentlySelectedEntityName] = true
		}
	} else if self.FocusedWidget == self.SqlWidget {
		if self.SqlWidget.Text == "" { // List of entities is empty.
			return
		}

		sql := self.SqlWidget.Text

		err := clipboard.WriteAll(sql)
		if err != nil {
			optErrPrompt = safego.Some("[Error copying to clipboard: " + err.Error() + "](fg:red)")
		} else {
			optErrPrompt = safego.Some("[Copied entire generated SQL to clipboard.](fg:green)")
		}
	}

	self.RenderWidgets(optErrPrompt)

}

// RenderWidgets is the final step in each event. It choses what to show based on the state current of the app.
func (self *PatchiRenderer) RenderWidgets(userPrompt safego.Option[string]) {
	termui.Clear()

	if self.TabPaneWidget.ActiveTabIndex == 0 { // Tables

		self.DiffWidget.Title = "Tables"

		if !self.ShowConfirmation && len(self.tabsData[0]) == 0 {
			diffResult := difftool.GetDiff(self.params.FirstDb.SqlConnection, self.params.SecondDb.SqlConnection, "tables", self.params.FirstDb.Info.Dialect)

			self.MessageBarWidget.Text = "Found " + strconv.Itoa(len(diffResult)) + " changes in " + strings.ToLower(self.DiffWidget.Title) + "."

			for _, tableDiff := range diffResult {
				text := "[" + tableDiff.TableName + "]"
				if tableDiff.DiffType == "created" {
					text += "(fg:green)"
				} else if tableDiff.DiffType == "deleted" {
					text += "(fg:red)"
				}

				self.DiffWidget.Rows = append(self.DiffWidget.Rows, text)
				self.tabsData[0] = self.DiffWidget.Rows
			}
		} else {
			// No need to generate the SQL. It has already been generated and saved for this current tab.
			self.DiffWidget.Rows = self.tabsData[0]
		}

	} else if self.TabPaneWidget.ActiveTabIndex == 1 { // Columns

		self.DiffWidget.Title = "Columns"
		self.DiffWidget.Rows = self.tabsData[1]

	} else if self.TabPaneWidget.ActiveTabIndex == 2 { // Views

		self.DiffWidget.Title = "Views"
		self.DiffWidget.Rows = self.tabsData[2]

	} else if self.TabPaneWidget.ActiveTabIndex == 3 { // Procedures

		self.DiffWidget.Title = "Procedures"
		self.DiffWidget.Rows = self.tabsData[3]

	} else if self.TabPaneWidget.ActiveTabIndex == 4 { // Functions

		self.DiffWidget.Title = "Functions"
		self.DiffWidget.Rows = self.tabsData[4]

	} else if self.TabPaneWidget.ActiveTabIndex == 5 { // Triggers

		self.DiffWidget.Title = "Triggers"
		self.DiffWidget.Rows = self.tabsData[5]

	}

	// Set the currently focused widget's border style to be green.
	self.ClearBorderStyles()

	// Set the BorderStyle for the currently FocusedWidget. Since it is of `any` type it needs to be casted to any
	// type that it could be. Simple `FocusedWidget.BorderStyle` returns an error.
	switch self.FocusedWidget.(type) {
	case *widgets.Paragraph:
		self.FocusedWidget.(*widgets.Paragraph).BorderStyle = focusedWidgetBorderStyle
	case *widgets.List:
		self.FocusedWidget.(*widgets.List).BorderStyle = focusedWidgetBorderStyle
	case *widgets.TabPane:
		self.FocusedWidget.(*widgets.TabPane).BorderStyle = focusedWidgetBorderStyle
	}

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

func getTabNameBasedOnIndex(index int) string {
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
