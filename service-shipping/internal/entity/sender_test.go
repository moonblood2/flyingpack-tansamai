package entity

import (
	"testing"
)

func TestSender_Validation(t *testing.T) {
	s := Sender{
		ContactInfo: ContactInfo{
			Name:        "123456",
			PhoneNumber: "191",
		},
		BirthDate: "1999-05-22",
	}
	t.Logf("s: %+v\n", s)
	err := s.Validate()
	t.Fatal(err)
}
