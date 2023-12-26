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

const (
	defaultBarMsg = "Press <h> or <?> for help."
)

var focusedWidgetBorderStyle = termui.NewStyle(termui.ColorGreen)

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

	// alertMsg holds the alert message to the user if it is set.
	alertMsg safego.Option[string]

	// tabsData holds the data for each tab.
	tabsData [6]tabData

	// params holds the data parameters that are passed to the PatchiRenderer.
	params *PatchiRendererParams

	// alreadyRenderedEntities Makes sure that we don't generate SQL for an entity twice.
	alreadyRenderedEntities map[string]map[string]bool
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
		alertMsg:                safego.None[string](),
		params:                  params,
		alreadyRenderedEntities: map[string]map[string]bool{},
	}

	// Initialize the map of maps with default values because an empty map is nil in Go for some reason.
	patchiRenderer.alreadyRenderedEntities["tables"] = map[string]bool{}
	patchiRenderer.alreadyRenderedEntities["columns"] = map[string]bool{}
	patchiRenderer.alreadyRenderedEntities["views"] = map[string]bool{}
	patchiRenderer.alreadyRenderedEntities["procedures"] = map[string]bool{}
	patchiRenderer.alreadyRenderedEntities["functions"] = map[string]bool{}
	patchiRenderer.alreadyRenderedEntities["triggers"] = map[string]bool{}

	patchiRenderer.FocusedWidget = patchiRenderer.DiffWidget

	// Setup the initial design/styling of the widgets. Sizing is done in ResizeWidgets() or in Render().
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
	patchiRenderer.MessageBarWidget.Text = defaultBarMsg

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
		`[<a>](fg:green)` + "\t \t \t \t \t \t \t \t \t \t \t \t on any tab to generate all the SQL at once.",
	}

	patchiRenderer.confirmationWidget.BorderTop = false

	patchiRenderer.confirmationWidget.Text = "Press Enter to fetch changes."

	// Set the ShowConfirmation flag to true on each view.
	for i := 0; i < len(patchiRenderer.tabsData); i += 1 {
		patchiRenderer.tabsData[i].ShowConfirmation = true
	}

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

	self.TabPaneWidget.SetRect(0, 0, width/2, height/tabPaneHeight)
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

// generateSqlFor is responsible for generating SQL for anything in the database (tables, columns, etc.)
func (self *PatchiRenderer) generateSqlFor(entityType string, entityRow string) string {
	// This is the status of the entity (created, deleted, modified)
	entityStatus := utils.ExtractExpressions(entityRow, "fg:(.*?)\\)")[0]
	if entityStatus == "green" {
		entityStatus = "created"
	} else if entityStatus == "red" {
		entityStatus = "deleted"
	}

	dialect := self.params.FirstDb.Info.Dialect

	// Extract the name from the row. `[users](fg:green)` -> `users` .
	entityName := utils.ExtractExpressions(entityRow, "\\[(.*?)\\]")[0]

	var generatedSql string
	if entityType == "tables" {

		generatedSql = sequelizer.GenerateSqlForTables(self.params.FirstDb.SqlConnection, dialect, entityName, entityStatus)

	} else if entityType == "columns" {

		// Since the `currentlySelectedEntityName` for columns contain a "tableName → columnName" pattern.
		// We need to extract the table name and column name using regex first.
		extractedExpressions := strings.Fields(entityName)
		tableName := extractedExpressions[0]
		columnName := extractedExpressions[2]

		var errMsg safego.Option[string]
		generatedSql, errMsg = sequelizer.GenerateSqlForColumns(self.params.FirstDb, dialect, columnName, tableName, entityStatus)
		if errMsg.IsSome() {
			self.alert(errMsg.Unwrap())
		}

	} else if entityType == "views" {

		generatedSql = sequelizer.GenerateSqlForViews(self.params.FirstDb.SqlConnection, dialect, entityName, entityStatus)

	} else if entityType == "procedures" {

		generatedSql = sequelizer.GenerateSqlForProcedures(self.params.FirstDb.SqlConnection, dialect, entityName, entityStatus)

	} else if entityType == "functions" {

		generatedSql = sequelizer.GenerateSqlForFunctions(self.params.FirstDb.SqlConnection, dialect, entityName, entityStatus)

	} else if entityType == "triggers" {

		generatedSql = sequelizer.GenerateSqlForTriggers(self.params.FirstDb.SqlConnection, dialect, entityName, entityStatus)

	}

	return generatedSql
}

func (self *PatchiRenderer) GenerateSqlForAllEntities() {
	var generatedSql string

	if len(self.DiffWidget.Rows) == 0 {
		return
	}

	if len(self.DiffWidget.Rows) <= self.DiffWidget.SelectedRow {
		return
	}

	if len(utils.ExtractExpressions(self.DiffWidget.Rows[self.DiffWidget.SelectedRow], "\\[(.*?)\\]")) == 0 {
		return
	}

	currentlySelectedTab := getTabNameBasedOnIndex(self.TabPaneWidget.ActiveTabIndex) // tables, columns, views, ..
	for i, entityRow := range self.DiffWidget.Rows {
		entityNameExtracted := utils.ExtractExpressions(entityRow, "\\[(.*?)\\]")[0]

		if ok, _ := self.alreadyRenderedEntities[currentlySelectedTab][entityNameExtracted]; !ok {
			leftMargin := "\n\n"

			if i == 0 && len(self.SqlWidget.Text) == 0 {
				leftMargin = ""
			}

			generatedSql += leftMargin + self.generateSqlFor(currentlySelectedTab, entityRow)

			self.alreadyRenderedEntities[currentlySelectedTab][entityNameExtracted] = true
		}
	}

	if self.SqlWidget.Text != "" {
		self.SqlWidget.Text += "\n\n" + generatedSql
	} else {
		self.SqlWidget.Text += generatedSql
	}
}

// HandleActionOnEnter Handle every case scenario of pressing the "action button" in any state of the app.
// It may sound obvious but this doesn't render anything. It just changes the state that the Render method relies on.
func (self *PatchiRenderer) HandleActionOnEnter() {
	optErrPrompt := safego.None[string]()

	if self.tabsData[self.TabPaneWidget.ActiveTabIndex].ShowConfirmation { // The user pressed <Enter> in the "Press Enter to start." stage.
		// Setting this to false means that the Render method will not render the confirmation message (but instead
		// render the actual db diff.)
		self.tabsData[self.TabPaneWidget.ActiveTabIndex].ShowConfirmation = false
	} else if self.FocusedWidget == self.DiffWidget { // user pressed <Enter> on an entity off the list on the Diff view.
		if len(self.DiffWidget.Rows) == 0 {
			return
		}

		if len(self.DiffWidget.Rows) <= self.DiffWidget.SelectedRow {
			return
		}

		if len(utils.ExtractExpressions(self.DiffWidget.Rows[self.DiffWidget.SelectedRow], "\\[(.*?)\\]")) == 0 {
			return
		}

		// generate the SQL for this entity (an entity could be a name of table that is created for example) based on its
		// type (table, column, ..etc.)

		currentlySelectedTab := getTabNameBasedOnIndex(self.TabPaneWidget.ActiveTabIndex) // tables, columns, views, ..
		// For columns, this value would be `users -> email`
		currentlySelectedEntityNameExtracted := utils.ExtractExpressions(self.DiffWidget.Rows[self.DiffWidget.SelectedRow], "\\[(.*?)\\]")[0]

		if ok, _ := self.alreadyRenderedEntities[currentlySelectedTab][currentlySelectedEntityNameExtracted]; !ok { // Check if we already generated the SQL for this entity.

			generatedSql := self.generateSqlFor(currentlySelectedTab, self.DiffWidget.Rows[self.DiffWidget.SelectedRow]) // Give the full row.

			if self.SqlWidget.Text != "" {
				self.SqlWidget.Text += "\n\n" + generatedSql
			} else {
				self.SqlWidget.Text += generatedSql
			}

			self.alreadyRenderedEntities[currentlySelectedTab][currentlySelectedEntityNameExtracted] = true
		}
	} else if self.FocusedWidget == self.SqlWidget {
		if self.SqlWidget.Text == "" { // No SQL is generated.
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

// alert opens a pop-up to the user showing a message. Used to report errors that don't panic the app to the user.
func (self *PatchiRenderer) alert(message string) {
	self.alertMsg = safego.Some("[" + message + "](fg:red)")
}

// ResetMsgBar resets the user message in the bottom left to the default one.
func (self *PatchiRenderer) ResetMsgBar() {
	self.alertMsg = safego.Some(defaultBarMsg)
}

// RenderWidgets is the final step in each event. It chooses what to show based on the state current of the app.
func (self *PatchiRenderer) RenderWidgets(userPrompt safego.Option[string]) {
	termui.Clear()

	tabName := getTabNameBasedOnIndex(self.TabPaneWidget.ActiveTabIndex)

	self.DiffWidget.Title = utils.CapitalizeWord(tabName)

	var numberOfChangesForEachTabToBeLoggedToUser int
	// Check this is the first time the user click "Enter" to generate the diff.
	if !self.tabsData[self.TabPaneWidget.ActiveTabIndex].ShowConfirmation && len(self.tabsData[self.TabPaneWidget.ActiveTabIndex].data) == 0 {

		if self.TabPaneWidget.ActiveTabIndex == 0 { // Tables
			// Get the diff data for the current tab that we're on.
			diffResult := difftool.GetTablesDiff(self.params.FirstDb, self.params.SecondDb, self.params.FirstDb.Info.Dialect)

			numberOfChangesForEachTabToBeLoggedToUser = len(diffResult)

			for _, tableDiff := range diffResult {
				text := "[" + tableDiff.TableName + "]"
				if tableDiff.DiffType == 1 {
					text += "(fg:green)"
				} else if tableDiff.DiffType == 0 {
					text += "(fg:red)"
				}

				self.DiffWidget.Rows = append(self.DiffWidget.Rows, text)
				self.tabsData[self.TabPaneWidget.ActiveTabIndex].data = self.DiffWidget.Rows
			}

		} else if self.TabPaneWidget.ActiveTabIndex == 1 { // Columns

			diffResult, errOpt := difftool.GetColumnsDiff(self.params.FirstDb, self.params.SecondDb, self.params.FirstDb.Info.Dialect)
			if errOpt.IsSome() {
				self.alert("An error occurred while fetching the diff for the columns: " + errOpt.Unwrap().Error())
			} else {

				numberOfChangesForEachTabToBeLoggedToUser = len(diffResult)

				for _, columnDiff := range diffResult {
					text := "[" + columnDiff.TableName + " → " + columnDiff.ColumnName + "]"
					if columnDiff.DiffType == 1 {
						text += "(fg:green)"
					} else if columnDiff.DiffType == 0 {
						text += "(fg:red)"
					}

					self.DiffWidget.Rows = append(self.DiffWidget.Rows, text)
					self.tabsData[self.TabPaneWidget.ActiveTabIndex].data = self.DiffWidget.Rows
				}
			}

		} else if self.TabPaneWidget.ActiveTabIndex == 2 { // Views

			diffResult := difftool.GetViewsDiff(self.params.FirstDb, self.params.SecondDb, self.params.FirstDb.Info.Dialect)

			numberOfChangesForEachTabToBeLoggedToUser = len(diffResult)

			for _, viewDiff := range diffResult {
				text := "[" + viewDiff.ViewName + "]"
				if viewDiff.DiffType == 1 {
					text += "(fg:green)"
				} else if viewDiff.DiffType == 0 {
					text += "(fg:red)"
				}

				self.DiffWidget.Rows = append(self.DiffWidget.Rows, text)
				self.tabsData[self.TabPaneWidget.ActiveTabIndex].data = self.DiffWidget.Rows
			}

		} else if self.TabPaneWidget.ActiveTabIndex == 3 { // Procedures

			diffResult := difftool.GetProceduresDiff(self.params.FirstDb, self.params.SecondDb, self.params.FirstDb.Info.Dialect)

			numberOfChangesForEachTabToBeLoggedToUser = len(diffResult)

			for _, viewDiff := range diffResult {
				text := "[" + viewDiff.ProcedureName + "]"
				if viewDiff.DiffType == 1 {
					text += "(fg:green)"
				} else if viewDiff.DiffType == 0 {
					text += "(fg:red)"
				}

				self.DiffWidget.Rows = append(self.DiffWidget.Rows, text)
				self.tabsData[self.TabPaneWidget.ActiveTabIndex].data = self.DiffWidget.Rows
			}

		} else if self.TabPaneWidget.ActiveTabIndex == 4 { // Functions

			diffResult := difftool.GetFunctionsDiff(self.params.FirstDb, self.params.SecondDb, self.params.FirstDb.Info.Dialect)

			numberOfChangesForEachTabToBeLoggedToUser = len(diffResult)

			for _, viewDiff := range diffResult {
				text := "[" + viewDiff.FunctionName + "]"
				if viewDiff.DiffType == 1 {
					text += "(fg:green)"
				} else if viewDiff.DiffType == 0 {
					text += "(fg:red)"
				}

				self.DiffWidget.Rows = append(self.DiffWidget.Rows, text)
				self.tabsData[self.TabPaneWidget.ActiveTabIndex].data = self.DiffWidget.Rows
			}

		} else if self.TabPaneWidget.ActiveTabIndex == 5 { // Triggers
			diffResult := difftool.GetTriggersDiff(self.params.FirstDb, self.params.SecondDb, self.params.FirstDb.Info.Dialect)

			numberOfChangesForEachTabToBeLoggedToUser = len(diffResult)

			for _, triggerDiff := range diffResult {
				text := "[" + triggerDiff.TriggerName + "]"
				if triggerDiff.DiffType == 1 {
					text += "(fg:green)"
				} else if triggerDiff.DiffType == 0 {
					text += "(fg:red)"
				}

				self.DiffWidget.Rows = append(self.DiffWidget.Rows, text)
				self.tabsData[self.TabPaneWidget.ActiveTabIndex].data = self.DiffWidget.Rows
			}
		}

		self.alertMsg = safego.None[string]()
		self.MessageBarWidget.Text = "Found " + strconv.Itoa(numberOfChangesForEachTabToBeLoggedToUser) + " changes in " + strings.ToLower(self.DiffWidget.Title) + "."
	} else {
		// No need to generate the SQL. It has already been generated and saved for this current tab.
		self.DiffWidget.Rows = self.tabsData[self.TabPaneWidget.ActiveTabIndex].data
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

	if self.tabsData[self.TabPaneWidget.ActiveTabIndex].ShowConfirmation {
		diffWidgetRec := self.DiffWidget.GetRect()
		self.confirmationWidget.SetRect(diffWidgetRec.Min.X, diffWidgetRec.Min.Y+1, diffWidgetRec.Max.X, diffWidgetRec.Max.Y)
	} else {
		self.confirmationWidget.SetRect(0, 0, 0, 0)
	}

	if self.alertMsg.IsSome() {
		self.MessageBarWidget.Text = self.alertMsg.UnwrapOr("")
		self.alertMsg = safego.None[string]()
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

type PatchiRendererParams struct {
	FirstDb  types.DbConnection
	SecondDb types.DbConnection
}

type tabData struct {
	ShowConfirmation bool
	data             []string
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
