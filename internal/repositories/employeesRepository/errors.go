package usersRepository

import (
	"fmt"

	"github.com/IiMDMiI/smartway/internal/dbservice"
)

var (
	ErrBadPhoneNumberFormat    = fmt.Errorf("incorrect phone format")
	ErrMissingCompanyId        = fmt.Errorf("companyId is required")
	ErrBadCompanyIdOrBadDepart = &BadCompanyIdOrBadDepartError{
		msg: "the company or department doesn't exist",
		err: dbservice.ErrForeignKeyViolation,
	}
)

type BadCompanyIdOrBadDepartError struct {
	msg string
	err error
}

func (e *BadCompanyIdOrBadDepartError) Error() string {
	return e.msg
}

func (e *BadCompanyIdOrBadDepartError) Unwrap() error {
	return e.err
}
