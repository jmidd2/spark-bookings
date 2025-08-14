package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"spark-bookings/client"
	"spark-bookings/config"
	"spark-bookings/utils"
	"time"
)

type EventPageData struct {
	PageTitle            string
	ActiveBooking        interface{}
	UpcomingAppointments []client.BookingAppointment
	TotalBookings        int
}

type App struct {
	Config config.Config
	Client *client.AppClient
}

var (
	//go:embed templates/*
	templatesFS embed.FS

	//go:embed assets/*
	assetsFS embed.FS
)

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%sHandling request:%s %s\n", utils.ColorYellow, utils.ColorReset, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func handleError(w http.ResponseWriter, err error) {
	fmt.Println(utils.ColorRed + err.Error() + utils.ColorReset)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%sNot found:%s %s\n", utils.ColorRed, utils.ColorReset, r.URL.Path)
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

var app App

func main() {
	utils.PrintBold(utils.Info("🚀 Starting server"))
	utils.PrintBold(utils.Info("Loading configuration"))

	c, err := config.New()
	if err != nil {
		utils.PrintError(err)
		panic(err)
	}

	app = App{Config: *c}

	app.Client = client.NewGraphClient(app.Config.Bookings)

	http.Handle("/assets/", logMiddleware(http.FileServerFS(assetsFS)))

	http.Handle("/", logMiddleware(http.HandlerFunc(handleIndex)))

	if app.Config.Port == 80 || app.Config.Port == 443 {
		utils.PrintSuccess(fmt.Sprintf("Visit %s in your browser\n", app.Config.Domain))
	} else {
		utils.PrintSuccess(fmt.Sprintf("Visit %s:%d in your browser\n", app.Config.Domain, app.Config.Port))
	}

	utils.PrintBold(utils.Warning("Press Ctrl+C to stop the server"))
	err = http.ListenAndServe(fmt.Sprintf(":%d", app.Config.Port), nil)
	if err != nil {
		utils.PrintError(err)
		panic(err)
	}
}
