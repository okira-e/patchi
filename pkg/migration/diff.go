package migration

import (
	"database/sql"
	"log"

	"github.com/Okira-E/patchi/pkg/types"
)

func TableDiff(firstDb *sql.DB, secondDb *sql.DB) []types.TableDiff {
	rows, err := firstDb.Query("SHOW TABLES")
	if err != nil {
		log.Fatalf("Error querying database: %s", err.Error())
	}

	tables := []string{}
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			log.Fatalf("Error scanning row: %s", err.Error())
		}

		tables = append(tables, table)
	}

	return []types.TableDiff{}
}
