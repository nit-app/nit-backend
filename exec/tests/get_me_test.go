package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/nit-app/nit-backend/models"
	"github.com/nit-app/nit-backend/models/responses"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type GetMeSuite struct {
	suite.Suite
	client *resty.Client
}

func (s *GetMeSuite) SetupSuite() {
	s.client = makeClient(s.T())
}

func (s *GetMeSuite) TestGetMe() {
	resp, err := s.client.R().SetResult(new(responses.BaseResponse[models.User])).Get("/getMe")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode())

	user, ok := resp.Result().(*responses.BaseResponse[models.User])
	s.Require().True(ok)

	s.Require().NotEmpty(user.Object.UUID, "uuid should not be empty")
	s.Require().True(user.Object.IsAdmin, "test user should be an admin")
}

func TestGetMe(t *testing.T) {
	suite.Run(t, new(GetMeSuite))
}
