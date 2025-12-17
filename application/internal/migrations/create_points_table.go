package migrations

import "darbelis.eu/persedimai/internal/database"

func CreatePointsTable(db *database.Database) error {
	conn, err := db.GetConnection()
	if err != nil {
		panic(err)
	}

	defer func() { _ = db.CloseConnection() }()

	sql := `create or replace table points (
		id varchar(64) not null primary key,
		 name varchar (128),  
		 x decimal(12,5),
		 y decimal(12,5))`

	_, err = conn.Exec(sql)

	return err
}
