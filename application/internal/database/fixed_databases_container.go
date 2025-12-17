package database

type FixedDatabasesContainer struct {
	database *Database
}

func (f FixedDatabasesContainer) GetDatabase(env string) (*Database, error) {
	return f.database, nil
}
