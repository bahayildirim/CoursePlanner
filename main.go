package main

import (
	"log"
	"net/http"
	"os"
	"program/router"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	File   *os.File
	Name   string
}

var err error

func (a *App) Run(addr string) error {
	return http.ListenAndServe(addr, a.Router)
}

func main() {
	rt := mux.NewRouter()
	router.RegisterRoutes(rt)
	if err != nil {
		log.Fatal(err)
	}

	files, _ := os.Open("courses.json")
	defer files.Close()

	app := &App{
		Router: rt,
		File:   files,
	}

	err = app.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
