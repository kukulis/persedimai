package migrations

import (
	"darbelis.eu/persedimai/internal/database"
	"errors"
	"fmt"
	"log"
	"time"
)

type ClustersCreator struct {
	db *database.Database
}

func NewClustersCreator(db *database.Database) *ClustersCreator {
	return &ClustersCreator{db: db}
}

func (creator *ClustersCreator) CreateClustersTableSQL(clustersTableNumber int) string {
	sql :=
		fmt.Sprintf(`create or replace table clustered_arrival_travels%d (
		travel_id varchar(64) not null,
		from_point varchar(64) not null,
		to_point varchar(64) not null,
		departure_cl int,
		arrival_cl int,
		index idx_from_departure_cl (from_point, departure_cl),
		index idx_to_arrival_cl (to_point, arrival_cl) )`, clustersTableNumber)

	return sql
}

func (creator *ClustersCreator) InsertClustersDataSQLs(clustersTableNumber int) []string {

	sqlDisableKeys := fmt.Sprintf(`ALTER TABLE clustered_arrival_travels%d DISABLE KEYS`, clustersTableNumber)

	fromTable := "travels"
	if clustersTableNumber > 2 {
		fromTable = fmt.Sprintf("clustered_arrival_travels%d", clustersTableNumber/2)
	}

	idField := "t.id"
	if clustersTableNumber > 2 {
		idField = "t.travel_id"
	}

	sqlInsert1 := fmt.Sprintf(`insert into clustered_arrival_travels%d
		select %s, t.from_point, t.to_point, t.departure_cl, t.arrival_cl
			from %s t`, clustersTableNumber, idField, fromTable)

	sqlInsert2 := fmt.Sprintf(`insert into clustered_arrival_travels%d
		select %s, t.from_point, t.to_point, t.departure_cl, t.arrival_cl+%d
			from %s t`, clustersTableNumber, idField, clustersTableNumber/2, fromTable)

	sqlEnableKeys := fmt.Sprintf(`ALTER TABLE clustered_arrival_travels%d ENABLE KEYS`, clustersTableNumber)

	return []string{sqlDisableKeys, sqlInsert1, sqlInsert2, sqlEnableKeys}
}

func (creator *ClustersCreator) CreateClustersTables() error {
	dbConn, err := creator.db.GetConnection()
	if err != nil {
		return err
	}
	var i = 2
	for i <= 32 {
		sql := creator.CreateClustersTableSQL(i)
		_, err := dbConn.Exec(sql)
		if err != nil {
			return errors.New("failed to create clusters : " + err.Error())
		}

		i = i * 2
	}
	return nil
}

func (creator *ClustersCreator) InsertClustersDatas() error {
	dbConn, err := creator.db.GetConnection()
	if err != nil {
		return err
	}
	var i = 2
	for i <= 32 {
		sqls := creator.InsertClustersDataSQLs(i)
		for _, sql := range sqls {
			log.Println("Running sql : " + sql)
			sqlStart := time.Now()
			_, err := dbConn.Exec(sql)
			if err != nil {
				return errors.New("failed to create clusters : " + err.Error())
			}
			sqlEnd := time.Now()

			duration := sqlEnd.Sub(sqlStart)
			log.Printf("sql execution duration %s", duration.String())
		}

		i = i * 2
	}
	return nil
}

func (creator *ClustersCreator) UpdateClustersOnTravels() error {
	dbConn, err := creator.db.GetConnection()
	if err != nil {
		return err
	}

	sql := `update travels set
                   departure_cl = floor(unix_timestamp(departure) / 3600),
                   arrival_cl = floor(unix_timestamp(arrival) / 3600)`

	log.Println("Running sql : " + sql)

	sqlStart := time.Now()

	_, err = dbConn.Exec(sql)
	sqlEnd := time.Now()

	duration := sqlEnd.Sub(sqlStart)
	log.Printf("sql execution duration %s", duration.String())

	return err
}
