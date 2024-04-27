package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/nit-app/nit-backend/env"
	wrappedErrors "github.com/nit-app/nit-backend/errors"
	"github.com/nit-app/nit-backend/models"
	"github.com/nit-app/nit-backend/models/status"
)

func GetByUuid(ctx context.Context, u string) (*models.User, error) {
	row := env.DB().QueryRowContext(ctx, `select "uuid", phoneNumber, firstName, lastName,
       registeredAt, isAdmin from users where "uuid" = ?`, u)

	return scanUser(row)
}

func scanUser(row env.Scanner) (*models.User, error) {
	u := &models.User{}

	err := row.Scan(&u.UUID, &u.PhoneNumber, &u.FirstName, &u.LastName, &u.RegisteredAt, &u.IsAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, wrappedErrors.New(status.NoSuchEvent, err)
		}

		return nil, wrappedErrors.New(status.InternalServerError, err)
	}

	return u, nil
}
