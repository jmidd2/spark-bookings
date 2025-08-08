package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type Event struct {
	Title string
	Done  bool
	Date  time.Time
}

type EventPageData struct {
	PageTitle string
	Events    []Event
}

var (
	//go:embed templates/*
	templatesFS embed.FS

	//go:embed assets/*
	assetsFS embed.FS
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServerFS(assetsFS)))

	test, err := template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		panic(err)
	}
	tmpl := template.Must(test, err)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Handling request: %s\n", r.URL.Path)
		data := EventPageData{PageTitle: "Todos", Events: []Event{
			{Title: "Some Meeting", Done: false},
			{Title: "All the people", Done: true},
		}}
		if err := tmpl.Execute(w, data); err != nil {
			println(err.Error())
			panic(err)
		}
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
