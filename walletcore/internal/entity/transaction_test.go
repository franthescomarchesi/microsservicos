package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {
	client1, _ := NewClient("John Doe", "j@j.com")
	account1 := NewAccount(client1)
	client2, _ := NewClient("John Doe 2", "j@j2.com")
	account2 := NewAccount(client2)
	account1.Credit(1000)
	account2.Credit(1000)
	transaction, err := NewTransaction(account1, account2, 300)
	assert.Nil(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, account1.Balance, float64(700))
	assert.Equal(t, account2.Balance, float64(1300))
}

func TestCreateTransactionWithInsuficientBalance(t *testing.T) {
	client1, _ := NewClient("John Doe", "j@j.com")
	account1 := NewAccount(client1)
	client2, _ := NewClient("John Doe 2", "j@j2.com")
	account2 := NewAccount(client2)
	account1.Credit(100)
	account2.Credit(200)
	transaction, err := NewTransaction(account1, account2, 300)
	assert.Error(t, err, "insufficient funds")
	assert.Nil(t, transaction)
	assert.Equal(t, account1.Balance, float64(100))
	assert.Equal(t, account2.Balance, float64(200))
}
