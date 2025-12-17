package di

import "darbelis.eu/persedimai/internal/database"

type DatabasesMapContainer struct {
	databases map[string]*database.Database
}

func NewDatabasesMapContainer() *DatabasesMapContainer {
	return &DatabasesMapContainer{
		databases: make(map[string]*database.Database),
	}
}

func (d DatabasesMapContainer) GetDatabase(env string) (*database.Database, error) {
	if db, ok := d.databases[env]; ok {
		return db, nil
	}

	db, err := NewDatabase(env)
	if err != nil {
		return nil, err
	}

	d.databases[env] = db

	return db, nil
}
