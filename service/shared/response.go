package shared

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func SuccessResponse(
	c echo.Context,
	message string,
	data interface{},
) error {
	return c.JSON(200, &Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func FailResponse(
	c echo.Context,
	message string,
	errorData interface{},
) error {
	return c.JSON(400, &Response{
		Success: false,
		Message: message,
		Error:   errorData,
	})
}

func FailResponseFromCustomError(
	c echo.Context,
	err error,
) error {

	if grpcErr, ok := status.FromError(err); ok {

		var httpStatus int

		switch grpcErr.Code() {
		case codes.InvalidArgument:
			httpStatus = http.StatusBadRequest
		case codes.NotFound:
			httpStatus = http.StatusNotFound
		case codes.AlreadyExists:
			httpStatus = http.StatusConflict
		default:
			httpStatus = http.StatusInternalServerError
		}

		err = &Error{
			HttpStatusCode: httpStatus,
			GrpcStatusCode: grpcErr.Code(),
			Message:        grpcErr.Message(),
		}
	}

	switch v := err.(type) {
	case *Error:
		return c.JSON(v.HttpStatusCode, &Response{
			Success: false,
			Message: v.Message,
			Data:    v.Data,
		})
	case Error:
		return c.JSON(v.HttpStatusCode, &Response{
			Success: false,
			Message: v.Message,
			Data:    v.Data,
		})
	}

	return c.JSON(500, &Response{
		Success: false,
		Message: "internal server error",
	})
}
