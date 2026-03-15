package events

const (
	MessagesChannel     = "gallery.messages.Message"
	TranslationsChannel = "gallery.messages.Translation"
	CommentsChannel     = "gallery.messages.Comment"

	MessageCreatedEvent   = "messagesapi.MessageCreated"
	MessageDeletedEvent   = "messagesapi.MessageDeleted"
	MessageUpdatedEvent   = "messagesapi.MessageUpdated"
	MessagesPublishEvent  = "messagesapi.MessagesPublished"
	MessagesPrivateEvent  = "messagesapi.MessagesPrivated"

	TranslationCreatedEvent   = "messagesapi.TranslationCreated"
	TranslationDeletedEvent   = "messagesapi.TranslationDeleted"
	TranslationUpdatedEvent   = "messagesapi.TranslationUpdated"

	CommentCreatedEvent         = "messagesapi.CommentCreated"
	CommentUpdatedEvent         = "messagesapi.CommentUpdated"
	CommentDeletedEvent         = "messagesapi.CommentDeleted"
	MessageCommentsDeletedEvent = "messagesapi.MessageCommentsDeleted"
)
