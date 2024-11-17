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

	events := s.lookupEvents(filter)

	s.Require().True(len(events.Object) > 0, "there should be at least one event")

	s.Require().NotEmpty(events.Object[0].UUID, "event uuid should be set")
	s.Require().NotEmpty(events.Object[0].Title, "event title should be set")
}

func (s *EventsSuite) TestAheadLookup() {
	aheadFilter := requests.EventLookupFilters{
		From: time.Date(2050, time.January, 1, 0, 0, 0, 0, time.UTC),
		To:   time.Date(2055, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	events := s.lookupEvents(aheadFilter)

	s.Require().True(len(events.Object) == 0, "there shouldn't be any forthcoming event")
}

func (s *EventsSuite) TestExcludePaid() {
	filter := requests.EventLookupFilters{
		From:        time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		To:          time.Date(2050, time.January, 1, 0, 0, 0, 0, time.UTC),
		ExcludePaid: true,
	}
	events := s.lookupEvents(filter)

	s.Require().True(len(events.Object) > 0, "there should be at least one event")

	for _, event := range events.Object {
		s.Require().True(event.PriceLow == 0, "events should be free")
	}
}

func (s *EventsSuite) TestExcludeAgeRestricted() {
	filter := requests.EventLookupFilters{
		From:                 time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		To:                   time.Date(2050, time.January, 1, 0, 0, 0, 0, time.UTC),
		ExcludeAgeRestricted: true,
	}

	events := s.lookupEvents(filter)

	s.Require().True(len(events.Object) > 0, "there should be at least one event")

	for _, event := range events.Object {
		s.Require().True(event.AgeLimitLow == 0 && event.AgeLimitHigh == 0, "events should not be restricted by age")
	}
}

func (s *EventsSuite) TestTagsLookup() {
	tag := "test"
	filter := requests.EventLookupFilters{
		From: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		To:   time.Date(2050, time.January, 1, 0, 0, 0, 0, time.UTC),
		Tags: []string{tag},
	}

	events := s.lookupEvents(filter)

	s.Require().True(len(events.Object) > 0, "there should be at least one event")

	for _, event := range events.Object {
		s.Require().Contains(event.Tags, tag, "events should have the tag from filters")
	}
}

func (s *EventsSuite) TestTagsEmptyLookup() {
	tag := "iAmNotATag,Actually"
	filter := requests.EventLookupFilters{
		From: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		To:   time.Date(2050, time.January, 1, 0, 0, 0, 0, time.UTC),
		Tags: []string{tag},
	}

	events := s.lookupEvents(filter)

	s.Require().Empty(events.Object, "there should be no events found")
}

func (s *EventsSuite) lookupEvents(filter requests.EventLookupFilters) *responses.BaseResponse[[]models.EventHeader] {
	resp, err := s.client.R().SetResult(new(responses.BaseResponse[[]models.EventHeader])).SetBody(filter).Post("/events/lookup")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusOK)

	events, ok := resp.Result().(*responses.BaseResponse[[]models.EventHeader])
	s.Require().True(ok)

	return events
}

func TestEvents(t *testing.T) {
	suite.Run(t, new(EventsSuite))
}
