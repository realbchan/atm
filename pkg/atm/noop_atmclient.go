package atm

import (
	"context"
	"errors"
	"time"

	"github.com/atm/pkg/domain"
	"github.com/google/uuid"
)

// Go does not support const arrays
// https://qvault.io/golang/golang-constant-maps-slices/
var emptyUUID = [16]byte{}

// NoOpATMClient is a noop implementation of ATMClient
// in reality, you would have a more intelligent ATMClient,
// maybe something that makes http calls
type NoOpATMClient struct {
	LoginAccounts   map[string]int64
	AccountBalances map[uuid.UUID]float64
}

// Login verifies that the provided login information is good, returns a token
// that is saved by AccountBalances
func (atm *NoOpATMClient) Login(ctx context.Context, username string, pin int64) (uuid.UUID, domain.ATMError) {
	if usernamePin, ok := atm.LoginAccounts[username]; ok && pin == usernamePin {
		token := uuid.New()
		atm.AccountBalances[token] = 0.0 // set to 0, this client isn't smart enough to have a persistent data store
		return token, nil
	}

	return emptyUUID, noOpATMError{
		originalError:   errors.New("not authenticated"),
		retry:           false,
		retryAfter:      time.Second * 0,
		isAuthenticated: false,
	}
}

// IsAuthenticated checks if token exists
func (atm *NoOpATMClient) IsAuthenticated(ctx context.Context, token uuid.UUID) domain.ATMError {
	if _, ok := atm.AccountBalances[token]; ok {
		return nil
	}
	return noOpATMError{
		originalError:   errors.New("not authenticated"),
		retry:           false,
		retryAfter:      time.Second * 0,
		isAuthenticated: false,
	}
}

// GetBalance returns balance for a given token that's logged in
func (atm *NoOpATMClient) GetBalance(ctx context.Context, token uuid.UUID) (float64, domain.ATMError) {
	if balance, ok := atm.AccountBalances[token]; ok {
		return balance, nil
	}
	return 0.0, noOpATMError{
		originalError:   errors.New("bad token"),
		retry:           false,
		retryAfter:      time.Second * 0,
		isAuthenticated: false,
	}
}

// DepositMoney deposits money into a given account determined by token
func (atm *NoOpATMClient) DepositMoney(ctx context.Context, token uuid.UUID, amount float64) (float64, domain.ATMError) {
	if balance, ok := atm.AccountBalances[token]; ok {
		newBalance := balance + amount
		atm.AccountBalances[token] = newBalance
		return newBalance, nil
	}
	return 0.0, noOpATMError{
		originalError:   errors.New("bad token"),
		retry:           false,
		retryAfter:      time.Second * 0,
		isAuthenticated: false,
	}
}

// WithdrawMoney withdraws money into a given account determined by token
func (atm *NoOpATMClient) WithdrawMoney(ctx context.Context, token uuid.UUID, amount float64) (float64, domain.ATMError) {
	if balance, ok := atm.AccountBalances[token]; ok {
		newBalance := balance - amount
		atm.AccountBalances[token] = newBalance
		return newBalance, nil
	}
	return 0.0, noOpATMError{
		originalError:   errors.New("bad token"),
		retry:           false,
		retryAfter:      time.Second * 0,
		isAuthenticated: false,
	}
}

// noOpATMError satifies the ATM interface
type noOpATMError struct {
	retry           bool
	retryAfter      time.Duration
	originalError   error
	isAuthenticated bool
}

func (e noOpATMError) IsRetryable() bool {
	return e.retry
}

func (e noOpATMError) RetryAfter() time.Duration {
	return e.retryAfter
}

func (e noOpATMError) Error() string {
	return e.originalError.Error()
}

func (e noOpATMError) IsAuthenticated() bool {
	return e.isAuthenticated
}
