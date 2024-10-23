package dbservice

import (
	"fmt"
	"strings"

	em "github.com/IiMDMiI/smartway/api/emploeeManagment"
)

func createEmploeeArgs(emp *em.Employee) map[string]any {
	args := make(map[string]any)

	if emp.Name != "" {
		args["name"] = emp.Name
	}
	if emp.Surname != "" {
		args["surname"] = emp.Surname
	}
	if emp.Phone != "" {
		args["phone"] = emp.Phone
	}
	if emp.CompanyId != em.UnfilledId {
		args["companyid"] = emp.CompanyId
	}
	if emp.Department.Name != "" {
		args["departmentname"] = emp.Department.Name
	}
	return args
}

func createPassportArgs(pass *em.Passport) map[string]any {
	args := make(map[string]any)

	if pass.Type != "" {
		args["type"] = pass.Type
	}
	if pass.Number != "" {
		args["number"] = pass.Number
	}
	return args
}

func createUpdateQuery(table string, args map[string]any, condition string) string {
	var sb strings.Builder
	sb.WriteString("UPDATE " + table + " SET ")

	first := true
	for k, v := range args {
		if !first {
			sb.WriteString(", ")
		}
		first = false

		switch val := v.(type) {
		case string:
			sb.WriteString(k + " = '" + val + "'")
		case nil:
			sb.WriteString(k + " = NULL")
		default:
			sb.WriteString(fmt.Sprintf("%s = %v", k, val))
		}
	}

	sb.WriteString(" WHERE " + condition)
	return sb.String()
}
