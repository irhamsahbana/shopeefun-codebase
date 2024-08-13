package entity

type LoginRequest struct {
	Provider string  `json:"provider" validate:"required,oneof=google apple"`
	Email    string  `json:"email" validate:"required,email"`
	Fname    string  `json:"first_name" validate:"required"`
	Lname    string  `json:"last_name" validate:"required"`
	PicURL   *string `json:"picture_url"`
}

type UserInfoResponse struct {
	Id            string  `json:"id"`
	Email         string  `json:"email"`
	Name          string  `json:"name"`
	GivenName     string  `json:"given_name"`
	FamilyName    string  `json:"family_name"`
	Locale        string  `json:"locale"`
	PicURL        *string `json:"picture"`
	VerifiedEmail bool    `json:"verified_email"`
}
