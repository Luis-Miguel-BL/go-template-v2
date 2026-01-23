package vo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateValidDocumentNumber(t *testing.T) {
	validDocuments := []string{
		"56049184070",
		"73604599075",
	}

	for _, document := range validDocuments {
		documentNumber, err := NewDocumentNumber(document)
		assert.Nil(t, err)
		assert.Equal(t, document, documentNumber.String())
	}
}

func TestCannotCreateAnInvalidDocumentNumber(t *testing.T) {
	invalidDocuments := []string{
		"abcd",
		"26049184071",
		"438.585.896-06",
	}

	for _, document := range invalidDocuments {
		documentNumber, err := NewDocumentNumber(document)
		assert.Contains(t, err.Error(), "invalid document number")
		assert.Empty(t, documentNumber.String())
	}
}
