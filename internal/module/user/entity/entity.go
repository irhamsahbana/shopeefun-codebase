package entity

type CreateUserRequest struct {
	RoleId string `json:"role_id" validate:"required,uuid"`
	Email  string `json:"email" validate:"required,email"`
	Name   string `json:"name" validate:"required"`
}

type CreateUserResponse struct {
	Id string `json:"id"`
}

type XxxResult struct {
}
