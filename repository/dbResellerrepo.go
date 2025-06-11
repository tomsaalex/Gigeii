package repository

import (
	"context"

	"example.com/db"
	"example.com/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type DBResellerRepository struct {
	queries *db.Queries
	mapper  ResellerMapperDB
}

// Add implements ResellerRepository.
func (d *DBResellerRepository) Add(ctx context.Context, reseller model.Reseller) (*model.Reseller, error) {
	addResellerParams:= d.mapper.ResellerToAddResellerParams(reseller)
	dbReseller, err := d.queries.AddReseller(ctx, addResellerParams)
	if err != nil {
		return nil, &DuplicateEntityError{
			Message: "Add: Email or Username is already taken by a different reseller.",
		}
	}
	modelReseller := d.mapper.DBResellerToReseller(dbReseller)
	return modelReseller, nil
}

// DeleteReseller implements ResellerRepository.
func (d *DBResellerRepository) DeleteReseller(ctx context.Context, id uuid.UUID) error {
	dbID := pgtype.UUID{Bytes: id, Valid: true}
	err := d.queries.DeleteReseller(ctx, dbID)
	if err != nil {
		return &EntityNotFoundError{
			Message: "DeleteReseller: No reseller matches ID \"" + id.String() + "\"",
		}
	}
	return nil
}

// GetAll implements ResellerRepository.
func (d *DBResellerRepository) GetAll(ctx context.Context) ([]model.Reseller, error) {
	dbReseller, err:= d.queries.SelectResellers(ctx)
	if err != nil {
		return nil, &EntityNotFoundError{
			Message: "GetAll: No resellers found.",
		}
	}
	modelResellers := make([]model.Reseller, len(dbReseller))
	for i, r := range dbReseller {
		model := d.mapper.DBResellerToReseller(r)
		modelResellers[i] = *model
	}
	return modelResellers, nil
}

// GetByEmail implements ResellerRepository.
func (d *DBResellerRepository) GetByEmail(ctx context.Context, email string) (*model.Reseller, error) {
	emailParam := pgtype.Text{
		String: email,
		Valid:  email != "",
	}

	dbReseller, err := d.queries.GetResellerByEmail(ctx, emailParam)
	if err != nil {
		return nil, &EntityNotFoundError{
			Message: "GetResellerByEmail: No reseller matches email \"" + email + "\"",
		}
	}

	modelReseller := d.mapper.DBResellerToReseller(dbReseller)
	return modelReseller, nil
}

// GetById implements ResellerRepository.
func (d *DBResellerRepository) GetById(ctx context.Context, id uuid.UUID) (*model.Reseller, error) {
	dbID := pgtype.UUID{Bytes: id, Valid: true}
	dbReseller, err := d.queries.GetResellerByID(ctx, dbID)
	if err != nil {
		return nil, &EntityNotFoundError{
			Message: "GetResellerByID: No reseller matches ID \"" + id.String() + "\"",
		}
	}

	modelReseller := d.mapper.DBResellerToReseller(dbReseller)
	return modelReseller, nil
}

// GetByUsername implements ResellerRepository.
func (d *DBResellerRepository) GetByUsername(ctx context.Context, username string) (*model.Reseller, error) {
	dbReseller, err := d.queries.GetResellerByUsername(ctx, username)
	if err != nil {
		return nil, &EntityNotFoundError{
			Message: "GetResellerByUsername: No reseller matches username \"" + username + "\"",
		}
	}

	modelReseller := d.mapper.DBResellerToReseller(dbReseller)
	return modelReseller, nil
}
