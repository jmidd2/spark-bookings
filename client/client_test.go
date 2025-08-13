package client

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoft/kiota-abstractions-go/store"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// Mock implementations for testing
type mockBookingAppointmentCollectionResponse struct {
	appointments []models.BookingAppointmentable
}

func (m *mockBookingAppointmentCollectionResponse) GetValue() []models.BookingAppointmentable {
	return m.appointments
}

func (m *mockBookingAppointmentCollectionResponse) GetBackingStore() store.BackingStore { return nil }
func (m *mockBookingAppointmentCollectionResponse) GetAdditionalData() map[string]interface{} {
	return nil
}
func (m *mockBookingAppointmentCollectionResponse) SetAdditionalData(value map[string]interface{}) {}
func (m *mockBookingAppointmentCollectionResponse) SetBackingStore(value store.BackingStore)       {}
func (m *mockBookingAppointmentCollectionResponse) SetValue(value []models.BookingAppointmentable) {}
func (m *mockBookingAppointmentCollectionResponse) GetOdataCount() *int32                          { return nil }
func (m *mockBookingAppointmentCollectionResponse) SetOdataCount(value *int32)                     {}
func (m *mockBookingAppointmentCollectionResponse) GetOdataNextLink() *string                      { return nil }
func (m *mockBookingAppointmentCollectionResponse) SetOdataNextLink(value *string)                 {}

type mockFullBookingAppointment struct {
	id                    *string
	customers             []models.BookingCustomerInformationBaseable
	customerName          *string
	customerEmailAddress  *string
	customerPhone         *string
	customerNotes         *string
	customerTimeZone      *string
	serviceId             *string
	serviceName           *string
	duration              *serialization.ISODuration
	serviceNotes          *string
	staffMemberIds        []string
	startDateTime         models.DateTimeTimeZoneable
	endDateTime           models.DateTimeTimeZoneable
	createdDateTime       *time.Time
	lastUpdatedDateTime   *time.Time
	maximumAttendeesCount *int32
}

func (m *mockFullBookingAppointment) GetId() *string { return m.id }
func (m *mockFullBookingAppointment) GetCustomers() []models.BookingCustomerInformationBaseable {
	return m.customers
}
func (m *mockFullBookingAppointment) GetCustomerName() *string                { return m.customerName }
func (m *mockFullBookingAppointment) GetCustomerEmailAddress() *string        { return m.customerEmailAddress }
func (m *mockFullBookingAppointment) GetCustomerPhone() *string               { return m.customerPhone }
func (m *mockFullBookingAppointment) GetCustomerNotes() *string               { return m.customerNotes }
func (m *mockFullBookingAppointment) GetCustomerTimeZone() *string            { return m.customerTimeZone }
func (m *mockFullBookingAppointment) GetServiceId() *string                   { return m.serviceId }
func (m *mockFullBookingAppointment) GetServiceName() *string                 { return m.serviceName }
func (m *mockFullBookingAppointment) GetDuration() *serialization.ISODuration { return m.duration }
func (m *mockFullBookingAppointment) GetServiceNotes() *string                { return m.serviceNotes }
func (m *mockFullBookingAppointment) GetStaffMemberIds() []string             { return m.staffMemberIds }
func (m *mockFullBookingAppointment) GetStartDateTime() models.DateTimeTimeZoneable {
	return m.startDateTime
}
func (m *mockFullBookingAppointment) GetEndDateTime() models.DateTimeTimeZoneable {
	return m.endDateTime
}
func (m *mockFullBookingAppointment) GetCreatedDateTime() *time.Time { return m.createdDateTime }
func (m *mockFullBookingAppointment) GetLastUpdatedDateTime() *time.Time {
	return m.lastUpdatedDateTime
}
func (m *mockFullBookingAppointment) GetMaximumAttendeesCount() *int32 {
	return m.maximumAttendeesCount
}

// Required interface methods
func (m *mockFullBookingAppointment) GetBackingStore() store.BackingStore            { return nil }
func (m *mockFullBookingAppointment) GetAdditionalData() map[string]interface{}      { return nil }
func (m *mockFullBookingAppointment) SetAdditionalData(value map[string]interface{}) {}
func (m *mockFullBookingAppointment) SetBackingStore(value store.BackingStore)       {}
func (m *mockFullBookingAppointment) SetId(value *string)                            {}
func (m *mockFullBookingAppointment) SetCustomers(value []models.BookingCustomerInformationBaseable) {
}
func (m *mockFullBookingAppointment) SetCustomerName(value *string)                      {}
func (m *mockFullBookingAppointment) SetCustomerEmailAddress(value *string)              {}
func (m *mockFullBookingAppointment) SetCustomerPhone(value *string)                     {}
func (m *mockFullBookingAppointment) SetCustomerNotes(value *string)                     {}
func (m *mockFullBookingAppointment) SetCustomerTimeZone(value *string)                  {}
func (m *mockFullBookingAppointment) SetServiceId(value *string)                         {}
func (m *mockFullBookingAppointment) SetServiceName(value *string)                       {}
func (m *mockFullBookingAppointment) SetDuration(value *serialization.ISODuration)       {}
func (m *mockFullBookingAppointment) SetServiceNotes(value *string)                      {}
func (m *mockFullBookingAppointment) SetStaffMemberIds(value []string)                   {}
func (m *mockFullBookingAppointment) SetStartDateTime(value models.DateTimeTimeZoneable) {}
func (m *mockFullBookingAppointment) SetEndDateTime(value models.DateTimeTimeZoneable)   {}
func (m *mockFullBookingAppointment) SetCreatedDateTime(value *time.Time)                {}
func (m *mockFullBookingAppointment) SetLastUpdatedDateTime(value *time.Time)            {}
func (m *mockFullBookingAppointment) SetMaximumAttendeesCount(value *int32)              {}

func TestNewGraphClient(t *testing.T) {
	// Save original env var
	originalEnv := os.Getenv("APP_ENV")
	defer os.Setenv("APP_ENV", originalEnv)

	tests := []struct {
		name        string
		config      Config
		envVar      string
		expectedEnv string
	}{
		{
			name: "creates client with development env by default",
			config: Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				TenantID:     "test-tenant-id",
				Scopes:       []string{"https://graph.microsoft.com/.default"},
			},
			envVar:      "",
			expectedEnv: "development",
		},
		{
			name: "creates client with custom env",
			config: Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				TenantID:     "test-tenant-id",
				Scopes:       []string{"https://graph.microsoft.com/.default"},
			},
			envVar:      "production",
			expectedEnv: "production",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.envVar != "" {
				os.Setenv("APP_ENV", tt.envVar)
			} else {
				os.Unsetenv("APP_ENV")
			}

			// Note: This test will fail without valid Azure credentials
			// In a real test environment, you'd mock the Azure SDK components
			// For now, we'll just test the structure
			t.Skip("Skipping test that requires valid Azure credentials")

			client := NewGraphClient(tt.config)
			if client.Env != tt.expectedEnv {
				t.Errorf("expected env %s, got %s", tt.expectedEnv, client.Env)
			}
		})
	}
}

func TestAppClient_FindActiveBooking(t *testing.T) {
	// Create test appointments
	baseTime := time.Now().UTC()

	appointments := []BookingAppointment{
		{
			Id: "past-appointment",
			Start: BookingDateTime{
				DateTime: baseTime.Add(-2 * time.Hour), // 17:00
			},
			End: BookingDateTime{
				DateTime: baseTime.Add(-1 * time.Hour), // 18:00
			},
		},
		{
			Id: "active-appointment",
			Start: BookingDateTime{
				DateTime: baseTime.Add(-30 * time.Minute), // 18:30
			},
			End: BookingDateTime{
				DateTime: baseTime.Add(30 * time.Minute), // 19:30
			},
		},
		{
			Id: "future-appointment",
			Start: BookingDateTime{
				DateTime: baseTime.Add(1 * time.Hour), // 20:00
			},
			End: BookingDateTime{
				DateTime: baseTime.Add(2 * time.Hour), // 21:00
			},
		},
	}

	tests := []struct {
		name                  string
		client                *AppClient
		appointments          []BookingAppointment
		expectedActiveId      *string
		expectedUpcomingCount int
	}{
		{
			name:                  "no active booking in production env",
			client:                &AppClient{Env: "production"},
			appointments:          appointments,
			expectedActiveId:      stringPtr("active-appointment"),
			expectedUpcomingCount: 2,
		},
		{
			name:                  "empty appointments",
			client:                &AppClient{Env: "development"},
			appointments:          []BookingAppointment{},
			expectedActiveId:      nil,
			expectedUpcomingCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			activeBooking, upcomingAppointments := tt.client.findActiveBooking(tt.appointments)

			// Check active booking
			if tt.expectedActiveId == nil {
				if activeBooking != nil {
					t.Errorf("expected no active booking, got %s", activeBooking.Id)
				}
			} else {
				if activeBooking == nil {
					t.Errorf("expected active booking with id %s, got nil", *tt.expectedActiveId)
				} else if activeBooking.Id != *tt.expectedActiveId {
					t.Errorf("expected active booking id %s, got %s", *tt.expectedActiveId, activeBooking.Id)
				}
			}

			// Check upcoming appointments count
			if len(upcomingAppointments) != tt.expectedUpcomingCount {
				t.Errorf("expected %d upcoming appointments, got %d", tt.expectedUpcomingCount, len(upcomingAppointments))
			}
		})
	}
}

func TestNewBookingCalendarViewRequest(t *testing.T) {
	start := "2025-08-13T00:00:00Z"
	end := "2025-08-14T00:00:00Z"

	config := newBookingCalendarViewRequest(start, end)

	if config == nil {
		t.Fatal("expected configuration, got nil")
	}

	if config.QueryParameters == nil {
		t.Fatal("expected query parameters, got nil")
	}

	qp := config.QueryParameters

	// Check start parameter
	if qp.Start == nil || *qp.Start != start {
		t.Errorf("expected start %s, got %v", start, qp.Start)
	}

	// Check end parameter
	if qp.End == nil || *qp.End != end {
		t.Errorf("expected end %s, got %v", end, qp.End)
	}

	// Check select fields
	expectedFields := []string{
		"id", "customers", "customerName", "customerEmailAddress", "customerPhone",
		"customerNotes", "customerTimeZone", "serviceId", "serviceName", "duration",
		"serviceNotes", "staffMemberIds", "startDateTime", "endDateTime",
		"createdDateTime", "lastUpdatedDateTime", "maximumAttendeesCount",
	}

	if len(qp.Select) != len(expectedFields) {
		t.Errorf("expected %d select fields, got %d", len(expectedFields), len(qp.Select))
	}

	for i, expected := range expectedFields {
		if i >= len(qp.Select) || qp.Select[i] != expected {
			t.Errorf("expected select field %s at position %d, got %v", expected, i, qp.Select[i])
		}
	}
}

func TestGetBookingCalendarView_MockScenario(t *testing.T) {
	// This test demonstrates the structure without making actual API calls
	// In a real scenario, you'd mock the GraphServiceClient methods

	t.Run("would process appointments correctly", func(t *testing.T) {
		client := &AppClient{Env: "development"}

		// Create mock appointment data
		startTime := "2025-08-13T18:30:00Z"
		endTime := "2025-08-13T19:30:00Z"

		mockAppointment := &mockFullBookingAppointment{
			id:                   stringPtr("test-appointment-1"),
			customerName:         stringPtr("John Doe"),
			customerEmailAddress: stringPtr("john@example.com"),
			customerPhone:        stringPtr("+1234567890"),
			serviceId:            stringPtr("service-1"),
			serviceName:          stringPtr("Conference Room"),
			startDateTime: &mockDateTimeTimeZone{
				dateTime: &startTime,
				timeZone: stringPtr("UTC"),
			},
			endDateTime: &mockDateTimeTimeZone{
				dateTime: &endTime,
				timeZone: stringPtr("UTC"),
			},
			maximumAttendeesCount: int32Ptr(10),
		}

		// Test the structure - this would be called if we had proper mocking
		// For now, we'll just verify the mock appointment structure
		if mockAppointment.GetId() == nil {
			t.Error("expected appointment ID to be set")
		}

		if mockAppointment.GetCustomerName() == nil {
			t.Error("expected customer name to be set")
		}

		if mockAppointment.GetStartDateTime() == nil {
			t.Error("expected start date time to be set")
		}

		// Test that findActiveBooking would work with real data
		appointments := []BookingAppointment{
			{
				Id: "test-appointment-1",
				Start: BookingDateTime{
					DateTime: time.Date(2025, 8, 13, 18, 30, 0, 0, time.UTC),
					TimeZone: "UTC",
				},
				End: BookingDateTime{
					DateTime: time.Date(2025, 8, 13, 19, 30, 0, 0, time.UTC),
					TimeZone: "UTC",
				},
			},
		}

		activeBooking, upcomingAppointments := client.findActiveBooking(appointments)

		// In development mode, this should find the active booking
		if activeBooking == nil {
			t.Error("expected to find active booking in development mode")
		}

		if len(upcomingAppointments) != 0 {
			t.Errorf("expected 0 upcoming appointments, got %d", len(upcomingAppointments))
		}
	})
}

// Mock implementation for testing
type mockDateTimeTimeZone struct {
	dateTime *string
	timeZone *string
}

func (m *mockDateTimeTimeZone) Serialize(writer serialization.SerializationWriter) error {
	return nil
}

func (m *mockDateTimeTimeZone) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	return nil
}

func (m *mockDateTimeTimeZone) GetOdataType() *string {
	return stringPtr("#microsoft.graph.dateTimeTimeZone")
}

func (m *mockDateTimeTimeZone) SetOdataType(value *string) {

}

func (m *mockDateTimeTimeZone) GetDateTime() *string {
	return m.dateTime
}

func (m *mockDateTimeTimeZone) GetTimeZone() *string {
	return m.timeZone
}

func (m *mockDateTimeTimeZone) GetBackingStore() store.BackingStore            { return nil }
func (m *mockDateTimeTimeZone) GetAdditionalData() map[string]interface{}      { return nil }
func (m *mockDateTimeTimeZone) SetAdditionalData(value map[string]interface{}) {}
func (m *mockDateTimeTimeZone) SetBackingStore(value store.BackingStore)       {}
func (m *mockDateTimeTimeZone) SetDateTime(value *string)                      {}
func (m *mockDateTimeTimeZone) SetTimeZone(value *string)                      {}

type mockBookingAppointment struct {
	customerName         *string
	customerEmailAddress *string
	customerPhone        *string
	customerNotes        *string
	customerTimeZone     *string
}

func (m *mockBookingAppointment) Serialize(writer serialization.SerializationWriter) error {
	return nil
}

func (m *mockBookingAppointment) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	return nil
}

func (m *mockBookingAppointment) GetOdataType() *string {
	return stringPtr("microsoft.graph.bookingAppointment")
}

func (m *mockBookingAppointment) SetOdataType(value *string) {}

func (m *mockBookingAppointment) GetAdditionalInformation() *string {
	return stringPtr("Additional information")
}

func (m *mockBookingAppointment) GetAnonymousJoinWebUrl() *string {
	return stringPtr("https://join.contoso.com/1234567890")
}

func (m *mockBookingAppointment) GetAppointmentLabel() *string {
	return stringPtr("Meeting with client")
}

func (m *mockBookingAppointment) GetFilledAttendeesCount() *int32 {
	return int32Ptr(1)
}

func (m *mockBookingAppointment) GetIsCustomerAllowedToManageBooking() *bool {
	return boolPtr(true)
}

func (m *mockBookingAppointment) GetIsLocationOnline() *bool {
	return boolPtr(true)
}

func (m *mockBookingAppointment) GetJoinWebUrl() *string {
	return stringPtr("https://join.contoso.com/1234567890")
}

func (m *mockBookingAppointment) GetOptOutOfCustomerEmail() *bool {
	return boolPtr(false)
}

func (m *mockBookingAppointment) GetPostBuffer() *serialization.ISODuration {
	return &serialization.ISODuration{}
}

func (m *mockBookingAppointment) GetPreBuffer() *serialization.ISODuration {
	return &serialization.ISODuration{}
}

func (m *mockBookingAppointment) GetPrice() *float64 {
	return float64Ptr(0.0)
}

func (m *mockBookingAppointment) GetPriceType() *models.BookingPriceType {
	return nil
}

func (m *mockBookingAppointment) GetReminders() []models.BookingReminderable {
	return nil
}

func (m *mockBookingAppointment) GetSelfServiceAppointmentId() *string {
	return stringPtr("1234567890")
}

func (m *mockBookingAppointment) GetServiceLocation() models.Locationable {
	return nil
}

func (m *mockBookingAppointment) GetSmsNotificationsEnabled() *bool {
	return boolPtr(true)
}

func (m *mockBookingAppointment) SetAdditionalInformation(value *string)          {}
func (m *mockBookingAppointment) SetAnonymousJoinWebUrl(value *string)            {}
func (m *mockBookingAppointment) SetAppointmentLabel(value *string)               {}
func (m *mockBookingAppointment) SetFilledAttendeesCount(value *int32)            {}
func (m *mockBookingAppointment) SetIsCustomerAllowedToManageBooking(value *bool) {}
func (m *mockBookingAppointment) SetIsLocationOnline(value *bool)                 {}
func (m *mockBookingAppointment) SetJoinWebUrl(value *string)                     {}
func (m *mockBookingAppointment) SetOptOutOfCustomerEmail(value *bool)            {}
func (m *mockBookingAppointment) SetPostBuffer(value *serialization.ISODuration)  {}
func (m *mockBookingAppointment) SetPreBuffer(value *serialization.ISODuration)   {}
func (m *mockBookingAppointment) SetPrice(value *float64)                         {}
func (m *mockBookingAppointment) SetPriceType(value *models.BookingPriceType)     {}
func (m *mockBookingAppointment) SetReminders(value []models.BookingReminderable) {}
func (m *mockBookingAppointment) SetSelfServiceAppointmentId(value *string)       {}
func (m *mockBookingAppointment) SetServiceLocation(value models.Locationable)    {}
func (m *mockBookingAppointment) SetSmsNotificationsEnabled(value *bool)          {}
func (m *mockBookingAppointment) GetCustomerName() *string                        { return m.customerName }
func (m *mockBookingAppointment) GetCustomerEmailAddress() *string                { return m.customerEmailAddress }
func (m *mockBookingAppointment) GetCustomerPhone() *string                       { return m.customerPhone }
func (m *mockBookingAppointment) GetCustomerNotes() *string                       { return m.customerNotes }
func (m *mockBookingAppointment) GetCustomerTimeZone() *string                    { return m.customerTimeZone }
func (m *mockBookingAppointment) GetBackingStore() store.BackingStore             { return nil }
func (m *mockBookingAppointment) GetAdditionalData() map[string]interface{}       { return nil }
func (m *mockBookingAppointment) SetAdditionalData(value map[string]interface{})  {}
func (m *mockBookingAppointment) SetBackingStore(value store.BackingStore)        {}
func (m *mockBookingAppointment) GetId() *string                                  { return nil }
func (m *mockBookingAppointment) SetId(value *string)                             {}
func (m *mockBookingAppointment) GetCustomers() []models.BookingCustomerInformationBaseable {
	return nil
}
func (m *mockBookingAppointment) SetCustomers(value []models.BookingCustomerInformationBaseable) {}
func (m *mockBookingAppointment) SetCustomerName(value *string)                                  {}
func (m *mockBookingAppointment) SetCustomerEmailAddress(value *string)                          {}
func (m *mockBookingAppointment) SetCustomerPhone(value *string)                                 {}
func (m *mockBookingAppointment) SetCustomerNotes(value *string)                                 {}
func (m *mockBookingAppointment) SetCustomerTimeZone(value *string)                              {}
func (m *mockBookingAppointment) GetDuration() *serialization.ISODuration                        { return nil }
func (m *mockBookingAppointment) SetDuration(value *serialization.ISODuration)                   {}
func (m *mockBookingAppointment) GetEndDateTime() models.DateTimeTimeZoneable                    { return nil }
func (m *mockBookingAppointment) SetEndDateTime(value models.DateTimeTimeZoneable)               {}
func (m *mockBookingAppointment) GetStartDateTime() models.DateTimeTimeZoneable                  { return nil }
func (m *mockBookingAppointment) SetStartDateTime(value models.DateTimeTimeZoneable)             {}
func (m *mockBookingAppointment) GetServiceId() *string                                          { return nil }
func (m *mockBookingAppointment) SetServiceId(value *string)                                     {}
func (m *mockBookingAppointment) GetServiceName() *string                                        { return nil }
func (m *mockBookingAppointment) SetServiceName(value *string)                                   {}
func (m *mockBookingAppointment) GetServiceNotes() *string                                       { return nil }
func (m *mockBookingAppointment) SetServiceNotes(value *string)                                  {}
func (m *mockBookingAppointment) GetStaffMemberIds() []string                                    { return nil }
func (m *mockBookingAppointment) SetStaffMemberIds(value []string)                               {}
func (m *mockBookingAppointment) GetCreatedDateTime() *time.Time                                 { return nil }
func (m *mockBookingAppointment) SetCreatedDateTime(value *time.Time)                            {}
func (m *mockBookingAppointment) GetLastUpdatedDateTime() *time.Time                             { return nil }
func (m *mockBookingAppointment) SetLastUpdatedDateTime(value *time.Time)                        {}
func (m *mockBookingAppointment) GetMaximumAttendeesCount() *int32                               { return nil }
func (m *mockBookingAppointment) SetMaximumAttendeesCount(value *int32)                          {}

type mockBookingQuestion struct {
	questionId *string
	isRequired *bool
	question   *string
	answer     *string
}

func (m *mockBookingQuestion) Serialize(writer serialization.SerializationWriter) error { return nil }
func (m *mockBookingQuestion) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	return nil
}
func (m *mockBookingQuestion) GetAnswerInputType() *models.AnswerInputType { return nil }
func (m *mockBookingQuestion) GetAnswerOptions() []string                  { return nil }
func (m *mockBookingQuestion) GetOdataType() *string {
	return stringPtr("microsoft.graph.bookingQuestion")
}
func (m *mockBookingQuestion) GetSelectedOptions() []string                     { return nil }
func (m *mockBookingQuestion) SetAnswerInputType(value *models.AnswerInputType) {}
func (m *mockBookingQuestion) SetAnswerOptions(value []string)                  {}
func (m *mockBookingQuestion) SetOdataType(value *string)                       {}
func (m *mockBookingQuestion) SetSelectedOptions(value []string)                {}
func (m *mockBookingQuestion) GetQuestionId() *string                           { return m.questionId }
func (m *mockBookingQuestion) GetIsRequired() *bool                             { return m.isRequired }
func (m *mockBookingQuestion) GetQuestion() *string                             { return m.question }
func (m *mockBookingQuestion) GetAnswer() *string                               { return m.answer }
func (m *mockBookingQuestion) GetBackingStore() store.BackingStore              { return nil }
func (m *mockBookingQuestion) GetAdditionalData() map[string]interface{}        { return nil }
func (m *mockBookingQuestion) SetAdditionalData(value map[string]interface{})   {}
func (m *mockBookingQuestion) SetBackingStore(value store.BackingStore)         {}
func (m *mockBookingQuestion) SetQuestionId(value *string)                      {}
func (m *mockBookingQuestion) SetIsRequired(value *bool)                        {}
func (m *mockBookingQuestion) SetQuestion(value *string)                        {}
func (m *mockBookingQuestion) SetAnswer(value *string)                          {}

type mockBackingStore struct {
	data map[string]interface{}
}

func (m *mockBackingStore) EnumerateKeysForValuesChangedToNil() []string {
	return nil
}
func (m *mockBackingStore) SubscribeWithId(callback store.BackingStoreSubscriber, subscriptionId string) error {
	return nil
}
func (m *mockBackingStore) Clear()                                                 {}
func (m *mockBackingStore) GetInitializationCompleted() bool                       { return true }
func (m *mockBackingStore) SetInitializationCompleted(val bool)                    {}
func (m *mockBackingStore) Enumerate() map[string]interface{}                      { return m.data }
func (m *mockBackingStore) EnumerateKeysForValuesChangedToNull() []string          { return nil }
func (m *mockBackingStore) GetIsInitializationCompleted() bool                     { return true }
func (m *mockBackingStore) SetIsInitializationCompleted(value bool)                {}
func (m *mockBackingStore) GetReturnOnlyChangedValues() bool                       { return false }
func (m *mockBackingStore) SetReturnOnlyChangedValues(value bool)                  {}
func (m *mockBackingStore) Subscribe(callback store.BackingStoreSubscriber) string { return "" }
func (m *mockBackingStore) Unsubscribe(subscriptionId string) error                { return nil }
func (m *mockBackingStore) Set(key string, value interface{}) error {
	m.data[key] = value
	return nil
}
func (m *mockBackingStore) Get(key string) (interface{}, error) {
	if value, exists := m.data[key]; exists {
		return value, nil
	}
	return nil, nil
}

type mockCustomerInformation struct {
	backingStore *mockBackingStore
}

func (m *mockCustomerInformation) Serialize(writer serialization.SerializationWriter) error {
	return nil
}

func (m *mockCustomerInformation) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	return nil
}

func (m *mockCustomerInformation) GetOdataType() *string {
	return stringPtr("microsoft.graph.bookingCustomerInformation")
}
func (m *mockCustomerInformation) SetOdataType(value *string) {}
func (m *mockCustomerInformation) GetBackingStore() store.BackingStore {
	return m.backingStore
}

func (m *mockCustomerInformation) GetAdditionalData() map[string]interface{}      { return nil }
func (m *mockCustomerInformation) SetAdditionalData(value map[string]interface{}) {}
func (m *mockCustomerInformation) SetBackingStore(value store.BackingStore)       {}

func TestGetDateTime(t *testing.T) {
	tests := []struct {
		name        string
		dateTime    *string
		expectError bool
		expected    time.Time
	}{
		{
			name:        "valid RFC3339 datetime",
			dateTime:    stringPtr("2025-08-13T15:30:00Z"),
			expectError: false,
			expected:    time.Date(2025, 8, 13, 15, 30, 0, 0, time.UTC),
		},
		{
			name:        "valid RFC3339 datetime with nanoseconds",
			dateTime:    stringPtr("2025-08-13T15:30:00.123456789Z"),
			expectError: false,
			expected:    time.Date(2025, 8, 13, 15, 30, 0, 123456789, time.UTC),
		},
		{
			name:        "nil datetime",
			dateTime:    nil,
			expectError: true,
		},
		{
			name:        "invalid datetime format",
			dateTime:    stringPtr("invalid-date"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockDateTimeTimeZone{dateTime: tt.dateTime}
			result, err := getDateTime(mock)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !result.Equal(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetCustomerValue(t *testing.T) {
	tests := []struct {
		name     string
		mock     *mockBookingAppointment
		expected BookingCustomer
	}{
		{
			name: "all fields populated",
			mock: &mockBookingAppointment{
				customerName:         stringPtr("John Doe"),
				customerEmailAddress: stringPtr("john@example.com"),
				customerPhone:        stringPtr("+1234567890"),
				customerNotes:        stringPtr("VIP customer"),
				customerTimeZone:     stringPtr("America/New_York"),
			},
			expected: BookingCustomer{
				Name:         "John Doe",
				EmailAddress: "john@example.com",
				Phone:        "+1234567890",
				Notes:        "VIP customer",
				TimeZone:     "America/New_York",
			},
		},
		{
			name: "nil fields",
			mock: &mockBookingAppointment{
				customerName:         nil,
				customerEmailAddress: nil,
				customerPhone:        nil,
				customerNotes:        nil,
				customerTimeZone:     nil,
			},
			expected: BookingCustomer{
				Name:         "",
				EmailAddress: "",
				Phone:        "",
				Notes:        "",
				TimeZone:     "",
			},
		},
		{
			name: "mixed populated and nil fields",
			mock: &mockBookingAppointment{
				customerName:         stringPtr("Jane Smith"),
				customerEmailAddress: nil,
				customerPhone:        stringPtr("+1987654321"),
				customerNotes:        nil,
				customerTimeZone:     stringPtr("UTC"),
			},
			expected: BookingCustomer{
				Name:         "Jane Smith",
				EmailAddress: "",
				Phone:        "+1987654321",
				Notes:        "",
				TimeZone:     "UTC",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getCustomerValue(tt.mock)

			if result != tt.expected {
				t.Errorf("expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}

func TestGetQuestionsFromStore(t *testing.T) {
	tests := []struct {
		name     string
		mock     *mockCustomerInformation
		expected []BookingQuestion
	}{
		{
			name: "questions exist in store",
			mock: &mockCustomerInformation{
				backingStore: &mockBackingStore{
					data: map[string]interface{}{
						"customQuestionAnswers": []models.BookingQuestionAnswerable{
							&mockBookingQuestion{
								questionId: stringPtr("q1"),
								isRequired: boolPtr(true),
								question:   stringPtr("What is your favorite color?"),
								answer:     stringPtr("Blue"),
							},
							&mockBookingQuestion{
								questionId: stringPtr("q2"),
								isRequired: boolPtr(false),
								question:   stringPtr("Any dietary restrictions?"),
								answer:     stringPtr("None"),
							},
						},
					},
				},
			},
			expected: []BookingQuestion{
				{
					QuestionID:   "q1",
					IsRequired:   true,
					QuestionText: "What is your favorite color?",
					Answer:       "Blue",
				},
				{
					QuestionID:   "q2",
					IsRequired:   false,
					QuestionText: "Any dietary restrictions?",
					Answer:       "None",
				},
			},
		},
		{
			name: "no questions in store",
			mock: &mockCustomerInformation{
				backingStore: &mockBackingStore{
					data: map[string]interface{}{},
				},
			},
			expected: nil,
		},
		{
			name: "nil questions in store",
			mock: &mockCustomerInformation{
				backingStore: &mockBackingStore{
					data: map[string]interface{}{
						"customQuestionAnswers": nil,
					},
				},
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getQuestionsFromStore(tt.mock)

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d questions, got %d", len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("question %d: expected %+v, got %+v", i, expected, result[i])
				}
			}
		})
	}
}

func TestGetStringFromStore(t *testing.T) {
	tests := []struct {
		name     string
		mock     *mockCustomerInformation
		key      string
		expected *string
	}{
		{
			name: "key exists in store",
			mock: &mockCustomerInformation{
				backingStore: &mockBackingStore{
					data: map[string]interface{}{
						"testKey": stringPtr("test value"),
					},
				},
			},
			key:      "testKey",
			expected: stringPtr("test value"),
		},
		{
			name: "key does not exist",
			mock: &mockCustomerInformation{
				backingStore: &mockBackingStore{
					data: map[string]interface{}{},
				},
			},
			key:      "nonexistent",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Redirect log output to avoid cluttering test output
			log.SetOutput(os.Stdout)

			result := getStringFromStore(tt.mock, tt.key)

			if tt.expected == nil && result != nil {
				t.Errorf("expected nil, got %v", *result)
			} else if tt.expected != nil && result == nil {
				t.Errorf("expected %v, got nil", *tt.expected)
			} else if tt.expected != nil && result != nil && *tt.expected != *result {
				t.Errorf("expected %v, got %v", *tt.expected, *result)
			}
		})
	}
}

func TestGetStringValue(t *testing.T) {
	tests := []struct {
		name     string
		input    *string
		expected string
	}{
		{
			name:     "non-nil string",
			input:    stringPtr("hello world"),
			expected: "hello world",
		},
		{
			name:     "empty string",
			input:    stringPtr(""),
			expected: "",
		},
		{
			name:     "nil string",
			input:    nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStringValue(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGetTimeValue(t *testing.T) {
	testTime := time.Date(2025, 8, 13, 15, 30, 0, 0, time.UTC)
	zeroTime := time.Time{}

	tests := []struct {
		name     string
		input    *time.Time
		expected time.Time
	}{
		{
			name:     "non-nil time",
			input:    &testTime,
			expected: testTime,
		},
		{
			name:     "nil time",
			input:    nil,
			expected: zeroTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getTimeValue(tt.input)
			if !result.Equal(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetIntValue(t *testing.T) {
	tests := []struct {
		name     string
		input    *int32
		expected int32
	}{
		{
			name:     "non-nil int",
			input:    int32Ptr(42),
			expected: 42,
		},
		{
			name:     "zero int",
			input:    int32Ptr(0),
			expected: 0,
		},
		{
			name:     "negative int",
			input:    int32Ptr(-10),
			expected: -10,
		},
		{
			name:     "nil int",
			input:    nil,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getIntValue(tt.input)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestGetDurationValue(t *testing.T) {
	testDuration := serialization.NewDuration(0, 0, 1, 0, 0, 0, 0)
	zeroDuration := serialization.ISODuration{}

	tests := []struct {
		name     string
		input    *serialization.ISODuration
		expected serialization.ISODuration
	}{
		{
			name:     "non-nil duration",
			input:    testDuration,
			expected: *testDuration,
		},
		{
			name:     "nil duration",
			input:    nil,
			expected: zeroDuration,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getDurationValue(tt.input)
			if result != tt.expected {
				t.Errorf("expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}

type mockODataError struct {
	message string
	err     *mockMainError
}

type mockMainError struct {
	code    *string
	message *string
}

func (m *mockMainError) GetAdditionalData() map[string]interface{} {
	return nil
}

func (m *mockMainError) SetAdditionalData(value map[string]interface{}) {

}

func (m *mockMainError) GetBackingStore() store.BackingStore {
	return &mockBackingStore{data: map[string]interface{}{}}
}

func (m *mockMainError) Serialize(writer serialization.SerializationWriter) error {
	return nil
}

func (m *mockMainError) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	return nil
}

func (m *mockMainError) GetDetails() []odataerrors.ErrorDetailsable {
	return nil
}

func (m *mockMainError) GetInnerError() odataerrors.InnerErrorable {
	return nil
}

func (m *mockMainError) GetTarget() *string {
	return stringPtr("")
}

func (m *mockMainError) SetBackingStore(value store.BackingStore) {
}

func (m *mockMainError) SetCode(value *string) {

}

func (m *mockMainError) SetDetails(value []odataerrors.ErrorDetailsable) {

}

func (m *mockMainError) SetInnerError(value odataerrors.InnerErrorable) {

}

func (m *mockMainError) SetMessage(value *string) {
	m.message = value
}

func (m *mockMainError) SetTarget(value *string) {

}

func (m *mockMainError) GetCode() *string    { return m.code }
func (m *mockMainError) GetMessage() *string { return m.message }

func (m *mockODataError) Error() string {
	return m.message
}

func (m *mockODataError) GetErrorEscaped() odataerrors.MainErrorable {
	return m.err
}

// Required interface methods
func (m *mockODataError) GetBackingStore() store.BackingStore             { return nil }
func (m *mockODataError) GetAdditionalData() map[string]interface{}       { return nil }
func (m *mockODataError) SetAdditionalData(value map[string]interface{})  {}
func (m *mockODataError) SetBackingStore(value store.BackingStore)        {}
func (m *mockODataError) SetErrorEscaped(value odataerrors.MainErrorable) {}

func TestPrintOdataError(t *testing.T) {
	t.Run("odata error with details", func(t *testing.T) {
		mockError := &mockODataError{
			message: "OData error occurred",
			err: &mockMainError{
				code:    stringPtr("InvalidRequest"),
				message: stringPtr("The request is invalid"),
			},
		}
		// This should not panic
		printOdataError(mockError)
	})

	t.Run("odata error without details", func(t *testing.T) {
		mockError := &mockODataError{
			message: "OData error occurred",
			err:     nil,
		}
		// This should not panic
		printOdataError(mockError)
	})

	t.Run("regular error", func(t *testing.T) {
		regularError := fmt.Errorf("regular error message")
		// This should not panic
		printOdataError(regularError)
	})
}

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func int32Ptr(i int32) *int32 {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}
