package repo

import (
	"context"

	"github.com/DENFNC/Zappy/internal/domain/models"
	psql "github.com/DENFNC/Zappy/internal/storage/postgres"
)

type User struct {
	psql.Storage
}

func (u *User) Create(
	ctx context.Context,
	user *models.User,
) error {
	return nil
}
