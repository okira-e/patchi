package migration

import (
	"database/sql"
	"log"

	"github.com/Okira-E/patchi/pkg/types"
)

// TableDiff represents the tables out of sync between two databases.
func TableDiff(firstDb *types.DbConnection, secondDb *types.DbConnection) []types.TableDiff {
	ret := []types.TableDiff{}

	// Fetch names of every table in the first database
	tablesInFirstDb := getTablesNames(firstDb.SqlConnection)
	tablesInSecondDb := getTablesNames(secondDb.SqlConnection)

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

// getTablesNames returns an array of table names retrieved from given database connection.
func getTablesNames(db *sql.DB) []string {
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		log.Fatalf("Error querying database: %s", err.Error())
	}

	ret := []string{}
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			log.Fatalf("Error scanning row: %s", err.Error())
		}

		ret = append(ret, table)
	}

	return ret
}
