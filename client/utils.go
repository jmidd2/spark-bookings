package client

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

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
