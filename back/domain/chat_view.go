package domain

type ChatView struct {
	Chat
	UserName string `validate:"required" db:"user_name"`
}
