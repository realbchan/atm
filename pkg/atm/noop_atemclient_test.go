package atm

import (
	context "context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetBalanceSuccess(t *testing.T) {
	token := uuid.New()
	accountBalances := map[uuid.UUID]float64{
		token: 5.1,
	}
	atmClient := NoOpATMClient{
		AccountBalances: accountBalances,
	}

	balance, err := atmClient.GetBalance(context.Background(), token)
	assert.Equal(t, balance, 5.1)
	assert.Nil(t, err)

}

func TestGetBalanceBadToken(t *testing.T) {
	token1 := uuid.New()
	accountBalances := map[uuid.UUID]float64{
		token1: 5.1,
	}
	token2 := uuid.New()
	atmClient := NoOpATMClient{
		AccountBalances: accountBalances,
	}

	_, err := atmClient.GetBalance(context.Background(), token2)
	assert.NotNil(t, err)

}

func TestDepositSuccess(t *testing.T) {
	token := uuid.New()
	accountBalances := map[uuid.UUID]float64{
		token: 5.1,
	}
	atmClient := NoOpATMClient{
		AccountBalances: accountBalances,
	}

	balance, err := atmClient.DepositMoney(context.Background(), token, 1.0)
	assert.Equal(t, balance, 6.1)
	assert.Nil(t, err)
	assert.Equal(t, accountBalances[token], 6.1)

}

func TestDepositBadtoken(t *testing.T) {
	token := uuid.New()
	accountBalances := map[uuid.UUID]float64{
		token: 5.1,
	}
	atmClient := NoOpATMClient{
		AccountBalances: accountBalances,
	}
	token2 := uuid.New()
	_, err := atmClient.DepositMoney(context.Background(), token2, 1.0)
	assert.NotNil(t, err)
	assert.Equal(t, accountBalances[token], 5.1)
}
