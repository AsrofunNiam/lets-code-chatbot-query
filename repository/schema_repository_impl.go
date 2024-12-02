package repository

import (
	"github.com/AsrofunNiam/lets-code-chatbot-query/model/domain"
	"gorm.io/gorm"
)

type SchemaRepositoryImpl struct {
}

func NewSchemaRepository() SchemaRepository {
	return &SchemaRepositoryImpl{}
}

func (repository *SchemaRepositoryImpl) FindAll(tx *gorm.DB) domain.SchemaTables {
	schemaTable := domain.SchemaTables{}
	query := `SELECT TABLE_NAME as table_name, COLUMN_NAME as column_name, DATA_TYPE as data_type
			  FROM INFORMATION_SCHEMA.COLUMNS 
			  WHERE TABLE_SCHEMA = 'voca_game'`
	tx.Raw(query).Scan(&schemaTable)
	return schemaTable
}
