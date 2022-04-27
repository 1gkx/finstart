package models

type Account struct {
	Id      string
	UserId  string
	Balance int
}

func (a *Account) GetId() string {
	return a.Id
}

func (a *Account) GetUserId() string {
	return a.UserId
}

func (a *Account) GetBalance() int {
	return a.Balance
}
