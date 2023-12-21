package difftool

import (
	"fmt"

	"github.com/Okira-E/patchi/pkg/types"
	"github.com/Okira-E/patchi/pkg/utils"
)

type viewDiff struct {
	ViewName string
	// DiffType represents the type of change that has occurred to the entity.
	// 0 -> Deleted.
	// 1 -> Created.
	DiffType int8
}

// GetViewsDiff returns the views out of sync between two databases.
func GetViewsDiff(firstDb types.DbConnection, secondDb types.DbConnection, dialect string) []viewDiff {
	ret := []viewDiff{}

	viewsInFirstDb := getAllViewsNamesInDb(firstDb, dialect)
	viewsInSecondDb := getAllViewsNamesInDb(secondDb, dialect)

	// Compare the two arrays of view names and return the difference.
	// Views that exist in the first database but not in the second database must have been created.
	// Views that exist in the second database but not in the first database must have been deleted.

	firstDbViewsBookKeeping := map[string]bool{}
	for _, viewName := range viewsInFirstDb {
		firstDbViewsBookKeeping[viewName] = true
	}

	secondDbViewsBookKeeping := map[string]bool{}
	for _, viewName := range viewsInSecondDb {
		secondDbViewsBookKeeping[viewName] = true
	}

	for _, viewName := range viewsInFirstDb {
		if _, ok := secondDbViewsBookKeeping[viewName]; !ok {
			ret = append(ret, viewDiff{
				ViewName: viewName,
				DiffType:  1,
			})
		}
	}

	for _, viewName := range viewsInSecondDb {
		if _, ok := firstDbViewsBookKeeping[viewName]; !ok {
			ret = append(ret, viewDiff{
				ViewName: viewName,
				DiffType:  0,
			})
		}
	}

	return ret
}

// getAllViewsNamesInDb returns an array of view names retrieved from given database connection.
func getAllViewsNamesInDb(db types.DbConnection, dialect string) []string {
	ret := []string{}

	if dialect == "mysql" || dialect == "mariadb" {
		ret = getAllViewsNamesInMysql(db, dialect)
	} else if dialect == "postgres" || dialect == "cockroachdb" {
		utils.AbortTui("UNIMPLEMENTED")
	}
	return ret
}

func getAllViewsNamesInMysql(db types.DbConnection, dialect string) []string {
	ret := []string{}

	rows, err := db.SqlConnection.Query(
		"SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_TYPE = 'VIEW'",
		db.Info.DatabaseName,
	)
	if err != nil {
		utils.Abort(fmt.Sprintf("Error querying database: %s", err.Error()))
	}

	for rows.Next() {
		var viewName string
		err := rows.Scan(&viewName)
		if err != nil {
			utils.Abort(fmt.Sprintf("Error scanning row: %s", err.Error()))
		}

		ret = append(ret, viewName)
	}

	return ret
}
