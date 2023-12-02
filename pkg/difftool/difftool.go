package difftool

import (
	"database/sql"

	"github.com/Okira-E/patchi/pkg/types"
	"github.com/Okira-E/patchi/pkg/utils"
)

func GetDiff(firstDb *sql.DB, secondDb *sql.DB, entityType string, dialect string) []types.TableDiff {
	var diffResult []types.TableDiff

	if entityType == "tables" {
		diffResult = getDiffInTables(firstDb, secondDb, dialect)
	} else if entityType == "columns" {
		utils.AbortTui("UNIMPLEMENTED")
	} else if entityType == "views" {
		utils.AbortTui("UNIMPLEMENTED")
	} else if entityType == "procedures" {
		utils.AbortTui("UNIMPLEMENTED")
	} else if entityType == "functions" {
		utils.AbortTui("UNIMPLEMENTED")
	} else if entityType == "triggers" {
		utils.AbortTui("UNIMPLEMENTED")
	}

	return diffResult
}
