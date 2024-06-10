package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	account := NewAccount("12345", 25.0)
	assert.Equal(t, account.ID, "12345")
	assert.Equal(t, account.Balance, 25.0)
}

func TestUpdateBalance(t *testing.T) {
	account := NewAccount("12345", 25.0)
	assert.Equal(t, account.ID, "12345")
	assert.Equal(t, account.Balance, 25.0)
	account.UpdateBalance(50.0)
	assert.Equal(t, account.Balance, 50.0)
}
