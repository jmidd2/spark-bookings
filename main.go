package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"spark-bookings/client"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type EventPageData struct {
	PageTitle            string
	ActiveBooking        interface{}
	UpcomingAppointments []client.BookingAppointment
	TotalBookings        int
}

type Config struct {
	Port     int
	Domain   string
	Secure   bool
	Env      string
	Bookings client.Config
}

type App struct {
	Config Config
	Client *client.AppClient
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

	nowUtc := time.Now().UTC()
	requestStartTime := nowUtc.Format("2006-01-02T15:04:05Z")
	requestEndTime := nowUtc.Add(time.Hour * 24 * 30).Format("2006-01-02T15:04:05Z")

	calendarView, err := app.Client.GetBookingCalendarView(context.Background(), app.Config.Bookings.BusinessID, requestStartTime, requestEndTime)
	if err != nil {
		handleError(w, err)
		return
	}

	data := EventPageData{
		PageTitle:            "Conference Room Events",
		UpcomingAppointments: calendarView.UpcomingAppointments,
		ActiveBooking:        calendarView.ActiveBooking,
		TotalBookings:        len(calendarView.BookingAppointments),
	}

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

func newConfig() (*Config, error) {
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
			return &Config{}, fmt.Errorf("DOMAIN environment variable is required")
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
				return &Config{}, err
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
				return &Config{}, err
			}
		}
	}
	config.Port = port
	config.Domain = url + config.Domain

	//config.Bookings.ClientID = os.Getenv("BOOKINGS_CLIENT_ID")
	config.Bookings.ClientID = GetEnvStringRequired("BOOKINGS_CLIENT_ID")
	//config.Bookings.ClientSecret = os.Getenv("BOOKINGS_CLIENT_SECRET")
	config.Bookings.ClientSecret = GetEnvStringRequired("BOOKINGS_CLIENT_SECRET")
	//config.Bookings.TenantID = os.Getenv("BOOKINGS_TENANT_ID")
	config.Bookings.TenantID = GetEnvStringRequired("BOOKINGS_TENANT_ID")
	//config.Bookings.BusinessID = os.Getenv("BOOKINGS_BUSINESS_ID")
	config.Bookings.BusinessID = GetEnvStringRequired("BOOKINGS_BUSINESS_ID")
	//authUrl := os.Getenv("BOOKINGS_AUTH_URL")
	authUrl := GetEnvStringRequired("BOOKINGS_AUTH_URL")

	config.Bookings.AuthURL = authUrl + "/" + config.Bookings.TenantID
	config.Bookings.Scopes = []string{"https://graph.microsoft.com/.default"}

	return &config, nil
}

func GetEnvString(key string) string {
	return os.Getenv(key)
}

func GetEnvInt(key string) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return 0
	}
	return value
}

func GetEnvBool(key string) bool {
	value, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		return false
	}
	return value
}

func GetEnvStringRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Environment variable %s is required", key))
	}
	return value
}

func GetEnvIntRequired(key string) int {
	value := GetEnvInt(key)
	if value == 0 {
		panic(fmt.Sprintf("Environment variable %s is required", key))
	}
	return value
}

func GetEnvBoolRequired(key string) bool {
	value := GetEnvBool(key)
	if !value {
		panic(fmt.Sprintf("Environment variable %s is required", key))
	}
	return value
}

var app App

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Print(colorBlue + "🚀 Starting server...\n" + colorReset)

	config, err := newConfig()
	if err != nil {
		fmt.Println("Error creating config")
		panic(err)
	}

	app = App{Config: *config}

	fmt.Printf("Config: %+v\n", app.Config)

	app.Client = client.NewGraphClient(app.Config.Bookings)

	http.Handle("/assets/", logMiddleware(http.FileServerFS(assetsFS)))

	http.Handle("/", logMiddleware(http.HandlerFunc(handleIndex)))

	if app.Config.Port == 80 || app.Config.Port == 443 {
		fmt.Printf("%sVisit %s in your browser%s\n", colorGreen, app.Config.Domain, colorReset)
	} else {
		fmt.Printf("%sVisit %s:%d in your browser%s\n", colorGreen, app.Config.Domain, app.Config.Port, colorReset)
	}

	fmt.Println(colorBold + colorYellow + "Press Ctrl+C to stop the server" + colorReset)
	err = http.ListenAndServe(fmt.Sprintf(":%d", app.Config.Port), nil)
	if err != nil {
		fmt.Println(colorRed + err.Error() + colorReset)
		panic(err)
	}
}
