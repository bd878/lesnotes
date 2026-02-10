package events

const (
	MessagesChannel = "gallery.messages.Message"
	TranslationsChannel = "gallery.messages.Translation"

	MessageCreatedEvent   = "messagesapi.MessageCreated"
	MessageDeletedEvent   = "messagesapi.MessageDeleted"
	MessageUpdatedEvent   = "messagesapi.MessageUpdated"
	MessagesPublishEvent  = "messagesapi.MessagesPublished"
	MessagesPrivateEvent  = "messagesapi.MessagesPrivated"

	TranslationCreatedEvent   = "messagesapi.TranslationCreated"
	TranslationDeletedEvent   = "messagesapi.TranslationDeleted"
	TranslationUpdatedEvent   = "messagesapi.TranslationUpdated"
)
