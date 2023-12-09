package difftool

import (
	"database/sql"
	"fmt"
	"github.com/Okira-E/patchi/pkg/types"
	"github.com/Okira-E/patchi/pkg/utils"
)

// getDiffInTables represents the tables out of sync between two databases.
func getDiffInTables(firstDb *sql.DB, secondDb *sql.DB, dialect string) []types.TableDiff {
	ret := []types.TableDiff{}

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
