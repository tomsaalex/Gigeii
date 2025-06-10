package repository

import (
	"example.com/db"
	"example.com/model"
	"github.com/jackc/pgx/v5/pgtype"
)

type ResellerMapperDB struct {
}

func (m *ResellerMapperDB) ResellerToAddResellerParams(reseller model.Reseller) db.AddResellerParams {
	
	email := pgtype.Text{
		String: reseller.Email,
		Valid:  reseller.Email != "",
	}

	return db.AddResellerParams{
		Name	:	 reseller.Name,
		Username: reseller.Username,
		Email:   email,
		PasswordHash: string(reseller.PasswordHash),
		
	}
}


func (m *ResellerMapperDB) DBResellerToReseller(dbReseller db.Reseller) *model.Reseller {
	reseller := model.Reseller{
		Id:          dbReseller.ID.Bytes,
		Name:        dbReseller.Name,
		Username:    dbReseller.Username,
		PasswordHash: []byte(dbReseller.PasswordHash),
		Email:       dbReseller.Email.String,
	}

	return &reseller
}
