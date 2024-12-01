package repository

import (
	"github.com/AsrofunNiam/lets-code-chatbot-query/helper"
	"github.com/AsrofunNiam/lets-code-chatbot-query/model/domain"
	"gorm.io/gorm"
)

type ProductRepositoryImpl struct {
}

func NewProductRepository() ProductRepository {
	return &ProductRepositoryImpl{}
}

func (repository *ProductRepositoryImpl) FindAll(db *gorm.DB, filters *map[string]string) domain.Products {
	products := domain.Products{}
	tx := db.Model(&domain.Product{})

	err := helper.ApplyFilter(tx, filters)
	helper.PanicIfError(err)

	err = tx.Find(&products).Error
	helper.PanicIfError(err)

	return products
}

func (repository *ProductRepositoryImpl) Create(db *gorm.DB, product *domain.Product) *domain.Product {
	err := db.Create(&product).Error
	helper.PanicIfError(err)
	return product
}
