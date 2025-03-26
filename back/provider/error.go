package provider

type CustomError interface {
	Error() string
	Unwrap() error
}
