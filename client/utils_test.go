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
