package main

import (
	"context"
	"database/sql"
	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/franthescomarchesi/walletcore/internal/database"
	"github.com/franthescomarchesi/walletcore/internal/event"
	"github.com/franthescomarchesi/walletcore/internal/event/handler"
	createaccount "github.com/franthescomarchesi/walletcore/internal/usecase/create_account"
	createclient "github.com/franthescomarchesi/walletcore/internal/usecase/create_client"
	createtransaction "github.com/franthescomarchesi/walletcore/internal/usecase/create_transaction"
	"github.com/franthescomarchesi/walletcore/internal/web"
	"github.com/franthescomarchesi/walletcore/internal/web/webserver"
	"github.com/franthescomarchesi/walletcore/pkg/events"
	"github.com/franthescomarchesi/walletcore/pkg/kafka"
	unitofwork "github.com/franthescomarchesi/walletcore/pkg/unit_of_work"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&multiStatements=true", "root", "root", "mysql-walletcore", "3306", "wallet"))
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
		"wallet",
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
		"group.id":          "wallet",
	}
	kafkaProducer := kafka.NewKafkaProducer(&configMap)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	eventDispatcher.Register("BalanceUpdated", handler.NewBalanceUpdatedKafkaHandler(kafkaProducer))
	transactionCreatedEvent := event.NewTransactionCreated()
	balanceUpdatedEvent := event.NewBalanceUpdated()

	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)

	ctx := context.Background()

	unitofwork := unitofwork.NewUow(ctx, db)
	unitofwork.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})
	unitofwork.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	createClientUseCase := createclient.NewCreateClientUseCase(clientDb)
	createAccountUseCase := createaccount.NewCreateAccountUseCase(accountDb, clientDb)
	createTransactionUseCase := createtransaction.NewCreateTransactionUseCase(unitofwork, eventDispatcher, transactionCreatedEvent, balanceUpdatedEvent)

	webServer := webserver.NewWebServer(":8080")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webServer.AddHandler("/clients", clientHandler.CreateClient)
	webServer.AddHandler("/accounts", accountHandler.CreateAccount)
	webServer.AddHandler("/transactions", transactionHandler.CreateTransaction)

	fmt.Println("Server is running")
	webServer.Start()

}
