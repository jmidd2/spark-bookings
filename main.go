package main

import (
	"html/template"
	"net/http"
)

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func main() {
	tmpl := template.Must(template.ParseFiles("index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := TodoPageData{PageTitle: "Todos", Todos: []Todo{
			{Title: "Buy milk", Done: false},
			{Title: "Buy eggs", Done: true},
		}}
		if err := tmpl.Execute(w, data); err != nil {
			panic(err)
		}
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
