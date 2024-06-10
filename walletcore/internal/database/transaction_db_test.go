package database

import (
	"database/sql"
	"testing"

	"github.com/franthescomarchesi/walletcore/internal/entity"
	"github.com/stretchr/testify/suite"
)

type TransactionDBTestSuite struct {
	suite.Suite
	db            *sql.DB
	client1       *entity.Client
	client2       *entity.Client
	account1      *entity.Account
	account2      *entity.Account
	TransactionDB *TransactionDB
}

func (s *TransactionDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	db.Exec("CREATE TABLE clients (id VARCHAR(255), name VARCHAR(255), email VARCHAR(255), created_at date, updated_at date)")
	db.Exec("CREATE TABLE accounts (id VARCHAR(255), client_id VARCHAR(255), balance int, created_at date, updated_at date)")
	db.Exec("CREATE TABLE transactions (id VARCHAR(255), account_id_from VARCHAR(255), account_id_to VARCHAR(255), amount int, created_at date)")
	client1, err := entity.NewClient("John Doe", "j@j.com")
	s.Nil(err)
	s.client1 = client1
	client2, err := entity.NewClient("John Doe 2", "j@j2.com")
	s.Nil(err)
	s.client2 = client2
	account1 := entity.NewAccount(client1)
	account1.Balance = 1000
	s.account1 = account1
	account2 := entity.NewAccount(client2)
	account2.Balance = 1000
	s.account2 = account2
	s.TransactionDB = NewTransactionDB(db)
}

func (s *TransactionDBTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE clients")
	s.db.Exec("DROP TABLE accounts")
	s.db.Exec("DROP TABLE transactions")
}

func TestTransactionDBTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionDBTestSuite))
}

func (s *TransactionDBTestSuite) TestCreate() {
	transaction, err := entity.NewTransaction(s.account1, s.account2, 300)
	s.Nil(err)
	err = s.TransactionDB.Create(transaction)
	s.Nil(err)
}
