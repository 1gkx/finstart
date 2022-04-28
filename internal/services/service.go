package services

import (
	"context"
	"fmt"

	"github.com/1gkx/finstar/internal/repositories"
	"github.com/1gkx/finstar/internal/repositories/models"
	"github.com/go-kit/kit/log"
)

type Service interface {
	GetAll(ctx context.Context) ([]models.Account, error)
	IncreaseBalance(context.Context, models.IncreaseBanaceRequest) (*models.Account, error)
	TransferMoney(context.Context, models.TransferMoneyRequest) error
}

type service struct {
	log  log.Logger
	repo repositories.Repository
}

func New(r repositories.Repository, log log.Logger) Service {
	return &service{
		log:  log,
		repo: r,
	}
}

func (s *service) GetAll(ctx context.Context) ([]models.Account, error) {
	return s.repo.GetAll(ctx)
}

func (s *service) IncreaseBalance(ctx context.Context, req models.IncreaseBanaceRequest) (*models.Account, error) {

	s.log.Log("event", "increase_balance", "user_id", req.GetUserId())

	if req.GetAmount() <= 0 {
		s.log.Log("event", "error", "desc", "amount cannot be negative or zero")
		return nil, fmt.Errorf("amount cannot be negative or zero")
	}

	userAccount, err := s.repo.FindAccount(ctx, req.GetUserId())
	if err != nil {
		s.log.Log("event", "error", "desc", err)
		return nil, err
	}

	return s.repo.IncreaseBalance(ctx, userAccount.GetId(), req.GetAmount())
}

func (s *service) TransferMoney(ctx context.Context, req models.TransferMoneyRequest) error {

	s.log.Log("event", "transfer_money", "sender_id", req.GetSenderId())

	if req.GetAmount() <= 0 {
		s.log.Log("event", "error", "desc", "amount cannot be negative or zero")
		return fmt.Errorf("amount cannot be negative or zero")
	}

	if req.GetSenderId() == req.GetReseiverId() {
		s.log.Log("event", "error", "desc", "transfer money yourself denied")
		return fmt.Errorf("transfer money yourself denied")
	}

	senderAcc, err := s.repo.FindAccount(ctx, req.GetSenderId())
	if err != nil {
		s.log.Log("event", "error", "desc", err)
		return err
	}

	s.log.Log("event", "transfer_money", "reseiver_id", req.GetReseiverId())

	receiverAcc, err := s.repo.FindAccount(ctx, req.GetReseiverId())
	if err != nil {
		s.log.Log("event", "error", "desc", err)
		return err
	}

	return s.repo.TransferMoney(
		ctx,
		senderAcc.GetId(),
		receiverAcc.GetId(),
		req.GetAmount(),
	)
}
