package types

import "fmt"

type UserConfig struct {
	DbConnections map[string]*DbConnectionInfo `json:"db_connections,omitempty"`
}

func (uc *UserConfig) String() string {
	var str string

	for _, dbConfig := range uc.DbConnections {
		str += fmt.Sprintf("%s\n", dbConfig)
	}

	return str
}
