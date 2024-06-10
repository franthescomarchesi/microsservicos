package createtransaction

import (
	"context"
	"testing"

	"github.com/franthescomarchesi/walletcore/internal/entity"
	"github.com/franthescomarchesi/walletcore/internal/event"
	"github.com/franthescomarchesi/walletcore/internal/usecase/mocks"
	"github.com/franthescomarchesi/walletcore/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransactionUseCase_Execute(t *testing.T) {
	client1, _ := entity.NewClient("John Doe", "j@j.com")
	account1 := entity.NewAccount(client1)
	account1.Credit(1000)
	client2, _ := entity.NewClient("John Doe 2", "j@j2.com")
	account2 := entity.NewAccount(client2)
	account2.Credit(1000)
	mockUow := &mocks.UnitOfWorkMock{}
	mockUow.On("Do", mock.Anything, mock.Anything).Return(nil)
	inputDTO := CreateTransactionInputDTO{
		AccountIDFrom: account1.ID,
		AccountIDTo:   account2.ID,
		Amount:        100,
	}
	dispatcher := events.NewEventDispatcher()
	eventTransaction := event.NewTransactionCreated()
	eventBalance := event.NewBalanceUpdated()
	ctx := context.Background()
	uc := NewCreateTransactionUseCase(mockUow, dispatcher, eventTransaction, eventBalance)
	output, err := uc.Execute(ctx, inputDTO)
	assert.Nil(t, err)
	assert.NotNil(t, output)
	mockUow.AssertExpectations(t)
	mockUow.AssertNumberOfCalls(t, "Do", 1)
}
