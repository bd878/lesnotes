package pkg

const (
	ThreadsChannel = "gallery.threads.Thread"

	CommandChannel = "gallery.threads.commands"

	ThreadCreatedEvent   = "threadsapi.ThreadCreated"
	ThreadDeletedEvent   = "threadsapi.ThreadDeleted"
	ThreadUpdatedEvent   = "threadsapi.ThreadUpdated"
	ThreadPublishedEvent   = "threadsapi.ThreadPublished"
	ThreadPrivatedEvent   = "threadsapi.ThreadPrivated"
	ThreadParentChangedEvent  = "threadsapi.ThreadParentChangedEvent"

	CreateThreadCommand = "threadsapi.CreateCommand"
)
