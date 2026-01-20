package vo

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockDomainError is used to compare error messages in tests.
func mockDomainError(msg string) error {
	return errors.New(msg)
}

func TestNewPersonName_ValidNames(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"john doe", "John Doe"},
		{"  alice smith  ", "Alice Smith"},
		{"MARIA GARCIA", "Maria Garcia"},
		{"joHN doE", "John Doe"},
		{"Anna Marie", "Anna Marie"},
	}

	for _, tt := range tests {
		got, err := NewPersonName(tt.input)
		assert.NoError(t, err)
		assert.Equal(t, tt.expected, got.String())
	}
}

func TestNewPersonName_InvalidNames(t *testing.T) {
	tests := []struct {
		input       string
		expectedErr string
	}{
		{"", "name must have at least two names"},
		{"John", "name must have at least two names"},
		{"J Doe", "name must have at least two characters in the first name"},
		{"John D", "name must have at least two characters in the last name"},
		{" A B ", "name must have at least two characters in the first name"},
	}

	for _, tt := range tests {
		_, err := NewPersonName(tt.input)
		assert.Error(t, err)
	}
}

func TestCapitalizeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"john doe", "John Doe"},
		{"ALICE SMITH", "Alice Smith"},
		{"", ""},
		{"a b", "a b"},
		{"joHN doE", "John Doe"},
		{"  maria   garcia  ", "Maria Garcia"},
		{"álvaro silva", "Álvaro Silva"},
	}

	for _, tt := range tests {
		got := capitalizeName(tt.input)
		if got != tt.expected {
			t.Errorf("capitalizeName(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
