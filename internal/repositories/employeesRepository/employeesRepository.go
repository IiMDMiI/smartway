package usersRepository

import (
	"errors"
	"log"

	em "github.com/IiMDMiI/smartway/api/emploeeManagment"
	"github.com/IiMDMiI/smartway/internal/dbservice"
)

var valid = NewValidator()

type EmploeesRepositoryService interface {
	Create(emp *em.Employee) (int, error)
	Update(emp *em.Employee) error
	Delete(id int) error
	GetCompanyEmployees(companyId int) ([]em.Employee, error)
	GetDepartmentEmployees(companyId int, departmentId string) ([]em.Employee, error)
}

func New(db dbservice.DBService) EmploeesRepositoryService {
	return &emploeesRepository{db}
}

type emploeesRepository struct {
	db dbservice.DBService
}

func (er *emploeesRepository) Create(emp *em.Employee) (int, error) {
	if err := valid.Validate(emp); err != nil {
		return 0, err
	}
	id, err := er.db.CreateEmployee(emp)

	if errors.Is(err, dbservice.ErrForeignKeyViolation) {
		return id, ErrBadCompanyIdOrBadDepart
	}
	return id, err
}

func (er *emploeesRepository) Update(emp *em.Employee) error {
	if emp.Id == em.UnfilledId {
		return ErrMissingCompanyId
	}
	if emp.Phone != "" {
		if err := valid.ValidatePhone(emp.Phone); err != nil {
			return ErrBadPhoneNumberFormat
		}
	}

	err := er.db.UpdateEmployee(emp)
	if errors.Is(err, dbservice.ErrForeignKeyViolation) {
		return ErrBadCompanyIdOrBadDepart
	}
	return err
}

func (er *emploeesRepository) Delete(id int) error {
	return er.db.DeleteEmployee(id)
}

func (er *emploeesRepository) GetCompanyEmployees(companyId int) ([]em.Employee, error) {
	employees, err := er.db.GetCompanyEmployees(companyId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return employees, nil
}

func (er *emploeesRepository) GetDepartmentEmployees(companyId int, department string) ([]em.Employee, error) {
	employees, err := er.db.GetDepartmentEmployees(companyId, department)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return employees, nil
}
