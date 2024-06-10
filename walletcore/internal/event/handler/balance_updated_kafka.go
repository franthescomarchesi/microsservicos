package handler

import (
	"fmt"
	"sync"

	"github.com/franthescomarchesi/walletcore/pkg/events"
	"github.com/franthescomarchesi/walletcore/pkg/kafka"
)

type BalanceUpdatedKafkaHandler struct {
	Kafka *kafka.Producer
}

func NewBalanceUpdatedKafkaHandler(kafka *kafka.Producer) *BalanceUpdatedKafkaHandler {
	return &BalanceUpdatedKafkaHandler{
		Kafka: kafka,
	}
}

func (b *BalanceUpdatedKafkaHandler) Handle(message events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	b.Kafka.Publish(message, nil, "balances")
	fmt.Println("BalanceUpdatedKafkaHandler called")
}
