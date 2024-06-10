package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/franthescomarchesi/balance/internal/database"
	getaccount "github.com/franthescomarchesi/balance/internal/usecase/get_account"
	saveaccount "github.com/franthescomarchesi/balance/internal/usecase/save_account"
	"github.com/franthescomarchesi/balance/internal/web"
	"github.com/franthescomarchesi/balance/internal/web/webserver"
	"github.com/franthescomarchesi/balance/pkg/kafka"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type KafkaMessage struct {
	Name    string          `json:"Name"`
	Payload json.RawMessage `json:"Payload"`
}

type balanceDTO struct {
	AccountIDFrom        string  `json:"account_id_from"`
	AccountIDTo          string  `json:"account_id_to"`
	BalanceAccountIDFrom float64 `json:"balance_account_id_from"`
	BalanceAccountIDTo   float64 `json:"balance_account_id_to"`
}

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&multiStatements=true", "root", "root", "mysql-balance", "3306", "balance"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations/",
		"balance",
		driver,
	)
	m.Force(1)
	if err != nil {
		panic(err)
	}
	err = m.Down()
	if err != nil {
		panic(err)
	}
	err = m.Up()
	if err != nil {
		panic(err)
	}
	fmt.Println("Migrations ran successfully")

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "balance",
	}
	kafkaConsumer := kafka.NewConsumer(&configMap, []string{"balances"})

	msgChan := make(chan *ckafka.Message)
	go func() {
		kafkaConsumer.Consume(msgChan)
	}()

	accountDb := database.NewAccountDB(db)
	saveAccountUseCase := saveaccount.NewSaveAccountUseCase(accountDb)

	go func() {
		for {
			msg := <-msgChan
			fmt.Println(string(msg.Value))
			var kafkaMessage KafkaMessage
			err := json.Unmarshal(msg.Value, &kafkaMessage)
			if err != nil {
				panic(err)
			}
			var payload balanceDTO
			err = json.Unmarshal(kafkaMessage.Payload, &payload)
			if err != nil {
				panic(err)
			}
			inputAccountTo := saveaccount.SaveAccountInputDTO{
				ID:      payload.AccountIDTo,
				Balance: payload.BalanceAccountIDTo,
			}
			_, err = saveAccountUseCase.Execute(inputAccountTo)
			if err != nil {
				panic(err)
			}
			inputAccountFrom := saveaccount.SaveAccountInputDTO{
				ID:      payload.AccountIDFrom,
				Balance: payload.BalanceAccountIDFrom,
			}
			_, err = saveAccountUseCase.Execute(inputAccountFrom)
			if err != nil {
				panic(err)
			}
		}
	}()

	getAccountUseCase := getaccount.NewGetAccountUseCase(accountDb)

	webServer := webserver.NewWebServer(":3003")

	accountHandler := web.NewWebAccountHandler(*getAccountUseCase)

	webServer.AddHandler("/balances/{account_id}", accountHandler.GetAccount)

	fmt.Println("Server is running")
	webServer.Start()
}
