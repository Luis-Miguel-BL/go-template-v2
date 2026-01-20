package vo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPhoneNumber(t *testing.T) {
	tests := []struct {
		phone          string
		wantErrMessage string
	}{
		{phone: "47992783455", wantErrMessage: ""},
		{phone: "11991123555", wantErrMessage: ""},
		{phone: "11945981396", wantErrMessage: ""},
		{phone: "11997421912", wantErrMessage: ""},
		{phone: "00000", wantErrMessage: "phone number must be length beetwen 11 and 11"},
		{phone: "1199231221352", wantErrMessage: "phone number must be length beetwen 11 and 11"},
		{phone: "11891123555", wantErrMessage: "phone number must start with 9"},
		{phone: "22199830866", wantErrMessage: "phone number must start with 9"},
		{phone: "1199112355a", wantErrMessage: "invalid phone number"},
		{phone: "119911a1a57", wantErrMessage: "invalid phone number"},
		{phone: "aa9aaaaaaaa", wantErrMessage: "invalid phone number"},
	}
	for _, tt := range tests {
		phone, err := NewPhoneNumber(tt.phone)
		if tt.wantErrMessage != "" {
			assert.ErrorContains(t, err, tt.wantErrMessage)
			continue
		}
		assert.Equal(t, tt.phone, phone.String())
		assert.NoError(t, err)
	}
}
