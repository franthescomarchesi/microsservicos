package database

import (
	"database/sql"

	"github.com/franthescomarchesi/balance/internal/entity"
)

type AccountDB struct {
	DB *sql.DB
}

func NewAccountDB(db *sql.DB) *AccountDB {
	return &AccountDB{
		DB: db,
	}
}

func (a *AccountDB) FindByID(id string) (*entity.Account, error) {
	var account entity.Account
	stmt, err := a.DB.Prepare("SELECT id, balance FROM accounts WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	err = row.Scan(
		&account.ID,
		&account.Balance,
	)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (a *AccountDB) Save(account *entity.Account) error {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM accounts WHERE id = ?)"
	err := a.DB.QueryRow(query, account.ID).Scan(&exists)
	if err != nil {
		return err
	}
	var stmt *sql.Stmt
	if exists {
		stmt, err = a.DB.Prepare("UPDATE accounts SET balance = ? WHERE id = ?")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(account.Balance, account.ID)
	} else {
		stmt, err = a.DB.Prepare("INSERT INTO accounts (id, balance) VALUES (?,?)")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(account.ID, account.Balance)
	}
	if err != nil {
		return err
	}
	return nil
}
