package difftool

import (
	"database/sql"
	"fmt"

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
func GetTablesDiff(firstDb *sql.DB, secondDb *sql.DB, dialect string) []tableDiff {
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
func getAllTablesNamesInDb(db *sql.DB, dialect string) []string {
	ret := []string{}

	if dialect == "mysql" || dialect == "mariadb" {
		ret = getAllTablesNamesInMysql(db, dialect)
	} else if dialect == "postgres" || dialect == "cockroachdb" {
		ret = getAllTablesNamesInPostgres(db, dialect)
	}
	return ret
}

func getAllTablesNamesInMysql(db *sql.DB, dialect string) []string {
	ret := []string{}

	rows, err := db.Query("SHOW TABLES")
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

func getAllTablesNamesInPostgres(db *sql.DB, dialect string) []string {
	ret := []string{}

	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
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
