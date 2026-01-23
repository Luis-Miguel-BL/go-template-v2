package api

import (
	"net/http"

	"github.com/Luis-Miguel-BL/go-lm-template/test/util"
)

func (s *APISuite) TestHealthEndpoint() {
	s.Run("Should return 200", func() {
		req, _ := s.TestUtil.DoRequest(s.T().Context(), util.RequestParams{
			Path:   "/health",
			Method: http.MethodGet,
		})

		s.True(req.IsSuccess())
		s.Equal(http.StatusOK, req.StatusCode())
	})
}
