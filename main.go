package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"github.com/microsoftgraph/msgraph-sdk-go/solutions"

	"github.com/joho/godotenv"
)

type EventPageData struct {
	PageTitle            string
	ActiveBooking        interface{}
	UpcomingAppointments []BookingAppointment
	TotalBookings        int
}

type BookingBusiness struct {
	ID          string
	DisplayName string
	Email       string
	Phone       string
	WebSiteUrl  string
}

type BookingDateTime struct {
	DateTime time.Time
	TimeZone string
}

type BookingCustomer struct {
	EmailAddress string
	Name         string
	Phone        string
	Notes        string
	TimeZone     string
}

type BookingQuestion struct {
	QuestionID   string
	IsRequired   bool
	QuestionText string
	Answer       string
}

type BookingAppointment struct {
	Id                    string
	Customer              BookingCustomer
	CustomerTimeZone      string
	ServiceId             string
	ServiceName           string
	Duration              serialization.ISODuration
	ServiceNotes          string
	StaffMemberIds        []string
	Start                 BookingDateTime
	End                   BookingDateTime
	CreatedDateTime       time.Time
	LastUpdatedDateTime   time.Time
	Questions             []BookingQuestion
	MaximumAttendeesCount int32
}

type CalendarView struct {
	BookingAppointments  []BookingAppointment
	ActiveBooking        *BookingAppointment
	UpcomingAppointments []BookingAppointment
	BookingBusiness      BookingBusiness
}

type BookingConfig struct {
	ClientID     string
	ClientSecret string
	TenantID     string
	BusinessID   string
	AuthURL      string
	Scopes       []string
}

type Config struct {
	Port     int
	Domain   string
	Secure   bool
	Env      string
	Bookings BookingConfig
}

type AppClient msgraphsdk.GraphServiceClient

type App struct {
	Config Config
	Client *AppClient
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

func getStringFromStore(c models.BookingCustomerInformationBaseable, key string) *string {
	val, ok := c.GetBackingStore().Get(key)
	if ok != nil {
		log.Fatal("Error getting value from store: ", ok)
	}

	return val.(*string)
}

func getQuestionsFromStore(c models.BookingCustomerInformationBaseable) []BookingQuestion {
	questions, err := c.GetBackingStore().Get("customQuestionAnswers")
	if questions == nil || err != nil {
		return nil
	}

	q := make([]BookingQuestion, 0)

	for _, question := range questions.([]models.BookingQuestionAnswerable) {
		q = append(q, BookingQuestion{
			QuestionID:   *question.GetQuestionId(),
			IsRequired:   *question.GetIsRequired(),
			QuestionText: *question.GetQuestion(),
			Answer:       *question.GetAnswer(),
		})
	}

	return q
}

func newBookingCalendarViewRequest(start string, end string) *solutions.BookingBusinessesItemCalendarViewRequestBuilderGetRequestConfiguration {
	qp := &solutions.BookingBusinessesItemCalendarViewRequestBuilderGetQueryParameters{
		Start:  &start,
		End:    &end,
		Select: []string{"id", "customers", "customerName", "customerEmailAddress", "customerPhone", "customerNotes", "customerTimeZone", "serviceId", "serviceName", "duration", "serviceNotes", "staffMemberIds", "startDateTime", "endDateTime", "createdDateTime", "lastUpdatedDateTime", "maximumAttendeesCount"},
	}

	return &solutions.BookingBusinessesItemCalendarViewRequestBuilderGetRequestConfiguration{
		QueryParameters: qp,
	}
}

func (c *AppClient) getBookingCalendarView(ctx context.Context, businessID string, start string, end string) (CalendarView, error) {
	configuration := newBookingCalendarViewRequest(start, end)

	response, err := c.Solutions().BookingBusinesses().ByBookingBusinessId(businessID).CalendarView().Get(ctx, configuration)
	if err != nil {
		fmt.Println("Error getting user: ", err)
		printOdataError(err)
		return CalendarView{}, err
	}

	appointments := make([]BookingAppointment, 0)

	for _, a := range response.GetValue() {

		customers := a.GetCustomers()

		var questions []BookingQuestion

		if customers != nil || len(customers) > 0 {
			questions = getQuestionsFromStore(customers[0])
		}

		startTime, err := getDateTime(a.GetStartDateTime())
		if err != nil {
			return CalendarView{}, err
		}

		endTime, err := getDateTime(a.GetEndDateTime())
		if err != nil {
			return CalendarView{}, err
		}

		customer := getCustomerValue(a)

		appointment := BookingAppointment{
			Id:             *a.GetId(),
			Customer:       customer,
			ServiceId:      getStringValue(a.GetServiceId()),
			ServiceName:    getStringValue(a.GetServiceName()),
			Duration:       getDurationValue(a.GetDuration()),
			ServiceNotes:   getStringValue(a.GetServiceNotes()),
			StaffMemberIds: a.GetStaffMemberIds(),
			Start: BookingDateTime{
				DateTime: startTime,
				TimeZone: getStringValue(a.GetStartDateTime().GetTimeZone()),
			},
			End: BookingDateTime{
				DateTime: endTime,
				TimeZone: getStringValue(a.GetStartDateTime().GetTimeZone()),
			},
			CreatedDateTime:       getTimeValue(a.GetCreatedDateTime()),
			LastUpdatedDateTime:   getTimeValue(a.GetLastUpdatedDateTime()),
			MaximumAttendeesCount: getIntValue(a.GetMaximumAttendeesCount()),
			Questions:             questions,
		}

		appointments = append(appointments, appointment)
	}

	activeBooking, upcomingAppointments := findActiveBooking(appointments)

	return CalendarView{BookingAppointments: appointments, ActiveBooking: activeBooking, UpcomingAppointments: upcomingAppointments}, nil
}

func getDateTime(dt models.DateTimeTimeZoneable) (time.Time, error) {
	dateTime := dt.GetDateTime()
	if dateTime == nil {
		return time.Time{}, fmt.Errorf("DateTime is required")
	}
	return time.Parse(time.RFC3339Nano, *dateTime)
}

func getCustomerValue(a models.BookingAppointmentable) BookingCustomer {
	return BookingCustomer{
		Name:         getStringValue(a.GetCustomerName()),
		EmailAddress: getStringValue(a.GetCustomerEmailAddress()),
		Phone:        getStringValue(a.GetCustomerPhone()),
		Notes:        getStringValue(a.GetCustomerNotes()),
		TimeZone:     getStringValue(a.GetCustomerTimeZone()),
	}
}

func findActiveBooking(appointments []BookingAppointment) (*BookingAppointment, []BookingAppointment) {
	currentTime := time.Now()

	// for testing, so we can see the booking in the future
	if app.Config.Env == "development" {
		currentTime = time.Date(2025, 8, 13, 19, 0, 0, 0, time.UTC)
	}

	for i, a := range appointments {
		if currentTime.Before(a.End.DateTime) && currentTime.After(a.Start.DateTime) {
			return &a, append(appointments[:i], appointments[i+1:]...)
		}
	}
	return nil, appointments
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/index.html" {
		handleNotFound(w, r)
		return
	}

	nowUtc := time.Now().UTC()
	requestStartTime := nowUtc.Format("2006-01-02T15:04:05Z")
	requestEndTime := nowUtc.Add(time.Hour * 24 * 30).Format("2006-01-02T15:04:05Z")

	calendarView, err := app.Client.getBookingCalendarView(context.Background(), app.Config.Bookings.BusinessID, requestStartTime, requestEndTime)
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

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func getTimeValue(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func getIntValue(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}

func getDurationValue(d *serialization.ISODuration) serialization.ISODuration {
	if d == nil {
		return serialization.ISODuration{}
	}
	return *d
}

func newGraphClient() *AppClient {
	cred, err := azidentity.NewClientSecretCredential(app.Config.Bookings.TenantID, app.Config.Bookings.ClientID, app.Config.Bookings.ClientSecret, nil)
	if err != nil {
		log.Fatal("Error creating client credentials: ", err)
	}

	graphClient, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, app.Config.Bookings.Scopes)
	if err != nil {
		log.Fatal("Error creating graph client: ", err)
	}

	return (*AppClient)(graphClient)
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

	app.Client = newGraphClient()

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

func printOdataError(err error) {
	switch err.(type) {
	case *odataerrors.ODataError:
		typed := err.(*odataerrors.ODataError)
		fmt.Printf("error:", typed.Error())
		if terr := typed.GetErrorEscaped(); terr != nil {
			fmt.Printf("code: %s", *terr.GetCode())
			fmt.Printf("msg: %s", *terr.GetMessage())
		}
	default:
		fmt.Printf("%T > error: %#v", err, err)
	}
}
