package profileservice

import (
	"context"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	"github.com/DENFNC/Zappy/user_service/internal/domain/repositories"
)

const (
	emptyValue = 0
)

type Profile struct {
	repo repositories.ProfileRepository
}

func New(
	repo repositories.ProfileRepository,
) *Profile {
	return &Profile{
		repo: repo,
	}
}

func (p *Profile) Create(
	ctx context.Context,
	authUserID uint32,
	firstName string,
	lastName string,
) (uint32, error) {
	profile := models.NewProfile(
		authUserID,
		firstName,
		lastName,
	)

	profileID, err := p.repo.Create(ctx, profile)
	if err != nil {
		return emptyValue, err
	}

	return profileID, nil
}

func (p *Profile) Delete(ctx context.Context, profileID uint32) (uint32, error) {
	panic("Implement me!")
}

func (p *Profile) GetByID(ctx context.Context, profileID uint32) (*models.Profile, error) {
	panic("Implement me!")
}

func (p *Profile) List(ctx context.Context) ([]*models.Profile, error) {
	panic("Implement me!")
}

func (p *Profile) Update(ctx context.Context, profileID uint32, firstName, lastName, phone string) (uint32, error) {
	panic("Implement me!")
}
