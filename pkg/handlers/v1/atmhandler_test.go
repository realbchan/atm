package v1

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestHandleWithdrawAmountSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockATM := NewMockATM(ctrl)
	responseWrite := NewMockResponseWriter(ctrl)
	handler := ATMHandler{
		ATMClient: mockATM,
	}
	u, _ := url.Parse("http://localhost")
	json := `{"amount": 10.0}`
	req, _ := http.NewRequest(http.MethodGet, u.String(), ioutil.NopCloser(strings.NewReader(json)))
	token := uuid.New()
	req.Header.Set("token", token.String())

	mockATM.EXPECT().WithdrawMoney(gomock.Any(), gomock.Any(), 10.0).Return(10.0, nil)
	responseWrite.EXPECT().Header().Return(http.Header{})
	responseWrite.EXPECT().Write(gomock.Any())
	responseWrite.EXPECT().WriteHeader(http.StatusOK)
	handler.HandleWithdrawAmount(responseWrite, req)
}
