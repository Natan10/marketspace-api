package models

type PaymentMethods struct {
	Id             int64 `json:"id,omitempty"`
	Boleto         bool  `json:"boleto"`
	Pix            bool  `json:"pix"`
	Cash           bool  `json:"cash"`
	CreditCard     bool  `json:"credit_card"`
	BankDeposit    bool  `json:"bank_deposit"`
	AnnouncementId int64 `json:"announcement_id,omitempty"`
}
