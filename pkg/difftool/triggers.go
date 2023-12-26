package difftool

import (
	"fmt"

	"github.com/Okira-E/patchi/pkg/types"
	"github.com/Okira-E/patchi/pkg/utils"
)

type triggerDiff struct {
	TriggerName string
	// DiffType represents the type of change that has occurred to the entity.
	// 0 -> Deleted.
	// 1 -> Created.
	DiffType int8
}

// GetTriggersDiff returns the triggers out of sync between two databases.
func GetTriggersDiff(firstDb types.DbConnection, secondDb types.DbConnection, dialect string) []triggerDiff {
	ret := []triggerDiff{}

	triggersInFirstDb := getAllTriggersNamesInDb(firstDb, dialect)
	triggersInSecondDb := getAllTriggersNamesInDb(secondDb, dialect)

	// Compare the two arrays of trigger names and return the difference.
	// Triggers that exist in the first database but not in the second database must have been created.
	// Triggers that exist in the second database but not in the first database must have been deleted.

	firstDbTriggersBookKeeping := map[string]bool{}
	for _, triggerName := range triggersInFirstDb {
		firstDbTriggersBookKeeping[triggerName] = true
	}

	secondDbTriggersBookKeeping := map[string]bool{}
	for _, triggerName := range triggersInSecondDb {
		secondDbTriggersBookKeeping[triggerName] = true
	}

	for _, triggerName := range triggersInFirstDb {
		if _, ok := secondDbTriggersBookKeeping[triggerName]; !ok {
			ret = append(ret, triggerDiff{
				TriggerName: triggerName,
				DiffType:    1,
			})
		}
	}

	for _, triggerName := range triggersInSecondDb {
		if _, ok := firstDbTriggersBookKeeping[triggerName]; !ok {
			ret = append(ret, triggerDiff{
				TriggerName: triggerName,
				DiffType:    0,
			})
		}
	}

	return ret
}

// getAllTriggersNamesInDb returns an array of trigger names retrieved from given database connection.
func getAllTriggersNamesInDb(db types.DbConnection, dialect string) []string {
	ret := []string{}

	if dialect == "mysql" || dialect == "mariadb" {
		ret = getAllTriggersNamesInMysql(db, dialect)
	} else if dialect == "postgres" || dialect == "cockroachdb" {
		utils.AbortTui("UNIMPLEMENTED")
	}
	return ret
}

func getAllTriggersNamesInMysql(db types.DbConnection, dialect string) []string {
	ret := []string{}

	rows, err := db.SqlConnection.Query(
		"SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ?",
		db.Info.DatabaseName,
	)
	if err != nil {
		utils.Abort(fmt.Sprintf("Error querying database: %s", err.Error()))
	}

	for rows.Next() {
		var triggerName string
		err := rows.Scan(&triggerName)
		if err != nil {
			utils.Abort(fmt.Sprintf("Error scanning row: %s", err.Error()))
		}

		ret = append(ret, triggerName)
	}

	return ret
}
