package datastruct

import (
	"time"
)

type Customer struct {
	Id        string    `"json": "id"`
	Name      string    `"json": "name"`
	Phone     string    `"json": "phone"`
	Address   string    `"json": "address"`
	BirthDate time.Time `"json": "birth"`
}

type Product struct {
	Id          string `"json": "id"`
	Description string `"json": "description"`
	Type        string `"json": "type"`
	SubType     string `"json": "subType"`
}

type TransactionProduct struct {
	Id       string  `"json": "id"`
	Product  Product `"json": "product"`
	Price    float64 `"json": "price"`
	Quantity int16   `"json": "quantity"`
	Total    float64 `"json": "total"`
}

type Transaction struct {
	Id       string               `"json": "id"`
	Date     time.Time            `"json": "date"`
	Customer Customer             `"json": "customer"`
	Items    []TransactionProduct `"json": "items"`
	Total    float64              `"json": "total"`
}
