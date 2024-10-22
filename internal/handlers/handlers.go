package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	em "github.com/IiMDMiI/smartway/api/emploeeManagment"
	mw "github.com/IiMDMiI/smartway/internal/middleware"
	er "github.com/IiMDMiI/smartway/internal/repositories/employeesRepository"
)

const (
	UNEXPECTED_EROR = "Unexpected eror"
)

var empsRep er.EmploeesRepositoryService

func SetUp(newUsersRep er.EmploeesRepositoryService) {
	empsRep = newUsersRep
	setRoutes()
}

func setRoutes() {
	prefix := "/api/v1"

	http.Handle("POST "+prefix+"/emploee", mw.NewAuth(http.HandlerFunc(CreateEmployee)))
	http.Handle("DELETE "+prefix+"/emploee", mw.NewAuth(http.HandlerFunc(DeleteEmployee)))
	http.Handle("PUT "+prefix+"/emploee", mw.NewAuth(http.HandlerFunc(UpdateEmployee)))

	http.HandleFunc("GET "+prefix+"/emploees", GetCompanyEmployees)
}

func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	emp := em.NewEmptyEmploee()

	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := empsRep.Create(emp)
	if err != nil {
		if errors.Is(err, er.ErrBadCompanyIdOrBadDepart) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if errors.Is(err, er.ErrValidation) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprintf(w, "Id:%d\n", id)
}

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "ID is required and must be an integer", http.StatusBadRequest)
		return
	}
	if err2 := empsRep.Delete(id); err2 != nil {
		http.Error(w, err2.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Employee deleted\n")
}

func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	emp := em.NewEmptyEmploee()
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := empsRep.Update(emp); err != nil {
		if errors.Is(err, er.ErrMissingCompanyId) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if errors.Is(err, er.ErrBadPhoneNumberFormat) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, UNEXPECTED_EROR, http.StatusBadRequest)
			log.Println(err)
		}
		return
	}

	fmt.Fprint(w, "Emploee updated")
}

func GetCompanyEmployees(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "CompanyID must be an integer", http.StatusBadRequest)
		return
	}

	var emps []em.Employee
	var err2 error

	dep := r.URL.Query().Get("dep")
	if dep == "" {
		emps, err2 = empsRep.GetCompanyEmployees(id)
	} else {
		emps, err2 = empsRep.GetDepartmentEmployees(id, dep)
	}

	shouldReturn := handleGetEmpsDBErrors(err2, w, emps)
	if shouldReturn {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emps)
}

func handleGetEmpsDBErrors(err error, w http.ResponseWriter, emps []em.Employee) bool {
	if err != nil {
		http.Error(w, UNEXPECTED_EROR, http.StatusInternalServerError)
		log.Println(err)
		return true
	}
	if len(emps) == 0 {
		http.Error(w, "No emploees found", http.StatusNotFound)
		return true
	}
	return false
}
