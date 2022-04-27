package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/1gkx/finstar/internal/repositories/models"
	"github.com/go-kit/kit/log"
	"github.com/jackc/pgx/v4"
)

type Repository interface {
	GetAll(context.Context) ([]models.Account, error)
	FindAccount(context.Context, string) (*models.Account, error)
	IncreaseBalance(context.Context, string, int) (*models.Account, error)
	TransferMoney(context.Context, string, string, int) error
}

type reposytory struct {
	log log.Logger
	c   *pgx.Conn
}

func New(ctx context.Context, dsn string, l log.Logger) (Repository, error) {

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return &reposytory{
		log: l,
		c:   conn,
	}, nil
}

func (r *reposytory) GetAll(ctx context.Context) ([]models.Account, error) {

	var count int
	if err := r.c.QueryRow(ctx, `
		SELECT count(*) FROM finstar.accounts;
	`).Scan(&count); err != nil {
		r.log.Log("event", "error", "desc", err)
		return nil, fmt.Errorf("accounts not found")
	}

	rows, err := r.c.Query(ctx, `
		SELECT * FROM finstar.accounts;
	`)
	if err != nil {
		r.log.Log("event", "error", "desc", err)
		return nil, fmt.Errorf("accounts not found")
	}
	defer rows.Close()

	accounts := make([]models.Account, count)
	account := new(models.Account)
	for idx := 0; rows.Next(); idx++ {
		if err := rows.Scan(
			&account.Id,
			&account.UserId,
			&account.Balance,
		); err != nil {
			r.log.Log("event", "error", "desc", err)
			return nil, fmt.Errorf("accounts not found")
		}
		accounts[idx] = *account
	}

	return accounts, nil
}

func (r *reposytory) FindAccount(ctx context.Context, id string) (*models.Account, error) {

	account := new(models.Account)
	if err := r.c.QueryRow(ctx, `
		SELECT * FROM finstar.accounts
		WHERE user_id = $1;`,
		id,
	).Scan(
		&account.Id,
		&account.UserId,
		&account.Balance,
	); err != nil {
		r.log.Log("event", "error", "desc", err)
		return nil, fmt.Errorf("user %s doesn`t have account", id)
	}

	return account, nil
}

func (r *reposytory) IncreaseBalance(ctx context.Context, id string, amount int) (*models.Account, error) {

	r.log.Log("event", "increase_balance", "account_id", id, "amount", amount)

	account := new(models.Account)
	if err := r.c.QueryRow(ctx, `
		UPDATE finstar.accounts
		SET balance = (
			SELECT balance FROM finstar.accounts
			WHERE id = $1) + $2
		WHERE id = $1
		RETURNING *;`,
		id,
		amount,
	).Scan(
		&account.Id,
		&account.UserId,
		&account.Balance,
	); err != nil {
		r.log.Log("event", "error", "desc", err)
		return nil, err
	}

	return account, nil
}

func (r *reposytory) TransferMoney(ctx context.Context, senderAccount, receiverAccount string, amount int) (e error) {

	r.log.Log(
		"event", "increase_balance",
		"sender_account_id", senderAccount,
		"receiver_account_id", receiverAccount,
		"amount", amount,
	)

	var errTransferMoney = fmt.Errorf(
		"transfer money from %s to %s failed",
		senderAccount,
		receiverAccount,
	)

	tx, err := r.c.Begin(ctx)
	if err != nil {
		r.log.Log("event", "error", "desc", err)
		return errTransferMoney
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(pgx.ErrTxClosed, err) {
			e = err
		}
	}()

	if _, err := tx.Exec(ctx, `
		UPDATE finstar.accounts
		SET balance = (
			SELECT balance FROM finstar.accounts
			WHERE id = $1) - $2
		WHERE id = $1;`,
		senderAccount,
		amount,
	); err != nil {
		r.log.Log("event", "error", "desc", err)
		return errTransferMoney
	}

	if _, err := tx.Exec(ctx, `
		UPDATE finstar.accounts
		SET balance = (
			SELECT balance FROM finstar.accounts
			WHERE id = $1) + $2
		WHERE id = $1;`,
		receiverAccount,
		amount,
	); err != nil {
		r.log.Log("event", "error", "desc", err)
		return errTransferMoney
	}

	if err := tx.Commit(ctx); err != nil {
		r.log.Log("event", "error", "desc", err)
		return errTransferMoney
	}

	return e
}
