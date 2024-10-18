package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	em "github.com/IiMDMiI/smartway/api/emploeeManagment"
)

const (
	UnfilledId = -1
)

type Validator interface {
	Validate(r *http.Request) (*em.Employee, error)
	MandatoryFieldsPresent(employee *em.Employee) error
	ValidatePhone(phone string) error
}

func NewValidator() Validator {
	return &JsonValidator{}
}

type JsonValidator struct {
}

func (jv *JsonValidator) Validate(r *http.Request) (*em.Employee, error) {
	var emp em.Employee
	emp.CompanyId = UnfilledId
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		return nil, err
	}

	errChan := make(chan error)
	defer close(errChan)
	chanClosed := false
	go func() {
		err := jv.MandatoryFieldsPresent(&emp)
		if !chanClosed {
			errChan <- err
		}
	}()

	go func() {
		err := jv.ValidatePhone(emp.Phone)
		if !chanClosed {
			errChan <- err
		}
	}()

	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			chanClosed = true
			return nil, err
		}
	}
	return &emp, nil
}

func (jv *JsonValidator) MandatoryFieldsPresent(emp *em.Employee) error {
	//TODO: ask team about any additional restrictions on the fields
	if emp.Name == "" {
		return errors.New("name is required")
	}
	if emp.Surname == "" {
		return errors.New("surname is required")
	}
	if emp.Phone == "" {
		return errors.New("phone is required")
	}
	if emp.CompanyId == UnfilledId {
		return errors.New("companyId is required")
	}
	if emp.Passport.Type == "" {
		return errors.New("passport type is required")
	}
	if emp.Passport.Number == "" {
		return errors.New("passport number is required")
	}
	if emp.Department.Name == "" {
		return errors.New("department name is required")
	}

	return nil
}

func (jv *JsonValidator) ValidatePhone(phone string) error {
	re := regexp.MustCompile(`\D`)
	re.ReplaceAllString(phone, "")

	//TODO: ask team about the phone format
	re2 := regexp.MustCompile(`^\+[1-9]\d{5,14}$`)
	if !re2.MatchString(phone) {
		return errors.New("incorrect phone format")
	}
	return nil
}
