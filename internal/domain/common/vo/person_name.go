package vo

import (
	"strings"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain"
)

type PersonName struct {
	name string
}

func NewPersonName(name string) (personName PersonName, err error) {
	name = strings.TrimSpace(name)
	name = capitalizeName(name)

	names := strings.Split(name, " ")
	if len(names) < 2 {
		return personName, domain.InvalidInputError("person_name", "name must have at least two names")
	}
	if len(names[0]) < 2 {
		return personName, domain.InvalidInputError("person_name", "name must have at least two characters in the first name")
	}
	if len(names[len(names)-1]) < 2 {
		return personName, domain.InvalidInputError("person_name", "name must have at least two characters in the last name")
	}

	return PersonName{
		name: strings.TrimSpace(name),
	}, nil
}

func (e PersonName) String() string {
	return e.name
}

func capitalizeName(name string) string {
	if name == "" {
		return ""
	}
	parts := strings.Fields(name)
	for i, part := range parts {
		runes := []rune(part)
		if len(runes) > 1 {
			parts[i] = strings.ToUpper(string(runes[0])) + strings.ToLower(string(runes[1:]))
		}
	}
	return strings.Join(parts, " ")
}
