package models

type Announcement struct {
	Id             int64    `json:"id"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	IsNew          bool     `json:"is_new"`
	Price          float32  `json:"price"`
	IsExchangeable bool     `json:"is_exchangeable"`
	IsActive       bool     `json:"is_active"`
	UserId         int64    `json:"user_id"`
	Images         []string `json:"images"`
	PaymentMethods `json:"payment_methods"`
}
