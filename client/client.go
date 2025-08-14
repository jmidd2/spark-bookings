package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"github.com/microsoftgraph/msgraph-sdk-go/solutions"
)

type AppClient struct {
	msgraphsdk.GraphServiceClient

	Env string
}

type Config struct {
	ClientID     string   `yaml:"client_id"`
	ClientSecret string   `yaml:"client_secret"`
	TenantID     string   `yaml:"tenant_id"`
	BusinessID   string   `yaml:"business_id"`
	AuthURL      string   `yaml:"auth_url,omitempty"`
	Scopes       []string `yaml:"scopes,omitempty"`
}

type CalendarView struct {
	BookingAppointments  []BookingAppointment
	ActiveBooking        *BookingAppointment
	UpcomingAppointments []BookingAppointment
	BookingBusiness      BookingBusiness
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

func NewGraphClient(config Config) *AppClient {
	cred, err := azidentity.NewClientSecretCredential(config.TenantID, config.ClientID, config.ClientSecret, nil)
	if err != nil {
		log.Fatal("Error creating client credentials: ", err)
	}

	graphClient, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, config.Scopes)
	if err != nil {
		log.Fatal("Error creating graph client: ", err)
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	client := AppClient{
		GraphServiceClient: *graphClient,
		Env:                env,
	}

	return &client
}

func (c *AppClient) GetBookingCalendarView(ctx context.Context, businessID string, start string, end string) (CalendarView, error) {
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

	activeBooking, upcomingAppointments := c.findActiveBooking(appointments)

	return CalendarView{BookingAppointments: appointments, ActiveBooking: activeBooking, UpcomingAppointments: upcomingAppointments}, nil
}

func (c *AppClient) findActiveBooking(appointments []BookingAppointment) (*BookingAppointment, []BookingAppointment) {
	currentTime := time.Now()

	// for testing, so we can see the booking in the future
	if c.Env == "development" {
		currentTime = time.Date(2025, 8, 13, 19, 0, 0, 0, time.UTC)
	}

	for i, a := range appointments {
		if currentTime.Before(a.End.DateTime) && currentTime.After(a.Start.DateTime) {
			return &a, append(appointments[:i], appointments[i+1:]...)
		}
	}
	return nil, appointments
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

func getStringFromStore(c models.BookingCustomerInformationBaseable, key string) *string {
	val, ok := c.GetBackingStore().Get(key)
	if ok != nil {
		log.Fatal("Error getting value from store: ", ok)
	}

	if val == nil {
		return nil
	}

	return val.(*string)
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

func printOdataError(err error) {
	var ODataError *odataerrors.ODataError
	switch {
	case errors.As(err, &ODataError):
		if ODataError == nil {
			return
		}
		fmt.Printf("error: %s", ODataError.Error())
		if terr := ODataError.GetErrorEscaped(); terr != nil {
			fmt.Printf("code: %s", *terr.GetCode())
			fmt.Printf("msg: %s", *terr.GetMessage())
		}
	default:
		fmt.Printf("%T > error: %#v", err, err)
	}
}
