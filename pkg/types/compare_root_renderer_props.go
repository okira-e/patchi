package types

import "database/sql"

// CompareRootRendererProps is the props for the CompareRootRenderer.
type CompareRootRendererProps struct {
	FirstDb  DbConnection
	SecondDb DbConnection
}

// DbConnection is the connection info for a database.
type DbConnection struct {
	Info          *DbConnectionInfo
	SqlConnection *sql.DB
}
