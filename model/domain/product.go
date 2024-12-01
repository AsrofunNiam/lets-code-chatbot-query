package domain

import (
	"github.com/AsrofunNiam/lets-code-chatbot-query/model/web"
	"gorm.io/gorm"
)

type Products []Product
type Product struct {
	gorm.Model
	Name        string  `gorm:"type:varchar(255);not null"`
	Description string  `gorm:"type:text"`
	Price       float64 `gorm:"type:decimal(10,2);not null" `
	Available   bool    `gorm:"default:true"`
}

func (product *Product) ToProductResponse() web.ProductResponse {
	return web.ProductResponse{
		// Required Fields
		ID: product.ID,
		// Fields
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Available:   product.Available,
	}
}

func (users Products) ToProductResponses() []web.ProductResponse {
	productResponses := []web.ProductResponse{}
	for _, user := range users {
		productResponses = append(productResponses, user.ToProductResponse())
	}
	return productResponses
}
