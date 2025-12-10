package dao

import (
	"darbelis.eu/persedimai/database"
	"darbelis.eu/persedimai/tables"
	"fmt"
	"strings"
)

type PointDao struct {
	database *database.Database
}

func NewPointDao(database *database.Database) *PointDao {
	return &PointDao{database: database}
}

func (pointDao *PointDao) InsertMany(points []*tables.Point) error {

	connection := pointDao.database.GetConnection()

	lines := make([]string, len(points))
	for i, point := range points {
		line := fmt.Sprintf("(%s, %f,%f,'%s')", point.ID, point.X, point.Y, database.MysqlRealEscapeString(point.Name))
		lines[i] = line
	}

	valuesSubSql := strings.Join(lines, ",\n")

	sql := "insert into points (ID, x,y,name) values " + valuesSubSql

	_, err := connection.Exec(sql)

	return err
}

func (pointDao *PointDao) UpsertMany(points []*tables.Point) {
	// TODO
}
