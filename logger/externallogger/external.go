package externallogger

type ExternalLogger interface {
	Sendlog(chatId int, message string) error
}