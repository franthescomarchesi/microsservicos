package createaccount

import (
	"testing"

	"github.com/franthescomarchesi/walletcore/internal/entity"
	"github.com/franthescomarchesi/walletcore/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAccountUseCase_Execute(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "j@j.com")
	clientMock := &mocks.ClientGatewayMock{}
	clientMock.On("Get", client.ID).Return(client, nil)
	accountMock := &mocks.AccountGatewayMock{}
	accountMock.On("Save", mock.Anything).Return(nil)
	uc := NewCreateAccountUseCase(accountMock, clientMock)
	inputDTO := CreateAccountInputDTO{
		ClientID: client.ID,
	}
	output, err := uc.Execute(inputDTO)
	assert.Nil(t, err)
	assert.NotNil(t, output)
	clientMock.AssertExpectations(t)
	accountMock.AssertExpectations(t)
	clientMock.AssertNumberOfCalls(t, "Get", 1)
	accountMock.AssertNumberOfCalls(t, "Save", 1)
}
