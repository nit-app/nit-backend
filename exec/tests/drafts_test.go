package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/nit-app/nit-backend/models"
	"github.com/nit-app/nit-backend/models/requests"
	"github.com/nit-app/nit-backend/models/responses"
	"github.com/stretchr/testify/suite"
	"net/http"
	"sort"
	"testing"
	"time"
)

type DraftsSuite struct {
	suite.Suite
	client        *resty.Client
	draftToCreate models.EventHeader
	draft         models.EventHeader
	links         []*models.EventExternalLink
	schedule      []*models.EventSchedule
}

func (s *DraftsSuite) SetupSuite() {
	s.client = makeClient(s.T())
	s.draftToCreate = models.EventHeader{
		Title:            "Testing new draft",
		PriceLow:         0,
		PriceHigh:        100,
		AgeLimitLow:      14,
		AgeLimitHigh:     99,
		Location:         "Testing department",
		OwnerInfo:        "Main tester",
		Tags:             []string{"testing"},
		PlainDescription: "Come and test",
	}

	s.draft = s.createDraft().Object
	s.links = []*models.EventExternalLink{{Title: "link1", URL: "https://vk.com/durov"}, {Title: "link2", URL: "https://pkg.go.dev/testing"}}
	s.schedule = []*models.EventSchedule{
		{BeginsAt: time.Date(2025, time.January, 7, 0, 0, 0, 0, time.UTC),
			EndsAt: time.Date(2025, time.January, 9, 0, 0, 0, 0, time.UTC)},
		{BeginsAt: time.Date(2025, time.January, 11, 0, 0, 0, 0, time.UTC),
			EndsAt: time.Date(2025, time.January, 13, 0, 0, 0, 0, time.UTC)}}

}

func (s *DraftsSuite) TestDraftCreation() {
	draft := s.createDraft()

	s.Require().NotEmpty(draft.Object.UUID, "draft uuid should be set")
	s.Require().True(draft.Object.IsDraft, "isDraft should be True")
	s.Require().Equal(draft.Object.Title, s.draftToCreate.Title, "title should be saved")
	s.Require().Empty(draft.Object.Schedule, "drafts schedule should be empty")

	draftResult := *s.getByUUID(draft.Object.UUID).Object.EventHeader

	s.Require().Equal(draftResult, draft.Object, "draft must be saved correctly")
}

func (s *DraftsSuite) TestFalseDraftCreation() {
	emptyDraft := models.EventHeader{
		Title: "Testing empty draft",
	}

	resp, err := s.client.R().SetResult(new(responses.BaseResponse[models.EventHeader])).SetBody(emptyDraft).Post("/eventAdmin/create")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusBadRequest, "draft without essential info should not be accepted")

	wrongDraft := models.EventHeader{
		Title:            "Testing wrong draft",
		PriceLow:         300,
		PriceHigh:        100,
		AgeLimitLow:      99,
		AgeLimitHigh:     14,
		Location:         "Testing department",
		OwnerInfo:        "Main tester",
		Tags:             []string{"testing"},
		PlainDescription: "Come and test",
	}

	_, err = s.client.R().SetResult(new(responses.BaseResponse[models.EventHeader])).SetBody(wrongDraft).Post("/eventAdmin/create")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusBadRequest, "draft with incorrect data should not be accepted")
}

func (s *DraftsSuite) TestSetTags() {
	s.setTags(s.draft.UUID)

	setDuplicateTags := requests.SetTags{EventUUID: s.draft.UUID, Tags: []string{"Best", "best"}}

	resp, err := s.client.R().SetBody(setDuplicateTags).Post("/eventAdmin/setTags")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusBadRequest, "duplicate tags should not be set")

	setTagsFalseEvent := requests.SetTags{EventUUID: uuid.New().String(), Tags: []string{"coolest", "best"}}

	resp, err = s.client.R().SetBody(setTagsFalseEvent).Post("/eventAdmin/setTags")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusNotFound, "non existent events should not get tags")
}

func (s *DraftsSuite) TestSetLinks() {
	s.setLinks(s.draft.UUID)

	setLinksFalseEvent := requests.SetLinks{EventUUID: uuid.New().String(), Links: s.links}

	resp, err := s.client.R().SetBody(setLinksFalseEvent).Post("/eventAdmin/setLinks")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusNotFound, "non existent events should not get links")

	duplicateLinks := []*models.EventExternalLink{{Title: "link1", URL: "https://vk.com/durov"}, {Title: "link2", URL: "https://vk.com/durov"}}

	setDuplicateLinks := requests.SetLinks{EventUUID: s.draft.UUID, Links: duplicateLinks}

	resp, err = s.client.R().SetBody(setDuplicateLinks).Post("/eventAdmin/setLinks")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusBadRequest, "duplicate links should not be set")

	nonURLLinks := []*models.EventExternalLink{{Title: "link1", URL: "durov is my man"}}

	setNonURLLinks := requests.SetLinks{EventUUID: s.draft.UUID, Links: nonURLLinks}

	resp, err = s.client.R().SetBody(setNonURLLinks).Post("/eventAdmin/setLinks")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusBadRequest, "non-URL links should not be accepted")
}

func (s *DraftsSuite) TestSetSchedule() {
	s.setSchedule(s.draft.UUID)

	wrongSchedule := []*models.EventSchedule{
		{BeginsAt: time.Date(2025, time.January, 9, 0, 0, 0, 0, time.UTC),
			EndsAt: time.Date(2025, time.January, 7, 0, 0, 0, 0, time.UTC)}}

	setWrongSchedule := requests.SetSchedule{EventUUID: s.draft.UUID, Schedule: wrongSchedule}

	resp, err := s.client.R().SetBody(setWrongSchedule).Post("/eventAdmin/setSchedule")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusBadRequest, "schedule, where start date is later than end date, should not be set")

	duplicateSchedule := []*models.EventSchedule{
		{BeginsAt: time.Date(2025, time.January, 9, 0, 0, 0, 0, time.UTC),
			EndsAt: time.Date(2025, time.January, 7, 0, 0, 0, 0, time.UTC)},
		{BeginsAt: time.Date(2025, time.January, 9, 0, 0, 0, 0, time.UTC),
			EndsAt: time.Date(2025, time.January, 7, 0, 0, 0, 0, time.UTC)}}

	setDuplicateSchedule := requests.SetSchedule{EventUUID: s.draft.UUID, Schedule: duplicateSchedule}

	resp, err = s.client.R().SetBody(setDuplicateSchedule).Post("/eventAdmin/setSchedule")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusBadRequest, "schedule with two identical dates should not be set")

	setScheduleFalseEvent := requests.SetSchedule{EventUUID: uuid.New().String(), Schedule: s.schedule}

	resp, err = s.client.R().SetBody(setScheduleFalseEvent).Post("/eventAdmin/setSchedule")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusNotFound, "non existent events should not get schedule")
}

func (s *DraftsSuite) TestEditDraft() {
	headerToEdit := s.draftToCreate
	headerToEdit.HasCertificate = true
	headerToEdit.Title = "Testing edited draft"

	editDraftRequest := requests.EditDraft{
		EventUUID:   s.draft.UUID,
		Object:      &headerToEdit,
		Description: "Now my draft has a description!",
	}
	resp, err := s.client.R().SetResult(new(responses.BaseResponse[models.EventHeader])).SetBody(editDraftRequest).Post("/eventAdmin/edit")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusOK)

	headerResult, ok := resp.Result().(*responses.BaseResponse[models.EventHeader])
	s.Require().True(ok)

	s.Require().Equal(headerResult.Object.Title, headerToEdit.Title, "new title should be set")
	s.Require().True(headerResult.Object.HasCertificate, "hasCertificate should be set correctly")

	draftResult := s.getByUUID(headerResult.Object.UUID)
	s.Require().Equal(draftResult.Object.Description, editDraftRequest.Description, "description should be set correctly")

	editDraftRequest.EventUUID = uuid.New().String()

	resp, err = s.client.R().SetResult(new(responses.BaseResponse[models.EventHeader])).SetBody(editDraftRequest).Post("/eventAdmin/edit")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusNotFound, "editing a non-existing event should raise an error")
}

func (s *DraftsSuite) TestPublishDraft() {
	draft := s.createDraft().Object

	resp, err := s.client.R().SetPathParam("uuid", draft.UUID).Post("/eventAdmin/publish/{uuid}")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusBadRequest, "drafts without tags and schedule should not be published")

	s.setTags(draft.UUID)
	s.setLinks(draft.UUID)
	s.setSchedule(draft.UUID)

	resp, err = s.client.R().SetPathParam("uuid", draft.UUID).Post("/eventAdmin/publish/{uuid}")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusOK, "drafts with tags, schedule and links can be published")

	publishedEvent := s.getByUUID(draft.UUID).Object
	s.Require().False(publishedEvent.IsDraft, "after publishing isDraft should be false")

	resp, err = s.client.R().SetPathParam("uuid", draft.UUID).Post("/eventAdmin/publish/{uuid}")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusBadRequest, "publishing an event twice should not be allowed")
}

func (s *DraftsSuite) TestGetAllDrafts() {
	resp, err := s.client.R().SetResult(new(responses.BaseResponse[[]models.EventHeader])).Get("/eventAdmin/drafts")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusOK)

	drafts, ok := resp.Result().(*responses.BaseResponse[[]models.EventHeader])
	s.Require().True(ok)

	s.Require().True(len(drafts.Object) > 0, "there should be at least one draft")

	containsAddedDraft := false

	for _, draft := range drafts.Object {
		if draft.UUID == s.draft.UUID {
			containsAddedDraft = true
		}

		s.Require().True(draft.IsDraft, "isDraft must be true")
	}
	s.Require().True(containsAddedDraft, "draft added in SetUp must be returned")
}

func (s *DraftsSuite) getByUUID(uuid string) *responses.BaseResponse[models.Event] {
	resp, err := s.client.R().SetPathParam("uuid", uuid).SetResult(new(responses.BaseResponse[models.Event])).Get("/events/get/{uuid}")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusOK)

	event, ok := resp.Result().(*responses.BaseResponse[models.Event])
	s.Require().True(ok)

	return event
}

func (s *DraftsSuite) createDraft() *responses.BaseResponse[models.EventHeader] {
	resp, err := s.client.R().SetResult(new(responses.BaseResponse[models.EventHeader])).SetBody(s.draftToCreate).Post("/eventAdmin/create")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusOK)

	draft, ok := resp.Result().(*responses.BaseResponse[models.EventHeader])
	s.Require().True(ok)

	return draft
}

func (s *DraftsSuite) setTags(uuid string) {
	setTags := requests.SetTags{EventUUID: uuid, Tags: []string{"coolest", "best"}}

	resp, err := s.client.R().SetBody(setTags).Post("/eventAdmin/setTags")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusOK)

	tagsResult := s.getByUUID(uuid).Object.Tags

	sort.Strings(tagsResult)
	sort.Strings(setTags.Tags)

	s.Require().Equal(tagsResult, setTags.Tags, "set tags must be saved correctly")
}

func (s *DraftsSuite) setLinks(uuid string) {
	setLinks := requests.SetLinks{EventUUID: uuid, Links: s.links}

	resp, err := s.client.R().SetBody(setLinks).Post("/eventAdmin/setLinks")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusOK)

	linksResult := s.getByUUID(uuid).Object.Links

	s.Require().Equal(len(linksResult), len(s.links), "set links must be saved")
}

func (s *DraftsSuite) setSchedule(uuid string) {
	setSchedule := requests.SetSchedule{EventUUID: uuid, Schedule: s.schedule}

	resp, err := s.client.R().SetBody(setSchedule).Post("/eventAdmin/setSchedule")
	s.Require().NoError(err)

	s.Require().Equal(resp.StatusCode(), http.StatusOK)

	scheduleResult := s.getByUUID(uuid).Object.Schedule

	s.Require().Equal(len(scheduleResult), len(s.schedule), "set schedule must be saved")
}

func TestDrafts(t *testing.T) {
	suite.Run(t, new(DraftsSuite))
}
