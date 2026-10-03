package domain

type Kind uint8

const (
	KindUnknown Kind = iota
	KindNotFound
	KindValidationError
	KindAlreadyExists
	KindAuthorizationError
	KindForbidden
	KindBadRequest
	KindConflict
)

type Error struct {
	Kind    Kind
	Message string
	Cause   error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func Unknown(message string) error {
	return &Error{
		Kind:    KindUnknown,
		Message: message,
	}
}

func ValidationError(message string) error {
	return &Error{
		Kind:    KindValidationError,
		Message: message,
	}
}

func AlreadyExists(message string) error {
	return &Error{
		Kind:    KindAlreadyExists,
		Message: message,
	}
}

func AuthorizationError() error {
	return &Error{
		Kind:    KindAuthorizationError,
		Message: "authorization error",
	}
}

func ForbiddenError() error {
	return &Error{
		Kind:    KindForbidden,
		Message: "forbidden",
	}
}

func NotFound(message string) error {
	return &Error{
		Kind:    KindNotFound,
		Message: message,
	}
}

func Conflict(message string) error {
	return &Error{
		Kind:    KindConflict,
		Message: message,
	}
}
