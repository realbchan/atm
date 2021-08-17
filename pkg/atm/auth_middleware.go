package atm

import (
	"net/http"

	"github.com/atm/pkg/domain"
	"github.com/google/uuid"
)

type atmAuthenticatorMiddleware struct {
	atmClient domain.ATM
	next      http.Handler
}

// ServeHTTP attempts to authenticate a given request, and
// if failure, returns a 401 unauthorized.
func (m *atmAuthenticatorMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("token")
	token, err := uuid.Parse(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	atmErr := m.atmClient.IsAuthenticated(r.Context(), token)
	if atmErr != nil {
		if atmErr.IsAuthenticated() {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	m.next.ServeHTTP(w, r)
}

func NewATMAuthenticatorMiddleware(atmClient domain.ATM) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return &atmAuthenticatorMiddleware{
			atmClient: atmClient,
			next:      next,
		}
	}
}
