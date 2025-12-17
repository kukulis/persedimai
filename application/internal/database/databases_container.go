package database

type DatabasesContainer interface {
	GetDatabase(env string) (*Database, error)
}
