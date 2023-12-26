package sequelizer

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Okira-E/patchi/pkg/types"
	"github.com/Okira-E/patchi/pkg/utils"
	"github.com/Okira-E/patchi/safego"
)

// GenerateSqlForColumns generates the SQL for a column based on it's status (created or deleted.)
func GenerateSqlForColumns(firstDb types.DbConnection, dialect string, columnName string, tableName string, status string) (string, safego.Option[string]) {
	var ret string
	errOpt := safego.None[string]()

	if dialect == "mysql" || dialect == "mariadb" {
		ret, errOpt = generateSqlForColumnsMysql(firstDb, columnName, tableName, status)
	} else if dialect == "postgres" || dialect == "cockroachdb" {
		utils.AbortTui("UNIMPLEMENTED")
	}

	return ret, errOpt
}

// UNFINISHED
func generateSqlForColumnsMysql(firstDb types.DbConnection, columnName string, tableName string, status string) (string, safego.Option[string]) {
	var result string

	if status == "deleted" {
		result = "ALTER TABLE " + tableName + " DROP COLUMN " + columnName
	} else if status == "created" {
		query := `
			SELECT
			    C.ORDINAL_POSITION,
			    C.COLUMN_DEFAULT,
			    C.EXTRA,
			    C.COLUMN_KEY,
			    C.IS_NULLABLE,
			    C.COLUMN_TYPE,
			    KCU.REFERENCED_TABLE_NAME,
			    KCU.REFERENCED_COLUMN_NAME
			FROM
			    information_schema.COLUMNS C
			LEFT JOIN
			    information_schema.KEY_COLUMN_USAGE KCU
			ON
			    C.TABLE_SCHEMA = KCU.TABLE_SCHEMA
			    AND C.TABLE_NAME = KCU.TABLE_NAME
			    AND C.COLUMN_NAME = KCU.COLUMN_NAME
			WHERE
			    C.TABLE_SCHEMA = ?
			    AND C.TABLE_NAME = ?
				AND C.COLUMN_NAME = ?
		`

		rows, err := firstDb.SqlConnection.Query(query, firstDb.Info.DatabaseName, tableName, columnName)
		if err != nil {
			return result, safego.Some("Failed to get info on column: " + err.Error())
		}

		var ordinalPos sql.NullInt32
		var columnDefault sql.NullString
		var columnExtra sql.NullString
		var columnKey sql.NullString
		var isNullable sql.NullString
		var columnType sql.NullString
		var referencedTableName sql.NullString
		var referencedColumnName sql.NullString
		for rows.Next() {
			err = rows.Scan(&ordinalPos, &columnDefault, &columnExtra, &columnKey, &isNullable, &columnType, &referencedTableName, &referencedColumnName)
			if err != nil {
				return result, safego.Some("Failed to scan results on column: " + err.Error())
			}
		}

		query = `
    		SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND ORDINAL_POSITION = ?
		`

		rows, err = firstDb.SqlConnection.Query(query, firstDb.Info.DatabaseName, tableName, ordinalPos.Int32-1)
		if err != nil {
			return result, safego.Some("Failed to get previous column name: " + err.Error())
		}

		var prevColumnName sql.NullString
		for rows.Next() {
			err = rows.Scan(&prevColumnName)
			if err != nil {
				return result, safego.Some("Failed to scan previous column name: " + err.Error())
			}
		}

		// Golang parsing sucks.
		result = "ALTER TABLE " + tableName + " ADD COLUMN " + columnName + " " + columnType.String + " " + utils.Ternary(isNullable.String == "YES", "NULL", "NOT NULL") + " "
		result += utils.Ternary(columnDefault.Valid, "DEFAULT "+columnDefault.String, "") + " "
		result += utils.Ternary(columnExtra.Valid, strings.ReplaceAll(columnExtra.String, "DEFAULT_GENERATED", ""), "") + " "
		result += utils.Ternary(columnKey.String == "PRI", "PRIMARY KEY", "") + " "
		// Handle if the column is a foreign key to a different table.
		// NOTE: This doesn't guarantee that the table it references exists in the other database env.
		if referencedTableName.Valid && referencedColumnName.Valid {
			result += fmt.Sprintf("REFERENCES %s(%s) ", referencedTableName.String, referencedColumnName.String)
		}

		result += utils.Ternary(prevColumnName.Valid, "AFTER "+prevColumnName.String, "FIRST") + " "
	}

	// Remove long spaces to make the query look nicer.
	ret := ""
	for i, word := range strings.Fields(result) { // `Fields()` returns a slice of words, nullifying any spaces or tabs between them.
		if i == 0 {
			ret += word
			continue
		}

		ret += " " + word
	}
	ret += ";"

	return ret, safego.None[string]()
}
