package handler

import (
	"fmt"
	"sync"

	"github.com/franthescomarchesi/walletcore/pkg/events"
	"github.com/franthescomarchesi/walletcore/pkg/kafka"
)

type TransactionCreatedKafkaHandler struct {
	Kafka *kafka.Producer
}

func NewTransactionCreatedKafkaHandler(kafka *kafka.Producer) *TransactionCreatedKafkaHandler {
	return &TransactionCreatedKafkaHandler{
		Kafka: kafka,
	}
}

func (t *TransactionCreatedKafkaHandler) Handle(message events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	t.Kafka.Publish(message, nil, "transactions")
	fmt.Println("TransactionCreatedKafkaHandler called")
}
