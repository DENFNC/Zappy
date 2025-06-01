package errpkg

type AppError struct {
	Code    string
	Message string
	Err     error
}

func (ae *AppError) Error() string {
	return ae.Message
}

func (ae *AppError) Unwrap() error {
	return ae.Err
}

func New(
	code string,
	message string,
	err error,
) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

var (
	ErrNotFound        = New("NOT_FOUND", "not found", nil)
	ErrInvalidArgument = New("INVALID_ARGUMENT", "invalid argument", nil)
	ErrInternal        = New("INTERNAL", "internal server error", nil)
	ErrConstraint      = New("ERR_FOREIGN_KEY_VIOLATION", "external switch violation", nil)
	ErrUniqueViolation = New("ERR_UNIQUE_VIOLATION", "violation of uniqueness", nil)
)
