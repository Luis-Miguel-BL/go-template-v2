package api

import (
	"net/http"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/http/request"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/http/response"
	"github.com/Luis-Miguel-BL/go-lm-template/test/util"
)

func (s *APISuite) Test_CreateLeadEndpoint() {

	s.Run("Should return 401 when accessToken is missing", func() {
		validBody := request.CreateLeadRequest{
			Email:          "fake@email.com",
			DocumentNumber: "12345678909",
			Name:           "Fake Name",
			Phone:          "99990000000",
		}
		resp, _ := s.TestUtil.DoRequest(s.T().Context(), util.RequestParams{
			Path:   "/v1/leads",
			Method: http.MethodPost,
			Body:   validBody,
		})

		s.Equal(http.StatusUnauthorized, resp.StatusCode())
	})

	s.Run("Should return 400 when email is invalid", func() {
		invalidBody := request.CreateLeadRequest{
			Email:          "invalid-email",
			DocumentNumber: "12345678909",
			Name:           "Fake Name",
			Phone:          "99990000000",
		}

		resp, err := s.TestUtil.RequestWithDefaultAccessToken(s.T().Context(), util.RequestParams{
			Path:   "/v1/leads",
			Method: http.MethodPost,
			Body:   invalidBody,
		})

		s.NoError(err)
		s.Equal(http.StatusBadRequest, resp.StatusCode())

		s.Contains(string(resp.Body()), "invalid email")
	})

	s.Run("Should return 400 when document is invalid", func() {
		invalidBody := request.CreateLeadRequest{
			Email:          "valid@email.com",
			DocumentNumber: "123",
			Name:           "Fake Name",
			Phone:          "99990000000",
		}

		resp, err := s.TestUtil.RequestWithDefaultAccessToken(s.T().Context(), util.RequestParams{
			Path:   "/v1/leads",
			Method: http.MethodPost,
			Body:   invalidBody,
		})

		s.NoError(err)
		s.Equal(http.StatusBadRequest, resp.StatusCode())

		s.Contains(string(resp.Body()), "invalid document")
	})

	s.Run("Should create a lead successfully", func() {
		validBody := request.CreateLeadRequest{
			Email:          "fake@email.com",
			DocumentNumber: "12345678909",
			Name:           "Fake Name",
			Phone:          "99990000000",
		}

		resp, err := s.TestUtil.RequestWithDefaultAccessToken(s.T().Context(), util.RequestParams{
			Path:   "/v1/leads",
			Method: http.MethodPost,
			Body:   validBody,
		})

		s.NoError(err)
		s.Equal(http.StatusCreated, resp.StatusCode())

		var respBody response.CreateLeadResponse
		err = resp.Unmarshal(&respBody)
		s.NoError(err)

		s.NotEmpty(respBody.AccessToken)
		s.NotEmpty(respBody.ExpiresIn)
		s.NotEmpty(respBody.LeadID)
	})
}
