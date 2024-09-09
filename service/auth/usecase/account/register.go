package account

import (
	"shared"

	validation "github.com/go-ozzo/ozzo-validation"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r RegisterRequest) Validate() error {
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

	if err := validation.Validate(r.Name, validation.Required); err != nil {
		return &shared.Error{
			HttpStatusCode: 400,
			Message:        "name required",
		}
	}

	return nil
}
