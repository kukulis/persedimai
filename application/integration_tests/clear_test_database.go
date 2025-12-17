package integration_tests

import (
	"darbelis.eu/persedimai/internal/database"
	"fmt"
	"strings"
)

func ClearTestDatabase(database *database.Database, tableNames ...string) bool {
	databaseName := database.GetDatabaseName()
	if !strings.HasPrefix(databaseName, "test") {
		return false
	}

	conn, err := database.GetConnection()
	if err != nil {
		fmt.Println(err)
		return false
	}

	for _, tableName := range tableNames {
		sql := fmt.Sprintf("truncate table `%s`", tableName)

		_, err = conn.Exec(sql)

		if err != nil {
			fmt.Println(err)
		}
	}

	return true
}
