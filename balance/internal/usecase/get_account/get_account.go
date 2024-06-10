package getaccount

import (
	"github.com/franthescomarchesi/balance/internal/gateway"
)

type GetAccountInputDTO struct {
	ID string `json:"id"`
}

type GetAccountOutputDTO struct {
	ID      string  `json:"id"`
	Balance float64 `json:"balance"`
}

type GetAccountUseCase struct {
	AccountGateway gateway.AccountGateway
}

func NewGetAccountUseCase(accountGateway gateway.AccountGateway) *GetAccountUseCase {
	return &GetAccountUseCase{
		AccountGateway: accountGateway,
	}
}

func (u *GetAccountUseCase) Execute(input GetAccountInputDTO) (*GetAccountOutputDTO, error) {
	account, err := u.AccountGateway.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	return &GetAccountOutputDTO{
		ID:      account.ID,
		Balance: account.Balance,
	}, nil
}
