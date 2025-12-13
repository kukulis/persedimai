package dao

import (
	"darbelis.eu/persedimai/database"
	"darbelis.eu/persedimai/tables"
	"errors"
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
	connection, err := pointDao.database.GetConnection()
	if err != nil {
		return err
	}

	lines := make([]string, len(points))
	for i, point := range points {
		line := fmt.Sprintf("('%s', %f,%f,'%s')", point.ID, point.X, point.Y, database.MysqlRealEscapeString(point.Name))
		lines[i] = line
	}

	valuesSubSql := strings.Join(lines, ",\n")

	sql := "insert into points (ID, x,y,name) values " + valuesSubSql

	_, err = connection.Exec(sql)

	if err != nil {
		return errors.New(err.Error() + " for sql " + sql)
	}

	return nil
}

func (pointDao *PointDao) SelectAll() ([]*tables.Point, error) {
	connection, err := pointDao.database.GetConnection()
	if err != nil {
		return nil, err
	}

	sql := "SELECT id, x, y, name FROM points"
	rows, err := connection.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []*tables.Point
	for rows.Next() {
		point := &tables.Point{}
		err := rows.Scan(&point.ID, &point.X, &point.Y, &point.Name)
		if err != nil {
			return nil, err
		}
		points = append(points, point)
	}

	return points, nil
}

func (pointDao *PointDao) Count() (int, error) {
	connection, err := pointDao.database.GetConnection()
	if err != nil {
		return 0, err
	}

	sql := "SELECT COUNT(*) FROM points"
	var count int
	err = connection.QueryRow(sql).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (pointDao *PointDao) UpsertMany(points []*tables.Point) {
	// TODO
}
