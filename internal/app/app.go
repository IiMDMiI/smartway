package app

import (
	"net/http"

	"github.com/IiMDMiI/smartway/internal/handlers"

	_ "github.com/lib/pq"
)

type App struct {
}

func New() *App {
	return &App{}
}

func (a *App) Run() error {
	handlers.SetRoutes()

	port := ":8080"
	println("Server is running on port", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
	return nil
}
