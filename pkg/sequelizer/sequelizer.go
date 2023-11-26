package sequelizer

import (
	"github.com/Okira-E/patchi/pkg/sequelizer/mysql"
	"github.com/Okira-E/patchi/pkg/types"
)

// PatchSqlForEntity Generates/retrieves the SQL for creating/deleting/modifying an entity.
//
// - entityType: The type of entity (table, view, etc.)
//
// - entityName: The name of the entity
//
// - status: The status of the entity (created, deleted, modified)
func PatchSqlForEntity(firstDb types.DbConnection, secondDb types.DbConnection, entityType string, entityName string, status string) string {
	var sql string

	if entityType == "tables" {
		sql = patchSqlForTables(firstDb, secondDb, entityType, entityName, status)
	} else if entityType == "columns" {
		sql = patchSqlForColumns(firstDb, secondDb, entityType, entityName, status)
	} else if entityType == "views" {
		sql = patchSqlForViews(firstDb, secondDb, entityType, entityName, status)
	} else if entityType == "procedures" {
		sql = patchSqlForProcedures(firstDb, secondDb, entityType, entityName, status)
	} else if entityType == "functions" {
		sql = patchSqlForFunctions(firstDb, secondDb, entityType, entityName, status)
	} else if entityType == "triggers" {
		sql = patchSqlForTriggers(firstDb, secondDb, entityType, entityName, status)
	}

	return sql
}

func patchSqlForTables(firstDb types.DbConnection, secondDb types.DbConnection, entityType string, entityName string, status string) string {
	var ret string

	dialect := firstDb.Info.Dialect

	if dialect == "mysql" || dialect == "mariadb" {
		ret = mysql.RecreateSqlForTables(firstDb, secondDb, entityType, entityName, status)
	} else if dialect == "postgres" || dialect == "cockroachdb" {
		// UNIMPLEMENTED
	}

	return ret
}

func patchSqlForColumns(firstDb types.DbConnection, secondDb types.DbConnection, entityType string, entityName string, status string) string {
	var ret string

	dialect := firstDb.Info.Dialect

	if dialect == "mysql" || dialect == "mariadb" {
		// UNIMPLEMENTED
	} else if dialect == "postgres" || dialect == "cockroachdb" {
		// UNIMPLEMENTED
	}

	return ret
}

func patchSqlForViews(firstDb types.DbConnection, secondDb types.DbConnection, entityType string, entityName string, status string) string {
	var ret string

	dialect := firstDb.Info.Dialect

	if dialect == "mysql" || dialect == "mariadb" {
		// UNIMPLEMENTED
	} else if dialect == "postgres" || dialect == "cockroachdb" {
		// UNIMPLEMENTED
	}

	return ret
}

func patchSqlForProcedures(firstDb types.DbConnection, secondDb types.DbConnection, entityType string, entityName string, status string) string {
	var ret string

	dialect := firstDb.Info.Dialect

	if dialect == "mysql" || dialect == "mariadb" {
		// UNIMPLEMENTED
	} else if dialect == "postgres" || dialect == "cockroachdb" {
		// UNIMPLEMENTED
	}

	return ret
}

func patchSqlForFunctions(firstDb types.DbConnection, secondDb types.DbConnection, entityType string, entityName string, status string) string {
	var ret string

	dialect := firstDb.Info.Dialect

	if dialect == "mysql" || dialect == "mariadb" {
		// UNIMPLEMENTED
	} else if dialect == "postgres" || dialect == "cockroachdb" {
		// UNIMPLEMENTED
	}

	return ret
}

func patchSqlForTriggers(firstDb types.DbConnection, secondDb types.DbConnection, entityType string, entityName string, status string) string {
	var ret string

	dialect := firstDb.Info.Dialect

	if dialect == "mysql" || dialect == "mariadb" {
		// UNIMPLEMENTED
	} else if dialect == "postgres" || dialect == "cockroachdb" {
		// UNIMPLEMENTED
	}

	return ret
}
