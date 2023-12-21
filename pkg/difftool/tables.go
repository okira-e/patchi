package difftool

import (
	"fmt"

	"github.com/Okira-E/patchi/pkg/types"
	"github.com/Okira-E/patchi/pkg/utils"
)

type tableDiff struct {
	TableName string
	// DiffType represents the type of change that has occurred to the entity.
	// 0 -> Deleted.
	// 1 -> Created.
	DiffType int8
}

// GetTablesDiff returns the tables out of sync between two databases.
func GetTablesDiff(firstDb types.DbConnection, secondDb types.DbConnection, dialect string) []tableDiff {
	ret := []tableDiff{}

	tablesInFirstDb := getAllTablesNamesInDb(firstDb, dialect)
	tablesInSecondDb := getAllTablesNamesInDb(secondDb, dialect)

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
			ret = append(ret, tableDiff{
				TableName: tableName,
				DiffType:  1,
			})
		}
	}

	for _, tableName := range tablesInSecondDb {
		if _, ok := firstDbTablesBookKeeping[tableName]; !ok {
			ret = append(ret, tableDiff{
				TableName: tableName,
				DiffType:  0,
			})
		}
	}

	return ret
}

// getAllTablesNamesInDb returns an array of table names retrieved from given database connection.
func getAllTablesNamesInDb(db types.DbConnection, dialect string) []string {
	ret := []string{}

	if dialect == "mysql" || dialect == "mariadb" {
		ret = getAllTablesNamesInMysql(db, dialect)
	} else if dialect == "postgres" || dialect == "cockroachdb" {
		utils.AbortTui("UNIMPLEMENTED")
	}
	return ret
}

func getAllTablesNamesInMysql(db types.DbConnection, dialect string) []string {
	ret := []string{}

	rows, err := db.SqlConnection.Query(
		"SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_TYPE = 'BASE TABLE'",
		db.Info.DatabaseName,
	)
	if err != nil {
		utils.Abort(fmt.Sprintf("Error querying database: %s", err.Error()))
	}

	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			utils.Abort(fmt.Sprintf("Error scanning row: %s", err.Error()))
		}

		ret = append(ret, tableName)
	}

	return ret
}
