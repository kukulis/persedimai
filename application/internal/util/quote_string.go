package util

import (
	"darbelis.eu/persedimai/internal/database"
	"fmt"
)

func QuoteString(s string) string {
	return fmt.Sprintf("'%s'", database.MysqlRealEscapeString(s))
}

func QuoteStringOrNull(s string) string {
	if s == "" {
		return "NULL"
	}
	
	return QuoteString(s)
}
