package dtos

type AnnouncementDTO struct {
	Title             string   `json:"title"`
	Description       string   `json:"description"`
	IsNew             bool     `json:"is_new"`
	Price             float32  `json:"price"`
	IsExchangeable    bool     `json:"is_exchangeable"`
	IsActive          bool     `json:"is_active"`
	UserId            int64    `json:"user_id"`
	Images            []string `json:"images"`
	PaymentMethodsDTO `json:"paymentMethods"`
}
