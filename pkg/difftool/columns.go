package difftool

import (
	"encoding/json"
	"slices"

	"github.com/Okira-E/patchi/pkg/types"
	"github.com/Okira-E/patchi/pkg/utils"
	"github.com/Okira-E/patchi/safego"
)

type columnDiff struct {
	ColumnName string
	// The table the column belongs to.
	TableName  string
	// DiffType represents the type of change that has occurred to the entity.
	// 0 -> Deleted.
	// 1 -> Created.
	DiffType int8
}

// GetColumnsDiff returns the columns out of sync between two databases.
func GetColumnsDiff(firstDb types.DbConnection, secondDb types.DbConnection, dialect string) ([]columnDiff, safego.Option[error]) {
	ret := []columnDiff{}
	errOpt := safego.None[error]()

	// There has to be a more optimized version of this query that returns a better data-structer.
	sql := `
		SELECT JSON_OBJECTAGG(
		               TABLE_NAME,
		               (SELECT JSON_ARRAYAGG(JSON_OBJECT('COLUMN_NAME', COLUMN_NAME, 'ORDINAL_POSITION', ORDINAL_POSITION))
		                FROM information_schema.COLUMNS t2
		                WHERE t2.TABLE_SCHEMA = ?
		                  AND t2.TABLE_NAME = t1.TABLE_NAME)
		       )
		FROM information_schema.TABLES t1
		WHERE TABLE_SCHEMA = ?
		  AND TABLE_TYPE = 'BASE TABLE'
	`

	// Running the query for the first db.
	rows, err := firstDb.SqlConnection.Query(sql, firstDb.Info.DatabaseName, firstDb.Info.DatabaseName)
	if err != nil {
		return ret, safego.Some(err)
	}

	firstDbTablesAndColumnsData := map[string][]columnInfo{}
	for rows.Next() {
		var buffer string
		if err := rows.Scan(&buffer); err != nil {
			utils.AbortTui(err.Error())
		}

		err := json.Unmarshal([]byte(buffer), &firstDbTablesAndColumnsData)
		if err != nil {
			utils.AbortTui(err.Error())
		}
	}

	// Running the query for the second db.
	rows, err = secondDb.SqlConnection.Query(sql, secondDb.Info.DatabaseName, secondDb.Info.DatabaseName)
	if err != nil {
		return ret, safego.Some(err)
	}

	secondDbTablesAndColumnsData := map[string][]columnInfo{}
	for rows.Next() {
		var buffer string
		if err := rows.Scan(&buffer); err != nil {
			utils.AbortTui(err.Error())
		}

		err := json.Unmarshal([]byte(buffer), &secondDbTablesAndColumnsData)
		if err != nil {
			utils.AbortTui(err.Error())
		}
	}

	// Loop through the tables in the firstDb and create the diff for the columns that are not sync
	// between the two databases. Obviously, we want to only check tables that exist in both of the
	// environments.
	for tableName, firstDbColumnsForCurrentTable := range firstDbTablesAndColumnsData {
		// Operate only on tables that are in both databases.
		if _, ok := secondDbTablesAndColumnsData[tableName]; ok {
			secondDbColumnsForCurrentTable := secondDbTablesAndColumnsData[tableName]

			// Populate a slice containing just the column names in for the first env.
			columnsInFirstDbTable := []string{}
			for _, columnInfo := range firstDbColumnsForCurrentTable {
				columnsInFirstDbTable = append(columnsInFirstDbTable, columnInfo.ColumnName)
			}

			// Populate a slice containing just the column names in for the second env.
			columnsInSecondDbTable := []string{}
			for _, columnInfo := range secondDbColumnsForCurrentTable {
				columnsInSecondDbTable = append(columnsInSecondDbTable, columnInfo.ColumnName)
			}

			// Loop through the columns in the first env. Columns that do exist in the second env 
			// are discarded, while the rest are added to ret as 'created.'
			for _, columnName := range columnsInFirstDbTable {
				if !slices.Contains(columnsInSecondDbTable, columnName) {
					ret = append(ret, columnDiff{TableName: tableName, ColumnName: columnName, DiffType: 1})
				}
			}

			// Loop through the columns in the second env. Columns that do not exist in the first 
			// are added to ret as 'deleted.'
			for _, columnName := range columnsInSecondDbTable {
				if !slices.Contains(columnsInFirstDbTable, columnName) {
					ret = append(ret, columnDiff{TableName: tableName, ColumnName: columnName, DiffType: 0})
				}
			}
		}
	}

	return ret, errOpt
}

type columnInfo struct {
	ColumnName      string `json:"COLUMN_NAME"`
	OrdinalPosition int    `json:"ORDINAL_POSITION"`
}
