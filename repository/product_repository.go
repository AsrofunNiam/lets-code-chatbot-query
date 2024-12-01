package repository

import (
	"github.com/AsrofunNiam/lets-code-chatbot-query/model/domain"
	"gorm.io/gorm"
)

type ProductRepository interface {
	FindAll(db *gorm.DB, filters *map[string]string) domain.Products
	Create(db *gorm.DB, Product *domain.Product) *domain.Product
}
