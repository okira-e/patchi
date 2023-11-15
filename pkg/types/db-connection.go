package types

import "database/sql"

// DbConnection is the connection info for a database.
type DbConnection struct {
	Info          *DbConnectionInfo
	SqlConnection *sql.DB
}
