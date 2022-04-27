package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/1gkx/finstar/internal/endpoints"
	"github.com/1gkx/finstar/internal/repositories"
	"github.com/1gkx/finstar/internal/repositories/models"
	"github.com/1gkx/finstar/internal/services"
	"github.com/go-kit/kit/log"
	ht "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

type HttpServer struct {
	IncreaseBalance *ht.Server
	TransferMoney   *ht.Server
}

func New(repo repositories.Repository, log log.Logger) *mux.Router {

	e := endpoints.MakeEndpoints(services.New(repo, log))
	r := mux.NewRouter()

	r.Methods("GET").Path("/api/v1/get_all").Handler(ht.NewServer(
		e.GetAll,
		decodeEmptyRequest,
		encodeResponse,
	))
	r.Methods("POST").Path("/api/v1/increase").Handler(ht.NewServer(
		e.IncreaseBalance,
		decodeIncreaseBalanceRequest,
		encodeResponse,
	))
	r.Methods("POST").Path("/api/v1/transfer").Handler(ht.NewServer(
		e.TransferMoney,
		decodeTransferMoneyRequest,
		encodeResponse,
	))

	return r
}

func decodeEmptyRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return nil, nil
}

func decodeIncreaseBalanceRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req models.IncreaseBanaceRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeTransferMoneyRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req models.TransferMoneyRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if response == nil {
		w.WriteHeader(http.StatusAccepted)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	return json.NewEncoder(w).Encode(response)
}
