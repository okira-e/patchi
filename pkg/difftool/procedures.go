package difftool

import (
	"fmt"

	"github.com/Okira-E/patchi/pkg/types"
	"github.com/Okira-E/patchi/pkg/utils"
)

type procedureDiff struct {
	ProcedureName string
	// DiffType represents the type of change that has occurred to the entity.
	// 0 -> Deleted.
	// 1 -> Created.
	DiffType int8
}

// GetProceduresDiff returns the procedures out of sync between two databases.
func GetProceduresDiff(firstDb types.DbConnection, secondDb types.DbConnection, dialect string) []procedureDiff {
	ret := []procedureDiff{}

	proceduresInFirstDb := getAllProceduresNamesInDb(firstDb, dialect)
	proceduresInSecondDb := getAllProceduresNamesInDb(secondDb, dialect)

	// Compare the two arrays of procedure names and return the difference.
	// Procedures that exist in the first database but not in the second database must have been created.
	// Procedures that exist in the second database but not in the first database must have been deleted.

	firstDbProceduresBookKeeping := map[string]bool{}
	for _, procedureName := range proceduresInFirstDb {
		firstDbProceduresBookKeeping[procedureName] = true
	}

	secondDbProceduresBookKeeping := map[string]bool{}
	for _, procedureName := range proceduresInSecondDb {
		secondDbProceduresBookKeeping[procedureName] = true
	}

	for _, procedureName := range proceduresInFirstDb {
		if _, ok := secondDbProceduresBookKeeping[procedureName]; !ok {
			ret = append(ret, procedureDiff{
				ProcedureName: procedureName,
				DiffType:  1,
			})
		}
	}

	for _, procedureName := range proceduresInSecondDb {
		if _, ok := firstDbProceduresBookKeeping[procedureName]; !ok {
			ret = append(ret, procedureDiff{
				ProcedureName: procedureName,
				DiffType:  0,
			})
		}
	}

	return ret
}

// getAllProceduresNamesInDb returns an array of procedure names retrieved from given database connection.
func getAllProceduresNamesInDb(db types.DbConnection, dialect string) []string {
	ret := []string{}

	if dialect == "mysql" || dialect == "mariadb" {
		ret = getAllProceduresNamesInMysql(db, dialect)
	} else if dialect == "postgres" || dialect == "cockroachdb" {
		utils.AbortTui("UNIMPLEMENTED")
	}
	return ret
}

func getAllProceduresNamesInMysql(db types.DbConnection, dialect string) []string {
	ret := []string{}

	rows, err := db.SqlConnection.Query(
		"SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE = 'PROCEDURE'",
		db.Info.DatabaseName,
	)
	if err != nil {
		utils.Abort(fmt.Sprintf("Error querying database: %s", err.Error()))
	}

	for rows.Next() {
		var procedureName string
		err := rows.Scan(&procedureName)
		if err != nil {
			utils.Abort(fmt.Sprintf("Error scanning row: %s", err.Error()))
		}

		ret = append(ret, procedureName)
	}

	return ret
}
