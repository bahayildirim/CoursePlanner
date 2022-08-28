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

func (a *App) Run(addr string) error {
	return http.ListenAndServe(addr, a.Router)
}

func main() {
	rt := mux.NewRouter()
	http.Handle("/", rt)
	router.RegisterRoutes(rt)

	files, _ := os.Open("courses.json")
	defer files.Close()

	app := &App{
		Router: rt,
		File:   files,
	}

	err := app.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
