package database

import (
	"errors"
	"github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Dbname   string
	DbType   string
}

func (dbConfig *DBConfig) InitializeFromEnvMap(envMap map[string]string) error {
	var ok bool

	dbConfig.DbType, ok = envMap["DBTYPE"]
	if !ok {
		return errors.New("DBTYPE environment variable not set")
	}

	dbConfig.Username, ok = envMap["DBUSER"]
	if !ok {
		return errors.New("DBUSER environment variable not set")
	}

	dbConfig.Password, ok = envMap["DBPASS"]
	if !ok {
		return errors.New("DBPASS environment variable not set")
	}

	dbConfig.Host, ok = envMap["DBHOST"]
	if !ok {
		return errors.New("DBHOST environment variable not set")
	}

	dbConfig.Port, ok = envMap["DBPORT"]
	if !ok {
		return errors.New("DBPORT environment variable not set")
	}

	dbConfig.Dbname, ok = envMap["DBNAME"]
	if !ok {
		return errors.New("DBNAME environment variable not set")
	}

	return nil
}

func (dbConfig *DBConfig) FormatDsn() (string, error) {
	if dbConfig.DbType == "mysql" {
		config := mysql.NewConfig()

		config.User = dbConfig.Username
		config.Passwd = dbConfig.Password
		config.DBName = dbConfig.Dbname
		config.DBName = dbConfig.Dbname
		config.Net = "tcp"
		config.Addr = dbConfig.Host + ":" + dbConfig.Port
		config.ParseTime = true

		return config.FormatDSN(), nil
	}

	return "", errors.New("BuildDsn: Database type " + dbConfig.DbType + " not supported")
}

func (dbConfig *DBConfig) GetRequiredParamsNames() []string {
	return []string{
		"DBTYPE",
		"DBUSER",
		"DBPASS",
		"DBHOST",
		"DBPORT",
		"DBNAME",
	}
}
