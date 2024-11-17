package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/nit-app/nit-backend/models"
	"github.com/nit-app/nit-backend/models/responses"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type EventSuite struct {
	suite.Suite
	client   *resty.Client
	realUuid string
}

func (s *EventSuite) SetupSuite() {
	s.client = makeClient(s.T())
	s.realUuid = "f70dd14f-8e24-11ee-8542-fa163e445fa2"
}

func (s *EventSuite) TestGetEventByUUID() {
	resp, err := s.client.R().SetPathParam("uuid", s.realUuid).SetResult(new(responses.BaseResponse[models.Event])).Get("/events/get/{uuid}")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusOK)

	event, ok := resp.Result().(*responses.BaseResponse[models.Event])
	s.Require().True(ok)

	s.Require().NotEmpty(event.Object.EventHeader)
	s.Require().NotEmpty(event.Object.Schedule)
	s.Require().NotEmpty(event.Object.Tags)
	s.Require().NotEmpty(event.Object.Links)
}

func (s *EventSuite) TestGetEventByWrongUUID() {
	resp, err := s.client.R().SetPathParam("uuid", uuid.New().String()).SetResult(new(responses.BaseResponse[models.Event])).Get("/events/get/{uuid}")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusNotFound, "searching by a wrong UUID should raise an error")
}

func (s *EventSuite) TestGetEventByNonUUID() {
	resp, err := s.client.R().SetPathParam("uuid", "uuid").SetResult(new(responses.BaseResponse[models.Event])).Get("/events/get/{uuid}")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusBadRequest, "searching by a non-UUID param should raise an error")
}

func TestEvent(t *testing.T) {
	suite.Run(t, new(EventSuite))
}
