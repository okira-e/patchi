package sequelizer

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Okira-E/patchi/pkg/utils"
)

// GenerateSqlForTables is the interface for generating SQL for tables in general.
func GenerateSqlForTables(firstDb *sql.DB, dialect string, entityName string, status string) string {
	var ret string

	if dialect == "mysql" || dialect == "mariadb" {
		ret = generateSqlForTablesMysql(firstDb, entityName, status)
	} else if dialect == "postgres" || dialect == "cockroachdb" {
		utils.AbortTui("UNIMPLEMENTED")
	}

	return ret
}

// generateSqlForTablesMysql is responsible for generating SQL for tables in Mysql.
func generateSqlForTablesMysql(firstDb *sql.DB, entityName string, status string) string {
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
	}

	return ret
}

// generateSqlForTablesPostgres AI GENERATED DRAFT THAT ISN'T COMPLETE/CORRECT.
func generateSqlForTablesPostgres(firstDb *sql.DB, entityName string, status string) string {
	var ret string

	if status == "created" {
		// Query to retrieve column and key information
		rows, err := firstDb.Query(`
		SELECT columns.column_name, data_type, character_maximum_length, is_nullable,
		       column_default, columns.ordinal_position,
		       constraint_type, constraint_name, referenced_table_name, referenced_column_name
		FROM information_schema.columns
		LEFT JOIN information_schema.key_column_usage
		ON columns.column_name = key_column_usage.column_name
		WHERE table_name = $1
		ORDER BY ordinal_position;
	`, entityName)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Collect column and key information
		var columns []string
		var primaryKeyCols []string
		var foreignKeyCols []string

		for rows.Next() {
			var columnName, dataType, isNullable, columnDefault, constraintType, constraintName sql.NullString
			var characterMaximumLength, ordinalPosition sql.NullInt64
			var referencedTableName, referencedColumnName sql.NullString

			if err := rows.Scan(&columnName, &dataType, &characterMaximumLength, &isNullable,
				&columnDefault, &ordinalPosition, &constraintType, &constraintName,
				&referencedTableName, &referencedColumnName); err != nil {
				log.Fatal(err)
			}

			columnDef := fmt.Sprintf("%s %s", columnName.String, dataType.String)
			if characterMaximumLength.Valid {
				columnDef += fmt.Sprintf("(%d)", characterMaximumLength.Int64)
			}
			if isNullable.String == "NO" {
				columnDef += " NOT NULL"
			}
			if columnDefault.Valid {
				columnDef += fmt.Sprintf(" DEFAULT %s", columnDefault.String)
			}

			columns = append(columns, columnDef)

			if constraintType.Valid {
				if constraintType.String == "PRIMARY KEY" {
					primaryKeyCols = append(primaryKeyCols, columnName.String)
				} else if constraintType.String == "FOREIGN KEY" {
					foreignKeyCols = append(foreignKeyCols, fmt.Sprintf("%s REFERENCES %s(%s)",
						columnName.String, referencedTableName.String, referencedColumnName.String))
				}
			}
		}

		// Check for errors from iterating over rows.
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		// Create the CREATE TABLE statement
		createTableStatement := fmt.Sprintf("CREATE TABLE %s (\n\t%s", entityName, strings.Join(columns, ",\n\t"))

		// Add primary key constraint
		if len(primaryKeyCols) > 0 {
			createTableStatement += fmt.Sprintf(",\n\tPRIMARY KEY (%s)", strings.Join(primaryKeyCols, ", "))
		}

		// Add foreign key constraints
		if len(foreignKeyCols) > 0 {
			createTableStatement += fmt.Sprintf(",\n\t%s", strings.Join(foreignKeyCols, ",\n\t"))
		}

		createTableStatement += "\n);"
	} else if status == "deleted" {
		ret = "DROP TABLE IF EXISTS `" + entityName + "`;"
	}

	return ret
}
