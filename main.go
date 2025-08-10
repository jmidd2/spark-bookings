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
	"github.com/microsoft/kiota-abstractions-go/store"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"github.com/microsoftgraph/msgraph-sdk-go/solutions"

	"github.com/joho/godotenv"
)

type Event struct {
	Title string
	Done  bool
	Date  time.Time
}

type EventPageData struct {
	PageTitle            string
	ActiveBooking        BookingAppointment
	UpcomingAppointments []BookingAppointment
}

type BookingConfig struct {
	ClientID     string
	ClientSecret string
	TenantID     string
	BusinessID   string
	AuthURL      string
	//RedirectURI  string
	//Scopes       []string
}

type Config struct {
	Port     int
	Domain   string
	Secure   bool
	Env      string
	Bookings BookingConfig
}

type App struct {
	Config Config
	//MsftAcct msalPublic.Client
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

func getStringFromStore(store store.BackingStore, key string) *string {
	val, ok := store.Get(key)
	if ok != nil {
		log.Fatal("Error getting value from store: ", ok)
	}

	return val.(*string)
}

func getQuestionsFromStore(store store.BackingStore) []BookingQuestion {
	questions, err := store.Get("customQuestionAnswers")
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

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/index.html" {
		handleNotFound(w, r)
		return
	}

	cred, err := azidentity.NewClientSecretCredential(app.Config.Bookings.TenantID, app.Config.Bookings.ClientID, app.Config.Bookings.ClientSecret, nil)
	if err != nil {
		log.Fatal("Error creating client credentials: ", err)
	}
	graphClient, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		log.Fatal("Error creating graph client: ", err)
	}

	var top int32 = 10

	requestOptions := solutions.BookingBusinessesItemAppointmentsRequestBuilderGetRequestConfiguration{
		QueryParameters: &solutions.BookingBusinessesItemAppointmentsRequestBuilderGetQueryParameters{
			Top: &top,
		},
	}

	result, err := graphClient.Solutions().BookingBusinesses().ByBookingBusinessId(app.Config.Bookings.BusinessID).Appointments().Get(context.Background(), &requestOptions)
	if err != nil {
		fmt.Println("Error getting user: ", err)
		printOdataError(err)
	}

	appointments := make([]BookingAppointment, 0)

	for _, appointment := range result.GetValue() {
		startTime, err := time.Parse(time.RFC3339Nano, *appointment.GetStartDateTime().GetDateTime())
		if err != nil {
			log.Fatal("Error parsing start time: ", err)
		}

		endTime, err := time.Parse(time.RFC3339Nano, *appointment.GetEndDateTime().GetDateTime())

		//duration := *appointment.GetDuration()

		customers := appointment.GetCustomers()

		if customers == nil || len(customers) == 0 {
			continue
		}

		questions := getQuestionsFromStore(customers[0].GetBackingStore())

		appointments = append(appointments, BookingAppointment{
			Id:                       *appointment.GetId(),
			SelfServiceAppointmentId: getStringValue(appointment.GetSelfServiceAppointmentId()),
			IsLocationOnline:         *appointment.GetIsLocationOnline(),
			CustomerName:             getStringValue(appointment.GetCustomerName()),
			CustomerEmailAddress:     getStringValue(appointment.GetCustomerEmailAddress()),
			CustomerPhone:            getStringValue(appointment.GetCustomerPhone()),
			CustomerNotes:            getStringValue(appointment.GetCustomerNotes()),
			JoinWebUrl:               getStringValue(appointment.GetJoinWebUrl()),
			Customer: BookingCustomer{
				ID:           *getStringFromStore(customers[0].GetBackingStore(), "customerId"),
				EmailAddress: *getStringFromStore(customers[0].GetBackingStore(), "emailAddress"),
				Name:         *getStringFromStore(customers[0].GetBackingStore(), "name"),
				Phone:        *getStringFromStore(customers[0].GetBackingStore(), "phone"),
				Notes:        *getStringFromStore(customers[0].GetBackingStore(), "notes"),
			},
			CustomerTimeZone:        *appointment.GetCustomerTimeZone(),
			SmsNotificationsEnabled: *appointment.GetSmsNotificationsEnabled(),
			ServiceId:               getStringValue(appointment.GetServiceId()),
			ServiceName:             getStringValue(appointment.GetServiceName()),
			Duration:                *appointment.GetDuration(),
			PreBuffer:               *appointment.GetPreBuffer(),
			PostBuffer:              *appointment.GetPostBuffer(),
			PriceType:               *appointment.GetPriceType(),
			Price:                   *appointment.GetPrice(),
			ServiceNotes:            *appointment.GetServiceNotes(),
			OptOutOfCustomerEmail:   *appointment.GetOptOutOfCustomerEmail(),
			AnonymousJoinWebUrl:     *appointment.GetAnonymousJoinWebUrl(),
			StaffMemberIds:          appointment.GetStaffMemberIds(),
			Start: BookingDateTime{
				DateTime: startTime,
				TimeZone: *appointment.GetStartDateTime().GetTimeZone(),
			},
			End: BookingDateTime{
				DateTime: endTime,
				TimeZone: *appointment.GetStartDateTime().GetTimeZone(),
			},
			ServiceLocation:       appointment.GetServiceLocation(),
			Reminders:             appointment.GetReminders(),
			CreatedDateTime:       *appointment.GetCreatedDateTime(),
			LastUpdatedDateTime:   *appointment.GetLastUpdatedDateTime(),
			MaximumAttendeesCount: *appointment.GetMaximumAttendeesCount(),
			Questions:             questions,
		})
	}

	data := EventPageData{PageTitle: "Conference Room Events", UpcomingAppointments: appointments[1:6], ActiveBooking: appointments[0]}

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
	ID           string
	EmailAddress string
	Name         string
	Phone        string
	Notes        string
}

type BookingQuestion struct {
	QuestionID   string
	IsRequired   bool
	QuestionText string
	Answer       string
}

type BookingAppointment struct {
	Id                       string
	SelfServiceAppointmentId string
	IsLocationOnline         bool
	CustomerName             string
	CustomerEmailAddress     string
	CustomerPhone            string
	CustomerNotes            string
	JoinWebUrl               string
	Customer                 BookingCustomer
	CustomerTimeZone         string
	SmsNotificationsEnabled  bool
	ServiceId                string
	ServiceName              string
	Duration                 serialization.ISODuration
	PreBuffer                serialization.ISODuration
	PostBuffer               serialization.ISODuration
	PriceType                models.BookingPriceType
	Price                    float64
	ServiceNotes             string
	OptOutOfCustomerEmail    bool
	AnonymousJoinWebUrl      string
	StaffMemberIds           []string
	Start                    BookingDateTime
	End                      BookingDateTime
	ServiceLocation          models.Locationable
	Reminders                []models.BookingReminderable
	CreatedDateTime          time.Time
	LastUpdatedDateTime      time.Time
	MaximumAttendeesCount    int32
	Questions                []BookingQuestion
}

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
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

	//publicClient, err := msalPublic.New(app.Config.Bookings.ClientID, msalPublic.WithAuthority(app.Config.Bookings.AuthURL))
	//if err != nil {
	//	fmt.Println("Error creating public client")
	//	panic(err)
	//}
	//app.MsftAcct = publicClient
	//
	//fmt.Printf("Public client: %v\n", publicClient)
	//
	//cred, err := confidential.NewCredFromSecret(app.Config.Bookings.ClientSecret)
	//if err != nil {
	//	log.Fatal("Error creating client credentials: ", err)
	//}
	//
	//client, err := confidential.New(app.Config.Bookings.AuthURL, app.Config.Bookings.ClientID, cred)
	//if err != nil {
	//	log.Fatal("Error creating confidential client: ", err)
	//}
	//
	//fmt.Printf("Confidential client: %+v\n", client)
	//
	//scopes := []string{"https://graph.microsoft.com/.default"}
	//ctx := context.Background()
	//result, err := client.AcquireTokenSilent(ctx, scopes)
	//if err != nil {
	//	// Cache miss
	//	result, err = client.AcquireTokenByCredential(ctx, scopes)
	//	if err != nil {
	//		log.Fatal("Error acquiring token: ", err)
	//	}
	//}

	//fmt.Printf("Result: %+v\n", result)
	//fmt.Printf("Result.AccessToken: %+v\n", result.AccessToken)

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
