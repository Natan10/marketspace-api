package dtos

type PaymentMethodsDTO struct {
	Boleto      bool `json:"boleto"`
	Pix         bool `json:"pix"`
	Cash        bool `json:"cash"`
	CreditCard  bool `json:"credit_card"`
	BankDeposit bool `json:"bank_deposit"`
}
