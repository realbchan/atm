package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ATM is an interface for interacting with a defined ATM
type ATM interface {
	Login(ctx context.Context, username string, pin int64) (uuid.UUID, ATMError)
	IsAuthenticated(ctx context.Context, token uuid.UUID) ATMError
	GetBalance(ctx context.Context, token uuid.UUID) (float64, ATMError)
	DepositMoney(ctx context.Context, token uuid.UUID, amount float64) (float64, ATMError)
	WithdrawMoney(ctx context.Context, token uuid.UUID, amount float64) (float64, ATMError)
}

// ATMError is an interface for errors that an ATM might return
type ATMError interface {
	IsRetryable() bool
	RetryAfter() time.Duration
	Error() string
	IsAuthenticated() bool
}
