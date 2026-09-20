package pkg

const (
	MessagesChannel     = "gallery.messages.Message"
	TranslationsChannel = "gallery.messages.Translation"
	CommentsChannel     = "gallery.messages.Comment"

	CommandChannel = "gallery.messages.commands"

	MessageCreatedEvent  = "messagesapi.MessageCreated"
	MessageDeletedEvent  = "messagesapi.MessageDeleted"
	MessageRestoredEvent = "messagesapi.MessageRestored"
	MessageUpdatedEvent  = "messagesapi.MessageUpdated"
	MessagesPublishEvent = "messagesapi.MessagesPublished"
	MessagesPrivateEvent = "messagesapi.MessagesPrivated"

	TranslationCreatedEvent = "messagesapi.TranslationCreated"
	TranslationDeletedEvent = "messagesapi.TranslationDeleted"
	TranslationRestoredEvent = "messagesapi.TranslationRestored"
	TranslationUpdatedEvent = "messagesapi.TranslationUpdated"

	CommentCreatedEvent         = "messagesapi.CommentCreated"
	CommentUpdatedEvent         = "messagesapi.CommentUpdated"
	CommentDeletedEvent         = "messagesapi.CommentDeleted"
	MessageCommentsDeletedEvent = "messagesapi.MessageCommentsDeleted"

	DeleteMessageCommand = "messagesapi.DeleteMessage"
	RestoreMessageCommand = "messagesapi.RestoreMessage"
)
