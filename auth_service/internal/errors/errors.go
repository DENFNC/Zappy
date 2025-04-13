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
	ErrUserNotFound       = New("USR_NOT_FOUND", "user not found", nil)
	ErrInvalidCredentials = New("INVALID_CREDENTIALS", "invalid credentials", nil)
	ErrInternalServer     = New("INTERNAL_SERVER", "internal server error", nil)
	ErrInvalidToken       = New("INVALID_TOKEN", "token is invalid", nil)
	ErrInvalidArgument    = New("INVALID_ARGUMENT", "invalid argument", nil)
)
