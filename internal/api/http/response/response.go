package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"familyquest-backend/internal/domain"
)

type PublicResponder struct{}

func NewPublicResponder() *PublicResponder {
	return &PublicResponder{}
}

type Success[T any] struct {
	Data T `json:"data"`
}

type ErrorResponse struct {
	Error ErrorData `json:"error"`
}

type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (r *PublicResponder) Success(w http.ResponseWriter, statusCode int, data any) {
	r.writeJSON(w, statusCode, Success[any]{
		Data: data,
	})
}

func (r *PublicResponder) Error(w http.ResponseWriter, err error) {
	statusCode, code, message := r.errorData(err)
	r.writeJSON(w, statusCode, ErrorResponse{
		Error: ErrorData{
			Code:    code,
			Message: message,
		},
	})
}

func (r *PublicResponder) errorData(err error) (int, string, string) {
	var domainErr *domain.Error
	if !errors.As(err, &domainErr) {
		return http.StatusInternalServerError, "INTERNAL_ERROR", "Внутренняя ошибка"
	}

	switch domainErr.Kind {
	case domain.KindValidationError, domain.KindBadRequest:
		return http.StatusBadRequest, "VALIDATION_ERROR", domainErr.Message
	case domain.KindAlreadyExists:
		return http.StatusConflict, "LOGIN_ALREADY_EXISTS", "Логин уже занят"
	case domain.KindAuthorizationError:
		return http.StatusUnauthorized, "AUTHORIZATION_ERROR", "Неверный логин или пароль"
	case domain.KindForbidden:
		return http.StatusForbidden, "FORBIDDEN", domainErr.Message
	case domain.KindNotFound:
		return http.StatusNotFound, "NOT_FOUND", domainErr.Message
	case domain.KindConflict:
		return http.StatusConflict, "CONFLICT", domainErr.Message
	default:
		return http.StatusInternalServerError, "INTERNAL_ERROR", "Внутренняя ошибка"
	}
}

func (r *PublicResponder) writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}
