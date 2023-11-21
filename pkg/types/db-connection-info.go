package types

import (
	"database/sql"
	"errors"
	"github.com/Okira-E/patchi/pkg/utils"
	"github.com/Okira-E/patchi/safego"
	"strconv"

	_ "github.com/cockroachdb/cockroach-go/v2/crdb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type DbConnectionInfo struct {
	Dialect      string `json:"dialect,omitempty"`
	Name         string `json:"name,omitempty"`
	Host         string `json:"host,omitempty"`
	Port         int    `json:"port,omitempty"`
	User         string `json:"user,omitempty"`
	Password     string `json:"password,omitempty"`
	DatabaseName string `json:"database,omitempty"`
}

// GetConnectionString retrieves the valid sql connection string for the current dialect. It returns
// an error if the dialect isn't supported.
func (self *DbConnectionInfo) GetConnectionString() (string, safego.Option[string]) {
	if self.Dialect == "mysql" || self.Dialect == "mariadb" {
		return self.User + ":" + self.Password + "@tcp(" + self.Host + ":" + strconv.Itoa(self.Port) + ")/" + self.DatabaseName, safego.None[string]()
	} else if self.Dialect == "postgres" || self.Dialect == "cockroachdb" {
		return "postgresql://" + self.User + ":" + self.Password + "@" + self.Host + ":" + strconv.Itoa(self.Port) + "/" + self.DatabaseName + "?sslmode=disable", safego.None[string]()
	}

	return "", safego.Some[string]("Invalid dialect")
}

// Connect connects to the database and returns the sql.DB object.
func (self *DbConnectionInfo) Connect() (*sql.DB, safego.Option[error]) {
	connStr, errOpt := self.GetConnectionString()
	if errOpt.IsSome() {
		return nil, safego.Some(errors.New(errOpt.Unwrap()))
	}

	driverName := utils.Ternary(self.Dialect == "cockroachdb", "postgres", self.Dialect)
	driverName = utils.Ternary(self.Dialect == "mariadb", "mysql", driverName)

	db, err := sql.Open(driverName, connStr)
	if err != nil {
		return nil, safego.Some(err)
	}

	return db, safego.None[error]()
}
