package httperror

import (
	"net/http"
)

type AppError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func (err AppError) Error() string {
	return err.Message
}

func BadRequestError() AppError {
	return AppError{
		Message: "Bad request",
		Code:    http.StatusBadRequest,
		Data:    nil,
	}
}

func InternalServerError() AppError {
	return AppError{
		Message: "Internal server error",
		Code:    http.StatusInternalServerError,
		Data:    nil,
	}
}

func NotFoundError() AppError {
	return AppError{
		Message: "Resource not found",
		Code:    http.StatusNotFound,
		Data:    nil,
	}
}

func TimeoutError() AppError {
	return AppError{
		Message: "Request timeout",
		Code:    http.StatusRequestTimeout,
		Data:    nil,
	}
}

func UnauthorizedError() AppError {
	return AppError{
		Message: "Unauthorized",
		Code:    http.StatusUnauthorized,
		Data:    nil,
	}
}

func ForbiddenError() AppError {
	return AppError{
		Message: "User not authorized",
		Code:    http.StatusForbidden,
		Data:    nil,
	}
}

func TooManyRequestsError() AppError {
	return AppError{
		Message: "Too many requests",
		Code:    http.StatusTooManyRequests, // 429
		Success: false,
		Data:    nil,
	}
}

func UnsupportedMediaType(msg string) AppError {
	return AppError{
		Message: msg,
		Code:    http.StatusUnsupportedMediaType,
		Data:    nil,
	}
}

func PayloadTooLarge(msg string) AppError {
	return AppError{
		Message: msg,
		Code:    http.StatusRequestEntityTooLarge,
	}
}

func GenericError(message string, code int) AppError {
	return AppError{
		Message: message,
		Code:    code,
		Success: false,
		Data:    nil,
	}
}
