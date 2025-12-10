package migrations

import "darbelis.eu/persedimai/database"

func CreatePointsTable(db *database.Database) error {
	conn, err := db.GetConnection()
	if err != nil {
		panic(err)
	}

	defer func() { _ = db.CloseConnection() }()

	sql := `create table if not exists points (
		id varchar(32) not null primary key,
		 name varchar (128),  
		 x decimal(10,5),
		 y decimal(10,5))`

	_, err = conn.Exec(sql)

	return err
}
