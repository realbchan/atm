package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/atm/pkg/atm"
	v1 "github.com/atm/pkg/handlers/v1"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// main establishes a noop atm client, and attaches it to the http handler for atm requests
// there are a lot of things that could be added, that could be environment-configurations driven, like:
// - what port do we want to listen to?
// - what level of logs should we emit?
// - where/how do stats get emitted?
func main() {
	noopATM := atm.NoOpATMClient{
		LoginAccounts: map[string]int64{
			"bchan": 1234,
		},
		AccountBalances: map[uuid.UUID]float64{},
	}
	atmHandler := v1.ATMHandler{
		ATMClient: &noopATM,
	}
	router := chi.NewRouter()
	router.Post("/login", atmHandler.HandleLogin)

	authMiddleware := atm.NewATMAuthenticatorMiddleware(&noopATM)

	// wrap the other endpoints with an auth middleware, keep the main business logic separate from checking auth
	authGroup := router.Group(nil)
	authGroup.Use(authMiddleware)
	authGroup.Get("/balance", atmHandler.HandleViewBalance)
	authGroup.Post("/balance/deposit", atmHandler.HandleDepositAmount)
	authGroup.Post("/balance/withdraw", atmHandler.HandleWithdrawAmount)

	// expose, listen and serve http
	endpoint := fmt.Sprintf("0.0.0.0:8080")

	if err := http.ListenAndServe(endpoint, router); err != nil {
		fmt.Printf("error running server: %s\n", err.Error())
		os.Exit(1)
	}
}
