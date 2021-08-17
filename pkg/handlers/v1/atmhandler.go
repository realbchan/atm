package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/atm/pkg/domain"
	"github.com/google/uuid"
)

type ATMHandler struct {
	ATMClient domain.ATM
}

// HandleLogin performs login given an ATMClient.
// ideally, we would actually want to use something like
// an auth header, such as basic auth (and not login creds in a payload)
func (h *ATMHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var loginPayload LoginRequest
	json.Unmarshal(body, &loginPayload)
	token, atmErr := h.ATMClient.Login(r.Context(), loginPayload.Username, loginPayload.PIN)
	if atmErr != nil {
		if atmErr.IsAuthenticated() {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	response := LoginResponse{
		Token: token,
	}
	b, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(b))
}

// HandleViewBalance performs GetBalance from the atm client
func (h *ATMHandler) HandleViewBalance(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("token")
	token, err := uuid.Parse(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	balance, atmErr := h.ATMClient.GetBalance(r.Context(), token)
	if atmErr != nil {
		if atmErr.IsAuthenticated() {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	response := BalanceResponse{
		Balance: balance,
	}
	b, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(b))
}

// HandleDepositAmount deposits money through ATMClient.DepositMoney
func (h *ATMHandler) HandleDepositAmount(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("token")
	token, err := uuid.Parse(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var depositRequest AmountATMRequest
	json.Unmarshal(body, &depositRequest)
	balance, atmErr := h.ATMClient.DepositMoney(r.Context(), token, depositRequest.Amount)
	if atmErr != nil {
		if atmErr.IsAuthenticated() {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	response := BalanceResponse{
		Balance: balance,
	}
	b, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(b))
}

// HandleWithdrawAmount withdraws money through ATMClient.DepositMoney
func (h *ATMHandler) HandleWithdrawAmount(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("token")
	token, err := uuid.Parse(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var depositRequest AmountATMRequest
	json.Unmarshal(body, &depositRequest)
	balance, atmErr := h.ATMClient.WithdrawMoney(r.Context(), token, depositRequest.Amount)
	if atmErr != nil {
		if atmErr.IsAuthenticated() {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	response := BalanceResponse{
		Balance: balance,
	}
	b, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(b))
}

type LoginRequest struct {
	Username string `json:"username"`
	PIN      int64  `json:"pin"`
}

type LoginResponse struct {
	Token uuid.UUID `json:"token"`
}

type AmountATMRequest struct {
	Amount float64 `json:"amount"`
}

type BalanceResponse struct {
	Balance float64 `json:"balance"`
}
