package vo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEmailAddress_ValidEmails(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Test@Example.com", "test@example.com"},
		{" user.name+tag@sub.domain.com ", "user.name+tag@sub.domain.com"},
		{"simple@example.co.uk", "simple@example.co.uk"},
		{"A_B-C.D+E@foo-bar.com", "a_b-c.d+e@foo-bar.com"},
	}

	for _, tt := range tests {
		email, err := NewEmailAddress(tt.input)
		assert.NoError(t, err, "should not error for valid email: %q", tt.input)
		assert.Equal(t, tt.expected, email.String(), "should normalize email")
	}
}

func TestNewEmailAddress_InvalidEmails(t *testing.T) {
	invalidEmails := []string{
		"",
		"plainaddress",
		"@missingusername.com",
		"username@.com",
		"username@com",
		"username@domain..com",
		"username@domain,com",
		"username@domain@domain.com",
		"username@domain.c",
		"username@-domain.com",
		"username@domain-.com",
		"username@.domain.com",
		"username@domain.com.",
		"username@domain..com",
		"username@domain.com-",
	}

	for _, input := range invalidEmails {
		_, err := NewEmailAddress(input)
		assert.Error(t, err, "should error for invalid email: %q", input)
	}
}
