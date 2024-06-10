package createtransaction

import (
	"context"

	"github.com/franthescomarchesi/walletcore/internal/entity"
	"github.com/franthescomarchesi/walletcore/internal/gateway"
	"github.com/franthescomarchesi/walletcore/pkg/events"
	unitofwork "github.com/franthescomarchesi/walletcore/pkg/unit_of_work"
)

type CreateTransactionInputDTO struct {
	AccountIDFrom string  `json:"account_id_from"`
	AccountIDTo   string  `json:"account_id_to"`
	Amount        float64 `json:"amount"`
}

type CreateTransactionOutputDTO struct {
	ID            string  `json:"id"`
	AccountIDFrom string  `json:"account_id_from"`
	AccountIDTo   string  `json:"account_id_to"`
	Amount        float64 `json:"amount"`
}

type BalanceUpdatedOutputDTO struct {
	AccountIDFrom        string  `json:"account_id_from"`
	AccountIDTo          string  `json:"account_id_to"`
	BalanceAccountIDFrom float64 `json:"balance_account_id_from"`
	BalanceAccountIDTo   float64 `json:"balance_account_id_to"`
}

type CreateTransactionUseCase struct {
	UnitOfWork         unitofwork.UowInterface
	EventDispatcher    events.EventDispatcherInterface
	TransactionCreated events.EventInterface
	BalanceUpdated     events.EventInterface
}

func NewCreateTransactionUseCase(
	unitOfWork unitofwork.UowInterface,
	eventDispatcher events.EventDispatcherInterface,
	transactionCreated events.EventInterface,
	balanceUpdated events.EventInterface,
) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		UnitOfWork:         unitOfWork,
		EventDispatcher:    eventDispatcher,
		TransactionCreated: transactionCreated,
		BalanceUpdated:     balanceUpdated,
	}
}

func (c *CreateTransactionUseCase) Execute(ctx context.Context, input CreateTransactionInputDTO) (*CreateTransactionOutputDTO, error) {
	output := &CreateTransactionOutputDTO{}
	balanceUpdatedOutput := &BalanceUpdatedOutputDTO{}
	err := c.UnitOfWork.Do(ctx, func(_ *unitofwork.Uow) error {
		accountRepository := c.getAccountRepository(ctx)
		transactionRepository := c.getTransactionRepository(ctx)
		accountFrom, err := accountRepository.FindByID(input.AccountIDFrom)
		if err != nil {
			return err
		}
		accountTo, err := accountRepository.FindByID(input.AccountIDTo)
		if err != nil {
			return err
		}
		transaction, err := entity.NewTransaction(accountFrom, accountTo, input.Amount)
		if err != nil {
			return err
		}
		err = accountRepository.UpdateBalance(accountFrom)
		if err != nil {
			return err
		}
		err = accountRepository.UpdateBalance(accountTo)
		if err != nil {
			return err
		}
		err = transactionRepository.Create(transaction)
		if err != nil {
			return err
		}
		output.ID = transaction.ID
		output.AccountIDFrom = input.AccountIDFrom
		output.AccountIDTo = input.AccountIDTo
		output.Amount = input.Amount
		balanceUpdatedOutput.AccountIDFrom = input.AccountIDFrom
		balanceUpdatedOutput.AccountIDTo = input.AccountIDTo
		balanceUpdatedOutput.BalanceAccountIDFrom = accountFrom.Balance
		balanceUpdatedOutput.BalanceAccountIDTo = accountTo.Balance
		return nil
	})
	if err != nil {
		return nil, err
	}
	c.TransactionCreated.SetPayload(output)
	c.EventDispatcher.Dispatch(c.TransactionCreated)
	c.BalanceUpdated.SetPayload(balanceUpdatedOutput)
	c.EventDispatcher.Dispatch(c.BalanceUpdated)
	return output, nil
}

func (c *CreateTransactionUseCase) getAccountRepository(ctx context.Context) gateway.AccountGateway {
	repo, err := c.UnitOfWork.GetRepository(ctx, "AccountDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.AccountGateway)
}

func (c *CreateTransactionUseCase) getTransactionRepository(ctx context.Context) gateway.TransactionGateway {
	repo, err := c.UnitOfWork.GetRepository(ctx, "TransactionDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.TransactionGateway)
}
