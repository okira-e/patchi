package mysql

import (
	"github.com/Okira-E/patchi/pkg/types"
	"github.com/Okira-E/patchi/pkg/utils"
)

func RecreateSqlForTables(firstDb types.DbConnection, secondDb types.DbConnection, entityType string, entityName string, status string) string {
	var ret string

	if status == "created" {
		rows, err := firstDb.SqlConnection.Query("SHOW CREATE TABLE " + entityName + ";")
		if err != nil {
			utils.AbortTui("Error getting create table statement: " + err.Error())
		}

		for rows.Next() {
			err = rows.Scan(&entityName, &ret)
			if err != nil {
				utils.AbortTui("Error getting create table statement: " + err.Error())
			}
		}
	} else if status == "deleted" {
		ret = "DROP TABLE IF EXISTS `" + entityName + "`;"
	} else if status == "modified" {
		// UNIMPLEMENTED
	}

	return ret
}
