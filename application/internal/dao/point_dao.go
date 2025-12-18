package dao

import (
	"darbelis.eu/persedimai/internal/data"
	"darbelis.eu/persedimai/internal/database"
	"darbelis.eu/persedimai/internal/tables"
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

// FindByCoordinates finds a point by exact X and Y coordinates
func (pointDao *PointDao) FindByCoordinates(x, y float64) (*tables.Point, error) {
	connection, err := pointDao.database.GetConnection()
	if err != nil {
		return nil, err
	}

	sql := "SELECT id, x, y, name FROM points WHERE x = ? AND y = ?"
	var point tables.Point
	err = connection.QueryRow(sql, x, y).Scan(&point.ID, &point.X, &point.Y, &point.Name)
	if err != nil {
		return nil, err
	}

	return &point, nil
}

// SelectWithFilter returns points filtered by the given PointsFilter
func (pointDao *PointDao) SelectWithFilter(filter *data.PointsFilter) ([]*tables.Point, error) {
	connection, err := pointDao.database.GetConnection()
	if err != nil {
		return nil, err
	}

	// Build SQL query with conditions
	sql := "SELECT id, x, y, name FROM points"
	var conditions []string

	// Add X coordinate filter
	if filter.X != nil {
		conditions = append(conditions, fmt.Sprintf("x = %f", *filter.X))
	}

	// Add Y coordinate filter
	if filter.Y != nil {
		conditions = append(conditions, fmt.Sprintf("y = %f", *filter.Y))
	}

	// Add name part filter
	if filter.NamePart != "" {
		escapedName := database.MysqlRealEscapeString(filter.NamePart)
		conditions = append(conditions, fmt.Sprintf("name LIKE '%%%s%%'", escapedName))
	}

	// Add ID part filter
	if filter.IdPart != "" {
		escapedId := database.MysqlRealEscapeString(filter.IdPart)
		conditions = append(conditions, fmt.Sprintf("id LIKE '%%%s%%'", escapedId))
	}

	// Append WHERE clause if there are conditions
	if len(conditions) > 0 {
		sql += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Add limit
	if filter.Limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d", filter.Limit)
	}

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

// PointBounds represents the min/max coordinate boundaries
type PointBounds struct {
	MinX float64
	MaxX float64
	MinY float64
	MaxY float64
}

// GetMinMaxX returns the minimum and maximum X coordinate values
func (pointDao *PointDao) GetMinMaxX() (minX, maxX float64, err error) {
	connection, err := pointDao.database.GetConnection()
	if err != nil {
		return 0, 0, err
	}

	sql := "SELECT MIN(x), MAX(x) FROM points"
	err = connection.QueryRow(sql).Scan(&minX, &maxX)
	if err != nil {
		return 0, 0, err
	}

	return minX, maxX, nil
}

// GetMinMaxY returns the minimum and maximum Y coordinate values
func (pointDao *PointDao) GetMinMaxY() (minY, maxY float64, err error) {
	connection, err := pointDao.database.GetConnection()
	if err != nil {
		return 0, 0, err
	}

	sql := "SELECT MIN(y), MAX(y) FROM points"
	err = connection.QueryRow(sql).Scan(&minY, &maxY)
	if err != nil {
		return 0, 0, err
	}

	return minY, maxY, nil
}

// GetBounds returns all coordinate boundaries in a single query (more efficient)
func (pointDao *PointDao) GetBounds() (*PointBounds, error) {
	connection, err := pointDao.database.GetConnection()
	if err != nil {
		return nil, err
	}

	sql := "SELECT MIN(x), MAX(x), MIN(y), MAX(y) FROM points"
	bounds := &PointBounds{}
	err = connection.QueryRow(sql).Scan(&bounds.MinX, &bounds.MaxX, &bounds.MinY, &bounds.MaxY)
	if err != nil {
		return nil, err
	}

	return bounds, nil
}
