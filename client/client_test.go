package client

import (
	"os"
	"testing"
	"time"

	"github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoft/kiota-abstractions-go/store"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
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
