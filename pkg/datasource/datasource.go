package datasource

import (
	"database/sql"
	"errors"
	"github.com/Okira-E/patchi/pkg/safego"
	"github.com/Okira-E/patchi/pkg/types"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func GetDataSource(dbInfo *types.DbConnectionInfo) (*sql.DB, safego.Option[error]) {
	connStr, errOpt := dbInfo.GetConnectionString()
	if errOpt.IsSome() {
		return nil, safego.Some(errors.New(errOpt.Unwrap()))
	}

	db, err := sql.Open(dbInfo.Dialect, connStr)
	if err != nil {
		return nil, safego.Some(err)
	}

	return db, safego.None[error]()
}
