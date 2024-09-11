package account

import (
	"shared"

	validation "github.com/go-ozzo/ozzo-validation"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r LoginRequest) Validate() error {
	if err := validation.Validate(r.Username, validation.Required); err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "username required",
		}
	}

	if err := validation.Validate(r.Password, validation.Required); err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "password required",
		}
	}

	return nil
}

type LoginResponse struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}
