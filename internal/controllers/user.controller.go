package controllers

import "sigo/internal/services"

func SendMessage(user *services.User, message []byte) {
	user.Receiver <- message
}

func ReadMessage(user *services.User) []byte {
	return <-user.Sender
}
