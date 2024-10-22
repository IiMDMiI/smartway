package app

import (
	"net/http"

	"github.com/IiMDMiI/smartway/internal/dbservice"
	"github.com/IiMDMiI/smartway/internal/handlers"
	er "github.com/IiMDMiI/smartway/internal/repositories/employeesRepository"
)

type App struct {
}

func New() *App {
	return &App{}
}

func (a *App) Run() error {
	db := dbservice.New()
	defer db.Close()

	emploeesRepository := er.New(db)
	handlers.SetUp(emploeesRepository)

	port := ":8080"
	println("Server is running on port", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
	return nil
}
