package entity

type Account struct {
	ID      string
	Balance float64
}

func NewAccount(id string, balance float64) *Account {
	return &Account{
		ID:      id,
		Balance: balance,
	}
}

func (a *Account) UpdateBalance(balance float64) {
	a.Balance = balance
}
