package integration_tests

import (
	"darbelis.eu/persedimai/di"
	"darbelis.eu/persedimai/internal/dao"
	"darbelis.eu/persedimai/internal/generator"
	"darbelis.eu/persedimai/internal/migrations"
	"fmt"
	"testing"
)

func TestGeneratePointsToDB(t *testing.T) {
	gf := generator.GeneratorFactory{}

	var idGenerator generator.IdGenerator
	idGenerator = &generator.SimpleIdGenerator{}
	g := gf.CreateGenerator(5, 1000, 0, idGenerator)

	db, err := di.NewDatabase("test")
	if err != nil {
		t.Error(err)
	}

	err = migrations.CreatePointsTable(db)
	if err != nil {
		t.Fatal(err)
	}

	if !ClearTestDatabase(db, "points") {
		fmt.Println("Failed to clear database")
	}

	pointDao := dao.NewPointDao(db)
	pointDbConsumer := generator.NewPointConsumer(pointDao, 10)
	g.GeneratePoints(pointDbConsumer)

	err = pointDbConsumer.Flush()
	if err != nil {
		t.Error(err)
	}
}
