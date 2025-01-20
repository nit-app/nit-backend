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

type FavSuite struct {
	suite.Suite
	client   *resty.Client
	realUuid string
}

func (s *FavSuite) SetupSuite() {
	s.client = makeClient(s.T())
	s.realUuid = makeEvent(s.T(), s.client)
}

func (s *FavSuite) TestAddFav() {
	s.addFavourite()

	resp, err := s.client.R().SetPathParam("uuid", "nonUuid").Post("/events/fav/{uuid}/add")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusBadRequest, "non uuid params should not be accepted")

	resp, err = s.client.R().SetPathParam("uuid", s.realUuid).Post("/events/fav/{uuid}/add")
	s.Require().NoError(err, "adding an event twice should not raise an error")

	s.Require().Equal(resp.StatusCode(), http.StatusOK)

	resp, err = s.client.R().SetPathParam("uuid", uuid.New().String()).Post("/events/fav/{uuid}/add")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusNotFound, "non existent events should not be added to favourites")
}

func (s *FavSuite) TestGetFav() {
	s.addFavourite()

	events := s.getFavourites()

	s.Require().True(len(events.Object) > 0, "there should be at least one favourite event")

	containsAddedFav := false

	for _, event := range events.Object {
		if event.UUID == s.realUuid {
			containsAddedFav = true
		}
		s.Require().True(event.FavCount > 0, "a favourite event must have a positive FavCount")

		s.Require().False(event.IsDraft, "drafts cannot be shown as favourites")
	}
	s.Require().True(containsAddedFav, "event recently saved to favourites must be returned by get")
}

func (s *FavSuite) TestDeleteFav() {
	s.addFavourite()

	resp, err := s.client.R().SetPathParam("uuid", "nonUuid").Post("/events/fav/{uuid}/remove")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusBadRequest, "non uuid params should not be accepted")

	events := s.getFavourites()

	for _, event := range events.Object {
		resp, err = s.client.R().SetPathParam("uuid", event.UUID).Post("/events/fav/{uuid}/remove")
		s.Require().NoError(err)

		s.Require().Equal(resp.StatusCode(), http.StatusOK)
	}

	events = s.getFavourites()

	s.Require().True(len(events.Object) == 0, "there should be no favourite events after full removal")
}

func (s *FavSuite) addFavourite() {
	resp, err := s.client.R().SetPathParam("uuid", s.realUuid).Post("/events/fav/{uuid}/add")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusOK)
}

func (s *FavSuite) getFavourites() *responses.BaseResponse[[]models.EventHeader] {
	resp, err := s.client.R().SetResult(new(responses.BaseResponse[[]models.EventHeader])).Get("/events/fav")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusOK)

	events, ok := resp.Result().(*responses.BaseResponse[[]models.EventHeader])
	s.Require().True(ok)

	return events
}

func TestFavourites(t *testing.T) {
	suite.Run(t, new(FavSuite))
}
