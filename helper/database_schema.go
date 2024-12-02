package helper

import (
	"fmt"

	"github.com/AsrofunNiam/lets-code-chatbot-query/model/domain"
	"gorm.io/gorm"
)

func GetDatabaseSchema(tx *gorm.DB) []domain.SchemaTable {
	var columns []domain.SchemaTable
	query := `SELECT TABLE_NAME as table_name, COLUMN_NAME as column_name, DATA_TYPE as data_type
			  FROM INFORMATION_SCHEMA.COLUMNS 
			  WHERE TABLE_SCHEMA = 'voca_game'`
	tx.Raw(query).Scan(&columns)
	return columns
}
func GenerateSchemaDescriptions(tx *gorm.DB) []string {
	databaseSchema := GetDatabaseSchema(tx)
	var descriptions []string
	for _, col := range databaseSchema {
		desc := fmt.Sprintf("Table: %s, Column: %s, Type: %s", col.TableName, col.ColumnName, col.DataType)
		descriptions = append(descriptions, desc)
	}

	return descriptions
}
