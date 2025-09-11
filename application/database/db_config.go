package database

import (
	"darbelis.eu/persedimai/env"
)

type DBConfig struct {
	username string
	password string
	host     string
	port     string
	dbname   string
}

func (D *DBConfig) Username() string {
	return D.username
}

func (D *DBConfig) SetUsername(username string) {
	D.username = username
}

func (D *DBConfig) Password() string {
	return D.password
}

func (D *DBConfig) SetPassword(password string) {
	D.password = password
}

func (D *DBConfig) Host() string {
	return D.host
}

func (D *DBConfig) SetHost(host string) {
	D.host = host
}

func (D *DBConfig) Port() string {
	return D.port
}

func (D *DBConfig) SetPort(port string) {
	D.port = port
}

func (D *DBConfig) Dbname() string {
	return D.dbname
}

func (D *DBConfig) SetDbname(dbname string) {
	D.dbname = dbname
}

func (config *DBConfig) LoadFromEnv(filePath string) {
	envFile, _ := env.EnvMap{}.Read(filePath)
	config.SetUsername(envFile.Getenv("DBUSER"))
	config.SetPassword(envFile.Getenv("DBPASS"))
	config.SetHost(envFile.Getenv("DBHOST"))
	config.SetPort(envFile.Getenv("DBPORT"))
	config.SetDbname(envFile.Getenv("DBNAME"))
}
