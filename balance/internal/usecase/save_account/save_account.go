package saveaccount

import (
	"github.com/franthescomarchesi/balance/internal/entity"
	"github.com/franthescomarchesi/balance/internal/gateway"
)

type SaveAccountInputDTO struct {
	ID      string  `json:"id"`
	Balance float64 `json:"balance"`
}

type SaveAccountOutputDTO struct {
	ID      string  `json:"id"`
	Balance float64 `json:"balance"`
}

type SaveAccountUseCase struct {
	AccountGateway gateway.AccountGateway
}

func NewSaveAccountUseCase(accountGateway gateway.AccountGateway) *SaveAccountUseCase {
	return &SaveAccountUseCase{
		AccountGateway: accountGateway,
	}
}

func (u *SaveAccountUseCase) Execute(input SaveAccountInputDTO) (*SaveAccountOutputDTO, error) {
	account := entity.NewAccount(input.ID, input.Balance)
	err := u.AccountGateway.Save(account)
	if err != nil {
		return nil, err
	}
	return &SaveAccountOutputDTO{
		ID:      account.ID,
		Balance: account.Balance,
	}, nil
}
