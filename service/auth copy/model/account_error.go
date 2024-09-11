package model

import "shared"

var ErrUserNotFound = shared.Error{
	HttpStatusCode: 400,
	Message:        "user not found",
}
