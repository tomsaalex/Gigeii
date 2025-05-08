package handler

import "example.com/model"

type UserDTOMapper struct {
}

func (m *UserDTOMapper) RegistrationDTOToUser(userDTO UserRegistrationDTO) model.User {
	return model.User{
		Username: userDTO.Username,
		Email:    userDTO.Email,
		PassHash: nil,
		PassSalt: nil,
	}
}
