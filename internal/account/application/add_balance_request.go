package application

type AddBalanceRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}
