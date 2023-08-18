package types

import "github.com/Okira-E/patchi/pkg/safego"

type DbConnectionInfo struct {
	Dialect  string `json:"dialect,omitempty"`
	Name     string `json:"name,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Database string `json:"database,omitempty"`
}

func (d *DbConnectionInfo) GetConnectionString() (string, safego.Option[string]) {
	if d.Dialect == "mysql" {
		return d.User + ":" + d.Password + "@(" + d.Host + ":" + string(rune(d.Port)) + ")/" + d.Database, safego.None[string]()
	} else if d.Dialect == "postgres" {
		return "postgres://" + d.User + ":" + d.Password + "@" + d.Host, safego.None[string]()
	} else if d.Dialect == "mssql" {
		return "sqlserver://" + d.User + ":" + d.Password + "@" + d.Host, safego.None[string]()
	}

	return "", safego.Some[string]("Invalid dialect")
}
