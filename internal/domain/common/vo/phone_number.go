package vo

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain"
)

type PhoneNumber struct {
	phone string
}

func NewPhoneNumber(phone string) (phoneNumber PhoneNumber, err error) {
	minSize := 11
	maxSize := 11
	re := regexp.MustCompile(`^\d{2}9\d{8}$`)
	if len(phone) < minSize || len(phone) > maxSize {
		return phoneNumber, domain.InvalidInputError("phone_number", fmt.Sprintf("phone number must be length beetwen %d and %d", minSize, maxSize))
	}
	if phone[2:3] != "9" {
		return phoneNumber, domain.InvalidInputError("phone_number", "phone number must start with 9")
	}
	if !re.MatchString(phone) {
		return phoneNumber, domain.InvalidInputError("phone_number", "invalid phone number")
	}

	return PhoneNumber{
		phone: strings.TrimSpace(phone),
	}, nil
}

func (e *PhoneNumber) String() string {
	return e.phone
}
