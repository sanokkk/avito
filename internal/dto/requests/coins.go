package requests

type SendCoinDto struct {
	ToUser string `json:"toUser" validate:"required"`
	Amount int    `json:"amount" validate:"gte=1"`
}
