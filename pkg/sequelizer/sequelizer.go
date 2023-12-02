package sequelizer

import (
	"database/sql"

	"github.com/Okira-E/patchi/pkg/utils"
)

// GenerateSqlForEntity Generates/retrieves the SQL for creating/deleting/modifying an entity.
//
// # Params
//
// - entityType: The type of entity (table, view, etc.)
//
// - entityName: The name of the entity
//
// - status: The status of the entity (created, deleted, modified)
func GenerateSqlForEntity(firstDb *sql.DB, secondDb *sql.DB, dialect string, entityType string, entityName string, status string) string {
	var sql string

	if entityType == "tables" {
		sql = generateSqlForTables(firstDb, secondDb, dialect, entityType, entityName, status)
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

	return sql
}
