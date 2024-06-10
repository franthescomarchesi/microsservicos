package gateway

import "github.com/franthescomarchesi/walletcore/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}
