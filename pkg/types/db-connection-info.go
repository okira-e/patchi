package types

import (
	"strconv"

	"github.com/Okira-E/patchi/pkg/safego"
)

type DbConnectionInfo struct {
	Dialect  string `json:"dialect,omitempty"`
	Name     string `json:"name,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Database string `json:"database,omitempty"`
}

// GetConnectionString retrieves the valid sql connection string for the current dialect. It returns 
// an error if the dialect isn't supported.
func (self *DbConnectionInfo) GetConnectionString() (string, safego.Option[string]) {
	if self.Dialect == "mysql" {
		return self.User + ":" + self.Password + "@tcp(" + self.Host + ":" + strconv.Itoa(self.Port) + ")/" + self.Database, safego.None[string]()
	} else if self.Dialect == "postgres" {
		return "postgres://" + self.User + ":" + self.Password + "@" + self.Host, safego.None[string]()
	} else if self.Dialect == "mssql" {
		return "sqlserver://" + self.User + ":" + self.Password + "@" + self.Host, safego.None[string]()
	}

	return "", safego.Some[string]("Invalid dialect")
}
