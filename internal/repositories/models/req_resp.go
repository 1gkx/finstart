package models

type IncreaseBanaceRequest struct {
	UserID string `json:"user_id"`
	Amount int    `json:"amount"`
}

func (r *IncreaseBanaceRequest) GetUserId() string {
	return r.UserID
}

func (r *IncreaseBanaceRequest) GetAmount() int {
	return r.Amount
}

type TransferMoneyRequest struct {
	SenderID string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   int    `json:"amount"`
}

func (r *TransferMoneyRequest) GetSenderId() string {
	return r.SenderID
}

func (r *TransferMoneyRequest) GetReseiverId() string {
	return r.Receiver
}

func (r *TransferMoneyRequest) GetAmount() int {
	return r.Amount
}
