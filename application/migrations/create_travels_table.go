package migrations

import "darbelis.eu/persedimai/database"

func CreateTravelsTable(db *database.Database) error {
	conn, err := db.GetConnection()
	if err != nil {
		panic(err)
	}

	defer func() { _ = db.CloseConnection() }()

	sql := `create or replace table travels (
		id varchar(64) not null primary key,
		from_point varchar(64) not null,
		to_point varchar(64) not null,
		departure datetime not null,
		arrival datetime not null,
		index idx_from_departure (from_point, departure),
		index idx_to_arrival (to_point, arrival)
	)`

	_, err = conn.Exec(sql)

	return err
}
