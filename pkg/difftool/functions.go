package difftool

import (
	"fmt"

	"github.com/Okira-E/patchi/pkg/types"
	"github.com/Okira-E/patchi/pkg/utils"
)

type functionDiff struct {
	FunctionName string
	// DiffType represents the type of change that has occurred to the entity.
	// 0 -> Deleted.
	// 1 -> Created.
	DiffType int8
}

// GetFunctionsDiff returns the functions out of sync between two databases.
func GetFunctionsDiff(firstDb types.DbConnection, secondDb types.DbConnection, dialect string) []functionDiff {
	ret := []functionDiff{}

	functionsInFirstDb := getAllFunctionsNamesInDb(firstDb, dialect)
	functionsInSecondDb := getAllFunctionsNamesInDb(secondDb, dialect)

	// Compare the two arrays of function names and return the difference.
	// Functions that exist in the first database but not in the second database must have been created.
	// Functions that exist in the second database but not in the first database must have been deleted.

	firstDbFunctionsBookKeeping := map[string]bool{}
	for _, functionName := range functionsInFirstDb {
		firstDbFunctionsBookKeeping[functionName] = true
	}

	secondDbFunctionsBookKeeping := map[string]bool{}
	for _, functionName := range functionsInSecondDb {
		secondDbFunctionsBookKeeping[functionName] = true
	}

	for _, functionName := range functionsInFirstDb {
		if _, ok := secondDbFunctionsBookKeeping[functionName]; !ok {
			ret = append(ret, functionDiff{
				FunctionName: functionName,
				DiffType:  1,
			})
		}
	}

	for _, functionName := range functionsInSecondDb {
		if _, ok := firstDbFunctionsBookKeeping[functionName]; !ok {
			ret = append(ret, functionDiff{
				FunctionName: functionName,
				DiffType:  0,
			})
		}
	}

	return ret
}

// getAllFunctionsNamesInDb returns an array of function names retrieved from given database connection.
func getAllFunctionsNamesInDb(db types.DbConnection, dialect string) []string {
	ret := []string{}

	if dialect == "mysql" || dialect == "mariadb" {
		ret = getAllFunctionsNamesInMysql(db, dialect)
	} else if dialect == "postgres" || dialect == "cockroachdb" {
		utils.AbortTui("UNIMPLEMENTED")
	}
	return ret
}

func getAllFunctionsNamesInMysql(db types.DbConnection, dialect string) []string {
	ret := []string{}

	rows, err := db.SqlConnection.Query(
		"SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE = 'FUNCTION'",
		db.Info.DatabaseName,
	)
	if err != nil {
		utils.Abort(fmt.Sprintf("Error querying database: %s", err.Error()))
	}

	for rows.Next() {
		var functionName string
		err := rows.Scan(&functionName)
		if err != nil {
			utils.Abort(fmt.Sprintf("Error scanning row: %s", err.Error()))
		}

		ret = append(ret, functionName)
	}

	return ret
}
