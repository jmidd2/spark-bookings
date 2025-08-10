package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
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

type Config struct {
	Port   int
	Domain string
	Secure bool
	Env    string
}

const (
	colorReset  = "\x1b[0m"
	colorBold   = "\x1b[1m"
	colorRed    = "\x1b[31m"
	colorGreen  = "\x1b[32m"
	colorYellow = "\x1b[33m"
	colorBlue   = "\x1b[34m"
)

var (
	//go:embed templates/*
	templatesFS embed.FS

	//go:embed assets/*
	assetsFS embed.FS
)

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%sHandling request:%s %s\n", colorYellow, colorReset, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func handleError(w http.ResponseWriter, err error) {
	fmt.Println(colorRed + err.Error() + colorReset)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%sNot found:%s %s\n", colorRed, colorReset, r.URL.Path)
	http.NotFound(w, r)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/index.html" {
		handleNotFound(w, r)
		return
	}

	data := EventPageData{PageTitle: "Conference Room Events", Events: []Event{
		{Title: "Some Meeting", Done: false},
		{Title: "All the people", Done: true},
	}}

	tmpl, err := template.ParseFS(templatesFS, "templates/index.html")
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	if err := tmpl.Execute(w, data); err != nil {
		handleError(w, err)
	}

}

func newConfig() (Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	fmt.Printf("%sRunning in %s mode%s\n", colorBlue, env, colorReset)

	config := Config{
		Secure: os.Getenv("HTTP_SECURE") == "true",
		Env:    env,
	}

	domain := os.Getenv("DOMAIN")
	if domain != "" {
		config.Domain = domain
	} else {
		if config.Env == "development" {
			config.Domain = "localhost"
		} else {
			return Config{}, fmt.Errorf("DOMAIN environment variable is required")
		}
	}

	var url string
	var port int
	if config.Secure {
		url = "https://"
		portStr := os.Getenv("PORT")
		port = 443
		if portStr != "" {
			var err error
			port, err = strconv.Atoi(portStr)
			if err != nil {
				return Config{}, err
			}
		}
	} else {
		url = "http://"
		portStr := os.Getenv("PORT")
		if config.Env == "development" {
			port = 8080
		} else {
			port = 80
		}
		if portStr != "" {
			var err error
			port, err = strconv.Atoi(portStr)
			if err != nil {
				return Config{}, err
			}
		}
	}
	config.Port = port
	config.Domain = url + config.Domain

	return config, nil
}

func main() {
	fmt.Print(colorBlue + "🚀 Starting server...\n" + colorReset)

	config, err := newConfig()
	if err != nil {
		fmt.Println("Error creating config")
		panic(err)
	}

	http.Handle("/assets/", logMiddleware(http.FileServerFS(assetsFS)))

	http.Handle("/", logMiddleware(http.HandlerFunc(handleIndex)))

	if config.Port == 80 || config.Port == 443 {
		fmt.Printf("%sVisit %s in your browser%s\n", colorGreen, config.Domain, colorReset)
	} else {
		fmt.Printf("%sVisit %s:%d in your browser%s\n", colorGreen, config.Domain, config.Port, colorReset)
	}

	fmt.Println(colorBold + colorYellow + "Press Ctrl+C to stop the server" + colorReset)
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
	if err != nil {
		fmt.Println(colorRed + err.Error() + colorReset)
		return
	}
}
