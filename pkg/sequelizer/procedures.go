package sequelizer

import (
	"database/sql"

	"github.com/Okira-E/patchi/pkg/utils"
)

// GenerateSqlForProcedures is the interface for generating SQL for procedures in general.
func GenerateSqlForProcedures(firstDb *sql.DB, secondDb *sql.DB, dialect string, procedureName string, status string) string {
	var ret string

	if dialect == "mysql" || dialect == "mariadb" {
		ret = generateSqlForProceduresMysql(firstDb, secondDb, procedureName, status)
	} else if dialect == "postgres" || dialect == "cockroachdb" {
		utils.AbortTui("UNIMPLEMENTED")
	}

	return ret
}

// generateSqlForProceduresMysql is responsible for generating SQL for procedures in Mysql.
func generateSqlForProceduresMysql(firstDb *sql.DB, secondDb *sql.DB, procedureName string, status string) string {
	var ret string

	if status == "created" {
		rows, err := firstDb.Query("SHOW CREATE PROCEDURE " + procedureName)
		if err != nil {
			utils.AbortTui("Error getting create procedure statement: " + err.Error())
		}

		for rows.Next() {
			var void string
			err = rows.Scan(&void, &void, &ret, &void, &void, &void)
			if err != nil {
				utils.AbortTui("Error getting create procedure statement: " + err.Error())
			}
		}

		ret += ";"
	} else if status == "deleted" {
		ret = "DROP PROCEDURE IF EXISTS `" + procedureName + "`;"
	}

	return ret
}

