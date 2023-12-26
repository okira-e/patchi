package sequelizer

import (
	"database/sql"

	"github.com/Okira-E/patchi/pkg/utils"
)

// GenerateSqlForTriggers is the interface for generating SQL for triggers in general.
func GenerateSqlForTriggers(firstDb *sql.DB, secondDb *sql.DB, dialect string, triggerName string, status string) string {
	var ret string

	if dialect == "mysql" || dialect == "mariadb" {
		ret = generateSqlForTriggersMysql(firstDb, secondDb, triggerName, status)
	} else if dialect == "postgres" || dialect == "cockroachdb" {
		utils.AbortTui("UNIMPLEMENTED")
	}

	return ret
}

// generateSqlForTriggersMysql is responsible for generating SQL for triggers in Mysql.
func generateSqlForTriggersMysql(firstDb *sql.DB, secondDb *sql.DB, triggerName string, status string) string {
	var ret string

	if status == "created" {
		rows, err := firstDb.Query("SHOW CREATE TRIGGER " + triggerName)
		if err != nil {
			utils.AbortTui("Error getting create trigger statement: " + err.Error())
		}

		for rows.Next() {
			var void string
			err = rows.Scan(&void, &void, &ret, &void, &void, &void, &void)
			if err != nil {
				utils.AbortTui("Error getting create trigger statement: " + err.Error())
			}
		}

		ret += ";"
	} else if status == "deleted" {
		ret = "DROP TRIGGER IF EXISTS `" + triggerName + "`;"
	}

	return ret
}
