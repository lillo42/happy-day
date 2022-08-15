package common

type (
	Error struct {
		code    string
		message string
	}
)

func NewError(code, message string) Error {
	return Error{
		code:    code,
		message: message,
	}
}

func (e Error) Code() string {
	return e.code
}

func (e Error) Message() string {
	return e.message
}

func (e Error) Error() string {
	return "code: " + e.code + " - message: " + e.message
}
