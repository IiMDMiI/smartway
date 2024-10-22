package middleware

import (
	"errors"
	"net/http"
)

type TokenAuthorizer struct {
	handler http.Handler
}

func NewAuth(handlerToWrap http.Handler) *TokenAuthorizer {
	return &TokenAuthorizer{handlerToWrap}
}

func (t *TokenAuthorizer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if err := t.tokenExists(&token); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	t.handler.ServeHTTP(w, r)
}

// TODO: get tocken from DB
func (ta *TokenAuthorizer) tokenExists(tocken *string) error {
	if *tocken == "smartway" {
		return nil
	} else {
		return errors.New("Unauthorized")
	}
}
