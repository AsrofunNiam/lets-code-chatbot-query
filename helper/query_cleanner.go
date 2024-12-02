package helper

import (
	"fmt"
	"strings"
)

func ExtractQuery(msg ...interface{}) string {
	test := ExtractSQLQuery(fmt.Sprint(msg...))
	cleanQuery := CleanQuery(fmt.Sprint(msg...))
	fmt.Println(test)
	fmt.Println(cleanQuery)
	return cleanQuery
}
func CleanQuery(query string) string {
	query = strings.TrimPrefix(query, "```sql")
	query = strings.TrimPrefix(query, "```")
	query = strings.TrimPrefix(query, " ```")
	query = strings.TrimSuffix(query, "```")
	query = strings.TrimSuffix(query, " ```")
	query = strings.TrimSpace(query)
	query = strings.ReplaceAll(query, "\n", " ")
	query = strings.ReplaceAll(query, "\r", "")

	query = strings.Trim(query, "`")
	query = strings.TrimSpace(query)

	return query
}

// ExtractSQLQuery
func ExtractSQLQuery(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "SELECT") || strings.HasPrefix(line, "select") {
			return line
		}
	}
	return ""
}
