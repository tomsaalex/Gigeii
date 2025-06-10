package dto

import "example.com/model"


func ToResellerResponse(r model.Reseller) ResellerResponse {
	return ResellerResponse{
		Id:       r.Id.String(),
		Name:     r.Name,
		Username: r.Username,
		Email:    r.Email,
	}
}
