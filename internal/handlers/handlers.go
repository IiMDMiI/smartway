package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	em "github.com/IiMDMiI/smartway/api/emploeeManagment"
	"github.com/IiMDMiI/smartway/internal/dbservice"
	mw "github.com/IiMDMiI/smartway/internal/middleware"

	"github.com/lib/pq"
)

const (
	UNEXPECTED_EROR = "Unexpected eror"
)

var (
	valid mw.Validator
	auth  mw.Authorizer
)

func init() {
	valid = mw.NewValidator()
	auth = mw.NewAuthorizer()
}

func SetRoutes() {
	prefix := "/api/v1"

	http.HandleFunc("POST "+prefix+"/emploee", CreateEmployee)
	http.HandleFunc("DELETE "+prefix+"/emploee", DeleteEmployee)
	http.HandleFunc("PUT "+prefix+"/emploee", UpdateEmployee)

	http.HandleFunc("GET "+prefix+"/emploees", GetCompanyEmployees)
}

func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	emp, status, err := mw.AuthorizeAndValidate(auth, valid, r)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	db := dbservice.New()
	defer db.Close()

	id, err := db.CreateEmployee(emp)
	if err != nil {

		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23503" {
				http.Error(w, "The company or department doesn't exist", http.StatusBadRequest)
				return
			}
		}
		log.Println(err)
		http.Error(w, UNEXPECTED_EROR, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Id:%d\n", id)
}

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	if err := auth.Authorize(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "ID is required and must be an integer", http.StatusBadRequest)
		return
	}

	db := dbservice.New()
	defer db.Close()
	db.DeleteEmployee(id)

	fmt.Fprintf(w, "Emploee deleted\n")
}

func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	if err := auth.Authorize(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var emp em.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if emp.Id == 0 {
		http.Error(w, "Id is required", http.StatusBadRequest)
		return
	}
	if emp.Phone != "" {
		if err := valid.ValidatePhone(emp.Phone); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	db := dbservice.New()
	defer db.Close()

	err := db.UpdateEmployee(&emp)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23503" {
				http.Error(w, "The company or department doesn't exist", http.StatusBadRequest)
				return
			}
		}
		log.Println(err)
		http.Error(w, UNEXPECTED_EROR, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Emploee updated")
	fmt.Printf("Emploee updated\n")
}

func GetCompanyEmployees(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "CompanyID must be an integer", http.StatusBadRequest)
		return
	}

	db := dbservice.New()
	defer db.Close()
	var emps []em.Employee
	var err2 error

	dep := r.URL.Query().Get("dep")
	if dep == "" {
		emps, err2 = db.GetCompanyEmployees(id)
	} else {
		emps, err2 = db.GetDepartmentEmployees(id, dep)
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
