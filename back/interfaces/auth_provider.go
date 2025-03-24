package interfaces

type AuthProvider interface {
	SignUp(email string, hashedPassword string) CustomError
	SignIn(id string, hashedPassword string) (string, error)
}
