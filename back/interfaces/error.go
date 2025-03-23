package interfaces

type CustomError interface {
	Error() string
	Unwrap() error
}
