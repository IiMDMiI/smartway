package dbservice

import (
	"database/sql"
	"fmt"
	"log"

	em "github.com/IiMDMiI/smartway/api/emploeeManagment"

	_ "github.com/lib/pq"
)

func init() {
	//TODO: get from env
	// host := os.Getenv("DB_HOST")
	// port := os.Getenv("DB_PORT")
	// user := os.Getenv("DB_USER")
	// password := os.Getenv("DB_PASSWORD")
	// dbname := os.Getenv("DB_NAME")

	host := "localhost"
	port := "5432"
	user := "griff"
	password := "1111"
	dbname := "smartway"

	psqlInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}

var (
	psqlInfo string
)

type DBInterface interface {
	CreateEmployee(emp *em.Employee) (int, error)
	DeleteEmployee(id int) error
	UpdateEmployee(emp *em.Employee) error
	GetCompanyEmployees(companyId int) ([]em.Employee, error)
	GetDepartmentEmployees(companyId int, departmentId string) ([]em.Employee, error)
	Close()
}

func New() DBInterface {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	return &DB{db}
}

type DB struct {
	db *sql.DB
}

// Close implements DBInterface.
func (db *DB) Close() {
	db.db.Close()
}

func (db *DB) CreateEmployee(emp *em.Employee) (int, error) {
	row := db.db.QueryRow(`INSERT INTO employee (name, surname, phone, companyid, departmentname)
	VALUES ($1, $2, $3, $4, $5) RETURNING id`, emp.Name, emp.Surname, emp.Phone, emp.CompanyId, emp.Department.Name)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return id, err
	}

	_, err = db.db.Exec(`INSERT INTO passport (type, number, employeeid) VALUES($1,$2,$3)`,
		emp.Passport.Type,
		emp.Passport.Number,
		id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (db *DB) DeleteEmployee(id int) error {
	_, err := db.db.Exec(`DELETE FROM employee WHERE id = $1`, id)
	return err
}

func (db *DB) UpdateEmployee(emp *em.Employee) error {
	_, err := db.db.Exec(`UPDATE employee SET
		name = CASE WHEN $1 != '' THEN $1 ELSE name END,
		surname = CASE WHEN $2 != '' THEN $2 ELSE surname END,
		phone = CASE WHEN $3 != '' THEN $3 ELSE phone END,
		companyid = CASE WHEN $4 != 0 THEN $4 ELSE companyid END,
		departmentname = CASE WHEN $5 != '' THEN $5 ELSE departmentname END
		WHERE id = $6`,
		emp.Name, emp.Surname, emp.Phone, emp.CompanyId, emp.Department.Name, emp.Id)
	if err != nil {
		return err
	}

	_, err = db.db.Exec(`UPDATE passport SET
		type = CASE WHEN $1 != '' THEN $1 ELSE type END,
		number = CASE WHEN $2 != '' THEN $2 ELSE number END
		WHERE employeeid = $3`,
		emp.Passport.Type, emp.Passport.Number, emp.Id)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetCompanyEmployees(companyId int) ([]em.Employee, error) {
	rows, err := db.db.Query(`SELECT e.*, d.phone AS dep_phone, p.type, p.number
		FROM employee e
		JOIN department d ON e.departmentname = d.name AND e.companyid = d.companyid
		JOIN passport p ON p.employeeid = e.id
		WHERE e.companyid = $1;`, companyId)

	if err != nil {
		return nil, err
	}

	return db.rowsToEmps(rows)
}

func (db *DB) GetDepartmentEmployees(companyId int, department string) ([]em.Employee, error) {
	rows, err := db.db.Query(`SELECT e.*, d.phone AS dep_phone, p.type, p.number
		FROM employee e
		JOIN department d ON e.departmentname = d.name AND e.companyid = d.companyid
		JOIN passport p ON p.employeeid = e.id
		WHERE e.companyid = $1 and e.departmentname = $2;`, companyId, department)

	if err != nil {
		return nil, err
	}

	return db.rowsToEmps(rows)
}

func (db *DB) rowsToEmps(rows *sql.Rows) ([]em.Employee, error) {
	defer rows.Close()

	var emps []em.Employee
	for rows.Next() {
		var emp em.Employee
		if err := rows.Scan(&emp.Id, &emp.Name, &emp.Surname, &emp.Phone,
			&emp.CompanyId, &emp.Department.Name, &emp.Department.Phone,
			&emp.Passport.Type, &emp.Passport.Number); err != nil {
			return nil, err
		}
		emps = append(emps, emp)
	}
	return emps, nil
}
