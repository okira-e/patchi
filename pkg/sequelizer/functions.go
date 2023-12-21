
package sequelizer

import (
	"database/sql"

	"github.com/Okira-E/patchi/pkg/utils"
)

// GenerateSqlForFunctions is the interface for generating SQL for functions in general.
func GenerateSqlForFunctions(firstDb *sql.DB, secondDb *sql.DB, dialect string, functionName string, status string) string {
	var ret string

	if dialect == "mysql" || dialect == "mariadb" {
		ret = generateSqlForFunctionsMysql(firstDb, secondDb, functionName, status)
	} else if dialect == "postgres" || dialect == "cockroachdb" {
		utils.AbortTui("UNIMPLEMENTED")
	}

	return ret
}

// generateSqlForFunctionsMysql is responsible for generating SQL for functions in Mysql.
func generateSqlForFunctionsMysql(firstDb *sql.DB, secondDb *sql.DB, functionName string, status string) string {
	var ret string

	if status == "created" {
		rows, err := firstDb.Query("SHOW CREATE FUNCTION " + functionName)
		if err != nil {
			utils.AbortTui("Error getting create function statement: " + err.Error())
		}

		for rows.Next() {
			var void string
			err = rows.Scan(&void, &void, &ret, &void, &void, &void)
			if err != nil {
				utils.AbortTui("Error getting create function statement: " + err.Error())
			}
		}

		ret += ";"
	} else if status == "deleted" {
		ret = "DROP FUNCTION IF EXISTS `" + functionName + "`;"
	}

	return ret
}

