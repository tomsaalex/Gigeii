package repository

import (
	"example.com/db"
	"example.com/model"
)

type UserMapperDB struct {
}

func (m *UserMapperDB) UserToAddUserParams(user model.User) db.AddUserParams {
	return db.AddUserParams{
		Username: user.Username,
		Email:    user.Email,
		PassHash: user.PassHash,
		PassSalt: user.PassSalt,
		Role:     user.Role,
	}
}

func (m *UserMapperDB) DBUserToUser(dbUser db.User) *model.User {
	user := model.User{
		Id:       dbUser.ID.Bytes,
		Email:    dbUser.Email,
		Username: dbUser.Username,
		PassHash: dbUser.PassHash,
		PassSalt: dbUser.PassSalt,
		Role:     dbUser.Role,
	}

	return &user
}
