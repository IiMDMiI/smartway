package middleware

import (
	"errors"
	"net/http"

	em "github.com/IiMDMiI/smartway/api/emploeeManagment"
)

func AuthorizeAndValidate(auth Authorizer, valid Validator, r *http.Request) (*em.Employee, int, error) {
	if err := auth.Authorize(r); err != nil {
		return nil, http.StatusUnauthorized, err
	}

	emp, err := valid.Validate(r)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return emp, http.StatusOK, nil
}

type Authorizer interface {
	Authorize(r *http.Request) error
}

func NewAuthorizer() Authorizer {
	return &TokenAuthorizer{}
}

type TokenAuthorizer struct {
}

func (ta *TokenAuthorizer) Authorize(r *http.Request) error {
	token := r.Header.Get("Authorization")
	return ta.tokenExists(&token)
}

// TODO: get tocken from DB
func (ta *TokenAuthorizer) tokenExists(tocken *string) error {
	if *tocken == "smartway" {
		return nil
	} else {
		return errors.New("Unauthorized")
	}
}
