package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/nit-app/nit-backend/models"
	"github.com/nit-app/nit-backend/models/requests"
	"github.com/nit-app/nit-backend/models/responses"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func makeEvent(t *testing.T, client *resty.Client) string {
	draftToCreate := models.EventHeader{
		Title:            "Mock Event for Setup",
		PriceLow:         0,
		PriceHigh:        230,
		AgeLimitLow:      10,
		AgeLimitHigh:     76,
		Location:         "Testing department",
		OwnerInfo:        "Main tester",
		Tags:             []string{"testing", "setup"},
		PlainDescription: "Come and test",
	}

	schedule := []*models.EventSchedule{
		{
			BeginsAt: time.Date(2025, time.January, 7, 0, 0, 0, 0, time.UTC),
			EndsAt:   time.Date(2025, time.January, 9, 0, 0, 0, 0, time.UTC),
		},
	}

	resp, err := client.R().SetResult(new(responses.BaseResponse[models.EventHeader])).SetBody(draftToCreate).Post("/eventAdmin/create")
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode())

	draft, ok := resp.Result().(*responses.BaseResponse[models.EventHeader])

	require.True(t, ok)

	uuid := draft.Object.UUID

	setSchedule := requests.SetSchedule{EventUUID: uuid, Schedule: schedule}

	resp, err = client.R().SetBody(setSchedule).Post("/eventAdmin/setSchedule")
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode())

	resp, err = client.R().SetPathParam("uuid", uuid).Post("/eventAdmin/publish/{uuid}")
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode())

	return uuid
}
