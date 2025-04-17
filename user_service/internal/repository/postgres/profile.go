package psqlrepoprofile

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	dto "github.com/DENFNC/Zappy/user_service/internal/dto/profile"
	psql "github.com/DENFNC/Zappy/user_service/internal/storage/postgres"
	"github.com/doug-martin/goqu/v9"
)

type ProfileRepo struct {
	*psql.Storage
	goqu *goqu.DialectWrapper
}

type cursor struct {
	ID        int32     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func New(
	db *psql.Storage,
	goqu *goqu.DialectWrapper,
) *ProfileRepo {
	return &ProfileRepo{
		db,
		goqu,
	}
}
func (r *ProfileRepo) Create(ctx context.Context, profile *models.Profile) (uint32, error) {
	var profileID uint32

	stmt, args, err := r.goqu.
		Insert("profile").
		Returning(goqu.C("profile_id")).
		Rows(goqu.Record{
			"auth_user_id": profile.AuthUserID,
			"first_name":   profile.FirstName,
			"last_name":    profile.LastName,
		}).
		Prepared(true).
		ToSQL()

	if err != nil {
		return 0, fmt.Errorf("failed to build SQL: %w", err)
	}

	err = r.DB.QueryRow(ctx, stmt, args...).Scan(&profileID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert profile: %w", err)
	}

	return profileID, nil
}

func (r *ProfileRepo) GetByID(ctx context.Context, id uint32) (*models.Profile, error) {
	var profile models.Profile

	stmt, args, err := r.goqu.
		Select(
			"profile_id",
			"auth_user_id",
			"first_name",
			"last_name",
			"created_at",
			"updated_at",
		).
		From("profile").
		Where(goqu.C("profile_id").Eq(id)).
		Prepared(true).
		ToSQL()

	if err != nil {
		return nil, fmt.Errorf("failed to build SQL: %w", err)
	}

	err = r.DB.QueryRow(ctx, stmt, args...).Scan(
		&profile.ProfileID,
		&profile.AuthUserID,
		&profile.FirstName,
		&profile.LastName,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch profile: %w", err)
	}

	return &profile, nil
}

func (r *ProfileRepo) List(
	ctx context.Context,
	params *dto.ListParams,
) (
	*dto.ListResult,
	error,
) {
	var cur cursor
	if params.PageToken != "" {
		data, err := base64.RawURLEncoding.DecodeString(params.PageToken)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, &cur); err != nil {
			return nil, err
		}
	}

	stmt, args, err := r.goqu.Select(
		"profile_id",
		"auth_user_id",
		"first_name",
		"last_name",
		"created_at",
		"updated_at",
	).From("profile").Where(
		goqu.L(
			"(?::timestamptz, ?::bigint) < (created_at, profile_id)",
			cur.CreatedAt, cur.ID,
		),
	).Order(
		goqu.I("created_at").Asc(),
		goqu.I("profile_id").Asc(),
	).Limit(uint(params.PageSize)).Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	profiles := make([]*models.Profile, 0, params.PageSize)

	var lastCur cursor
	for rows.Next() {
		var p models.Profile
		if err := rows.Scan(
			&p.ProfileID,
			&p.AuthUserID,
			&p.FirstName,
			&p.LastName,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		lastCur = cursor{CreatedAt: p.CreatedAt, ID: int32(p.ProfileID)}
		profiles = append(profiles, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	var nextToken string
	if len(profiles) > 0 {
		raw, _ := json.Marshal(lastCur)
		nextToken = base64.RawURLEncoding.EncodeToString(raw)
	}
	return &dto.ListResult{
		Items:         profiles,
		NextPageToken: nextToken,
	}, nil
}

func (r *ProfileRepo) Update(ctx context.Context, profile *models.Profile) (uint32, error) {
	stmt, args, err := r.goqu.
		Update("profile").
		Returning(goqu.C("profile_id")).
		Set(goqu.Record{
			"first_name": profile.FirstName,
			"last_name":  profile.LastName,
			"updated_at": profile.UpdatedAt,
		}).
		Where(goqu.C("profile_id").Eq(profile.ProfileID)).
		Prepared(true).
		ToSQL()

	if err != nil {
		return 0, fmt.Errorf("failed to build SQL: %w", err)
	}

	var profileID uint32
	err = r.DB.QueryRow(ctx, stmt, args...).Scan(&profileID)
	if err != nil {
		return 0, fmt.Errorf("failed to update profile: %w", err)
	}

	return profileID, nil
}

func (r *ProfileRepo) Delete(ctx context.Context, id uint32) (uint32, error) {
	stmt, args, err := r.goqu.
		Delete("profile").
		Returning(goqu.C("profile_id")).
		Where(goqu.C("profile_id").Eq(id)).
		Prepared(true).
		ToSQL()

	if err != nil {
		return 0, fmt.Errorf("failed to build SQL: %w", err)
	}

	var profileID uint32
	err = r.DB.QueryRow(ctx, stmt, args...).Scan(&profileID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete profile: %w", err)
	}

	return profileID, nil
}
