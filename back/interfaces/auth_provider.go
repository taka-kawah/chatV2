package interfaces

type AuthProvider interface {
	SignUp(email string, hashedPassword string) CustomError
	SignIn(email string, hashedPassword string) (string, CustomError)
}
