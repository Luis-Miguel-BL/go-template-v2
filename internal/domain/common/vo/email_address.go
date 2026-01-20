package vo

import (
	"regexp"
	"strings"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain"
)

type EmailAddress struct {
	email string
}

var emailPattern = regexp.MustCompile(`^[a-zA-Z0-9]+([._+-][a-zA-Z0-9]+)*[-_a-zA-Z0-9]*@[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*(\.[a-zA-Z]{2,}(-[a-zA-Z]+)*)+$`)

func NewEmailAddress(email string) (emailAddress EmailAddress, err error) {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)
	if !emailPattern.MatchString(email) {
		return emailAddress, domain.InvalidInputError("email_address", "invalid email address")
	}
	return EmailAddress{
		email: email,
	}, nil
}

func (e EmailAddress) String() string {
	return e.email
}
