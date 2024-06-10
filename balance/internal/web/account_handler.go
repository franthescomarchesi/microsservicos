package web

import (
	"encoding/json"
	"net/http"

	getaccount "github.com/franthescomarchesi/balance/internal/usecase/get_account"
	"github.com/go-chi/chi/v5"
)

type WebAccountHandler struct {
	GetAccountUseCase getaccount.GetAccountUseCase
}

func NewWebAccountHandler(getAccountUseCase getaccount.GetAccountUseCase) *WebAccountHandler {
	return &WebAccountHandler{
		GetAccountUseCase: getAccountUseCase,
	}
}

func (h *WebAccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	accountId := chi.URLParam(r, "account_id")
	dto := getaccount.GetAccountInputDTO{
		ID: accountId,
	}
	output, err := h.GetAccountUseCase.Execute(dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusFound)
}
