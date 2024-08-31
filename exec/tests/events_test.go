package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/nit-app/nit-backend/models"
	"github.com/nit-app/nit-backend/models/requests"
	"github.com/nit-app/nit-backend/models/responses"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
	"time"
)

type EventsSuite struct {
	suite.Suite
	client *resty.Client
}

func (s *EventsSuite) SetupSuite() {
	s.client = makeClient(s.T())
}

func (s *EventsSuite) TestBasicLookup() {
	filter := requests.EventLookupFilters{
		From: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		To:   time.Date(2050, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	resp, err := s.client.R().SetResult(new(responses.BaseResponse[[]models.EventHeader])).SetBody(filter).Post("/events/lookup")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusOK)

	events, ok := resp.Result().(*responses.BaseResponse[[]models.EventHeader])
	s.Require().True(ok)

	s.Require().True(len(events.Object) > 0, "there should be at least one event")

	s.Require().NotEmpty(events.Object[0].UUID, "event uuid should be set")
	s.Require().NotEmpty(events.Object[0].Title, "event title should be set")
}

func (s *EventsSuite) TestAheadLookup() {
	aheadFilter := requests.EventLookupFilters{
		From: time.Date(2050, time.January, 1, 0, 0, 0, 0, time.UTC),
		To:   time.Date(2055, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	resp, err := s.client.R().SetResult(new(responses.BaseResponse[[]models.EventHeader])).SetBody(aheadFilter).Post("/events/lookup")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusOK)

	events, ok := resp.Result().(*responses.BaseResponse[[]models.EventHeader])
	s.Require().True(ok)

	s.Require().True(len(events.Object) == 0, "there shouldn't be any forthcoming event")
}

func (s *EventsSuite) TestExcludePaid() {
	filter := requests.EventLookupFilters{
		From:        time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		To:          time.Date(2050, time.January, 1, 0, 0, 0, 0, time.UTC),
		ExcludePaid: true,
	}

	resp, err := s.client.R().SetResult(new(responses.BaseResponse[[]models.EventHeader])).SetBody(filter).Post("/events/lookup")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusOK)

	events, ok := resp.Result().(*responses.BaseResponse[[]models.EventHeader])
	s.Require().True(ok)

	s.Require().True(len(events.Object) > 0, "there should be at least one event")

	for _, event := range events.Object {
		s.Require().True(event.PriceLow == 0 && event.PriceHigh == 0, "events should be free")
	}
}

func TestEvents(t *testing.T) {
	suite.Run(t, new(EventsSuite))
}
