package endpoints

import (
	"context"

	"github.com/1gkx/finstar/internal/repositories/models"
	"github.com/1gkx/finstar/internal/services"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetAll          endpoint.Endpoint
	IncreaseBalance endpoint.Endpoint
	TransferMoney   endpoint.Endpoint
}

func MakeEndpoints(s services.Service) Endpoints {
	return Endpoints{
		GetAll:          makeGetAllEndpoint(s),
		IncreaseBalance: makeIncreaseBalanceEndpoint(s),
		TransferMoney:   makeTransferMoneyEndpoint(s),
	}
}

func makeGetAllEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return s.GetAll(ctx)
	}
}

func makeIncreaseBalanceEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(models.IncreaseBanaceRequest)
		return s.IncreaseBalance(ctx, req)
	}
}

func makeTransferMoneyEndpoint(s services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(models.TransferMoneyRequest)
		return nil, s.TransferMoney(ctx, req)
	}
}
