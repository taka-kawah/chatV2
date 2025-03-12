package usecase

import "back/domain"

type ChatView struct {
	domain.Chat
	UserName string
}
