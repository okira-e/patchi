package difftool

import (
	"database/sql"
	"fmt"
	"github.com/Okira-E/patchi/pkg/types"
	"github.com/Okira-E/patchi/pkg/utils"
)

// GetDiffInTablesBetweenSchemas represents the tables out of sync between two databases.
func GetDiffInTablesBetweenSchemas(firstDb *types.DbConnection, secondDb *types.DbConnection) []types.TableDiff {
	ret := []types.TableDiff{}

	tablesInFirstDb := getAllTablesNamesInDb(firstDb.SqlConnection)
	tablesInSecondDb := getAllTablesNamesInDb(secondDb.SqlConnection)

	// Compare the two arrays of table names and return the difference.
	// Tables that exist in the first database but not in the second database must have been created.
	// Tables that exist in the second database but not in the first database must have been deleted.

	firstDbTablesBookKeeping := map[string]bool{}
	for _, tableName := range tablesInFirstDb {
		firstDbTablesBookKeeping[tableName] = true
	}

	secondDbTablesBookKeeping := map[string]bool{}
	for _, tableName := range tablesInSecondDb {
		secondDbTablesBookKeeping[tableName] = true
	}

	for _, tableName := range tablesInFirstDb {
		if _, ok := secondDbTablesBookKeeping[tableName]; !ok {
			ret = append(ret, types.TableDiff{
				TableName: tableName,
				DiffType:  "created",
			})
		}
	}

	for _, tableName := range tablesInSecondDb {
		if _, ok := firstDbTablesBookKeeping[tableName]; !ok {
			ret = append(ret, types.TableDiff{
				TableName: tableName,
				DiffType:  "deleted",
			})
		}
	}

	return ret
}

// getAllTablesNamesInDb returns an array of table names retrieved from given database connection.
func getAllTablesNamesInDb(db *sql.DB) []string {
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		utils.Abort(fmt.Sprintf("Error querying database: %s", err.Error()))
	}

	ret := []string{}
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			utils.Abort(fmt.Sprintf("Error scanning row: %s", err.Error()))
		}

		ret = append(ret, table)
	}

	return ret
}
