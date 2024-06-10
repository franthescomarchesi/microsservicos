package getaccount

import (
	"testing"

	"github.com/franthescomarchesi/balance/internal/entity"
	"github.com/franthescomarchesi/balance/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
)

func TestSaveAccountUseCase_Execute(t *testing.T) {
	account := entity.NewAccount("123", 25.0)
	accountMock := &mocks.AccountGatewayMock{}
	accountMock.On("FindByID", account.ID).Return(account, nil)
	uc := NewGetAccountUseCase(accountMock)
	inputDTO := GetAccountInputDTO{
		ID: account.ID,
	}
	output, err := uc.Execute(inputDTO)
	assert.Nil(t, err)
	assert.Equal(t, output.ID, account.ID)
	assert.Equal(t, output.Balance, account.Balance)
	accountMock.AssertExpectations(t)
	accountMock.AssertNumberOfCalls(t, "FindByID", 1)
	accountMock.AssertNumberOfCalls(t, "Save", 0)
}
