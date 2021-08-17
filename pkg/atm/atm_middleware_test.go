package atm

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestServeHTTPSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockATM := NewMockATM(ctrl)
	mockHandle := NewMockHandler(ctrl)
	middleWare := NewATMAuthenticatorMiddleware(mockATM)
	responseWrite := NewMockResponseWriter(ctrl)
	handler := middleWare(mockHandle)
	u, _ := url.Parse("http://localhost")
	req, _ := http.NewRequest(http.MethodGet, u.String(), http.NoBody)
	token := uuid.New()
	req.Header.Set("token", token.String())
	mockATM.EXPECT().IsAuthenticated(gomock.Any(), gomock.Any()).Return(nil)
	mockHandle.EXPECT().ServeHTTP(gomock.Any(), gomock.Any())
	handler.ServeHTTP(responseWrite, req)
}

func TestServeHTTPUnAuthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockATM := NewMockATM(ctrl)
	mockHandle := NewMockHandler(ctrl)
	middleWare := NewATMAuthenticatorMiddleware(mockATM)
	responseWrite := NewMockResponseWriter(ctrl)
	handler := middleWare(mockHandle)
	u, _ := url.Parse("http://localhost")
	req, _ := http.NewRequest(http.MethodGet, u.String(), http.NoBody)
	token := uuid.New()
	req.Header.Set("token", token.String())
	mockATM.EXPECT().IsAuthenticated(gomock.Any(), gomock.Any()).Return(noOpATMError{})
	responseWrite.EXPECT().WriteHeader(http.StatusUnauthorized).Return()
	handler.ServeHTTP(responseWrite, req)
}
