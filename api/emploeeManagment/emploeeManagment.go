package employeeManagement

type Employee struct {
	Id         int
	Name       string
	Surname    string
	Phone      string
	CompanyId  int
	Passport   Passport
	Department Department
}

type Passport struct {
	Type   string
	Number string
}

type Department struct {
	Name  string
	Phone string
}
