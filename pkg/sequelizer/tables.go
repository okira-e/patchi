package sequelizer

import (
	"database/sql"

	"github.com/Okira-E/patchi/pkg/utils"
)

// generateSqlForTables is the interface for generating SQL for tables in general.
func generateSqlForTables(firstDb *sql.DB, secondDb *sql.DB, dialect string, entityType string, entityName string, status string) string {
	var ret string

	if dialect == "mysql" || dialect == "mariadb" {
		ret = generateSqlForTablesMysql(firstDb, secondDb, entityType, entityName, status)
	} else if dialect == "postgres" || dialect == "cockroachdb" {
		ret = generateSqlForTablesPostgres(firstDb, secondDb, entityType, entityName, status)
	}

	return ret
}

// generateSqlForTablesMysql is responsible for generating SQL for tables in Mysql.
func generateSqlForTablesMysql(firstDb *sql.DB, secondDb *sql.DB, entityType string, entityName string, status string) string {
	var ret string

	if status == "created" {
		rows, err := firstDb.Query("SHOW CREATE TABLE " + entityName)
		if err != nil {
			utils.AbortTui("Error getting create table statement: " + err.Error())
		}

		for rows.Next() {
			err = rows.Scan(&entityName, &ret)
			if err != nil {
				utils.AbortTui("Error getting create table statement: " + err.Error())
			}
		}

		ret += ";"
	} else if status == "deleted" {
		ret = "DROP TABLE IF EXISTS `" + entityName + "`;"
	} else if status == "modified" {
		utils.AbortTui("UNIMPLEMENTED")
	}

	return ret
}

// generateSqlForTablesPostgres is responsible for generating SQL for tables in Postgres.
func generateSqlForTablesPostgres(firstDb *sql.DB, secondDb *sql.DB, entityType string, entityName string, status string) string {
	var ret string

	utils.AbortTui("UNIMPLEMENTED")

	return ret
}
