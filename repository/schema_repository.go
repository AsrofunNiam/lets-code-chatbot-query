package repository

import (
	"github.com/AsrofunNiam/lets-code-chatbot-query/model/domain"
	"gorm.io/gorm"
)

type SchemaRepository interface {
	FindAll(tx *gorm.DB) domain.SchemaTables
}
