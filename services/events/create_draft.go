package events

import (
	"context"
	"github.com/google/uuid"
	"github.com/nit-app/nit-backend/env"
	"github.com/nit-app/nit-backend/models"
)

func CreateDraft(ctx context.Context, header *models.EventHeader) (*models.EventHeader, error) {
	// set draft manually
	header.IsDraft = true

	// create uuid
	header.UUID = uuid.New().String()

	// todo
	env.DB().ExecContext(ctx, "insert into events () values ()")

	panic("implement me!!!")
}
