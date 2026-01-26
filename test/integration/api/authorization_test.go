package api

import (
	"net/http"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/http/response"
	"github.com/Luis-Miguel-BL/go-lm-template/test/integration/util"
)

func (s *APISuite) TestAuthorizationEndpoint() {
	s.Run("Should return 200 and valid authorization response", func() {
		req, err := s.TestUtil.DoRequest(s.T().Context(), util.RequestParams{
			Path:   "/v1/authorization",
			Method: http.MethodPost,
			Headers: map[string]string{
				"x-api-key": s.TestUtil.GetConfig().Server.AppKey,
			},
		})

		s.NoError(err)

		s.True(req.IsSuccess())
		s.Equal(http.StatusOK, req.StatusCode())
		resBody := response.AuthorizationResponse{}

		err = req.Unmarshal(&resBody)
		s.NoError(err)

		s.NotEmpty(resBody.AccessToken)
		s.Greater(resBody.ExpiresIn, int64(0))
	})

	s.Run("Should return 401 when not authorized", func() {
		req, _ := s.TestUtil.DoRequest(s.T().Context(), util.RequestParams{
			Path:   "/v1/authorization",
			Method: http.MethodPost,
			Headers: map[string]string{
				"x-api-key": "invalid-key",
			},
		})

		s.False(req.IsSuccess())
		s.Equal(http.StatusUnauthorized, req.StatusCode())
	})

}
