package sequelizer

import (
	"database/sql"

	"github.com/Okira-E/patchi/pkg/utils"
)

// GenerateSqlForViews is the interface for generating SQL for views in general.
func GenerateSqlForViews(firstDb *sql.DB, dialect string, viewName string, status string) string {
	var ret string

	if dialect == "mysql" || dialect == "mariadb" {
		ret = generateSqlForViewsMysql(firstDb, viewName, status)
	} else if dialect == "postgres" || dialect == "cockroachdb" {
		utils.AbortTui("UNIMPLEMENTED")
	}

	return ret
}

// generateSqlForViewsMysql is responsible for generating SQL for views in Mysql.
func generateSqlForViewsMysql(firstDb *sql.DB, viewName string, status string) string {
	var ret string

	if status == "created" {
		rows, err := firstDb.Query("SHOW CREATE VIEW " + viewName)
		if err != nil {
			utils.AbortTui("Error getting create view statement: " + err.Error())
		}

		for rows.Next() {
			var void string
			err = rows.Scan(&viewName, &ret, &void, &void)
			if err != nil {
				utils.AbortTui("Error getting create view statement: " + err.Error())
			}
		}

		ret += ";"
	} else if status == "deleted" {
		ret = "DROP view IF EXISTS `" + viewName + "`;"
	}

	return ret
}
