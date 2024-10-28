// Copyright (c) 2024 Christopher Wolf. All rights reserved.
// This software is proprietary and confidential.
// Unauthorized copying of this file, via any medium, is strictly prohibited.

package models

import "encoding/xml"

// Transaction structure
type Transaction struct {
	XMLName       xml.Name `xml:"Transaction"`
	ID            string   `xml:"ID"`
	Timestamp     string   `xml:"Timestamp"`
	Amount        float64  `xml:"Amount"`
	Currency      string   `xml:"Currency"`
	Customer      Customer `xml:"Customer"`
	Items         []Item   `xml:"Items>Item"`
	Status        string   `xml:"Status"`
	PromotionCode *string  `xml:"PromotionCode,omitempty"`
	Discount      *float64 `xml:"Discount,omitempty"`
}

// Customer structure
type Customer struct {
	Name  string `xml:"Name"`
	Email string `xml:"Email"`
}

// Item structure
type Item struct {
	ProductID string  `xml:"ProductID"`
	Quantity  int     `xml:"Quantity"`
	Price     float64 `xml:"Price"`
}
