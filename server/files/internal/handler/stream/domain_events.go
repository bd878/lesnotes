package stream

import (
	"context"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/am"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/files/internal/domain"
	"github.com/bd878/gallery/server/files/pkg/events"
)

type domainHandler[T ddd.Event] struct {
	publisher am.MessagePublisher
}

var _ ddd.EventHandler[ddd.Event] = (*domainHandler[ddd.Event])(nil)

func NewDomainEventHandlers(publisher am.MessagePublisher) *domainHandler[ddd.Event] {
	return &domainHandler[ddd.Event]{publisher}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handler ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handler,
		domain.FileUploadedEvent,
		domain.FilesDeletedEvent,
		domain.FilesPublishedEvent,
		domain.FilesPrivatedEvent,
	)
}

func (h domainHandler[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.FileUploadedEvent:
		return h.onFileUploaded(ctx, event)
	case domain.FilesDeletedEvent:
		return h.onFilesDeleted(ctx, event)
	case domain.FilesPrivatedEvent:
		return h.onFilesPrivated(ctx, event)
	case domain.FilesPublishedEvent:
		return h.onFilesPublished(ctx, event)
	}
	return nil
}

func (h domainHandler[T]) onFileUploaded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.FileUploaded)
	data, err := proto.Marshal(&api.FileUploaded{
		Id:           payload.ID,
		UserId:       payload.UserID,
		Name:         payload.Name,
		Description:  payload.Description,
		Private:      payload.Private,
		Size:         payload.Size,
		Mime:         payload.Mime,
		CreatedAt:    payload.CreatedAt,
		UpdatedAt:    payload.UpdatedAt,
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.FilesChannel, am.NewRawMessage(event.ID(), events.FileUploadedEvent, data))
}

func (h domainHandler[T]) onFilesDeleted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.FilesDeleted)
	data, err := proto.Marshal(&api.FilesDeleted{
		Ids:         payload.IDs,
		UserId:      payload.UserID,
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.FilesChannel, am.NewRawMessage(event.ID(), events.FilesDeletedEvent, data))
}

func (h domainHandler[T]) onFilesPrivated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.FilesPrivated)
	data, err := proto.Marshal(&api.FilesPrivated{
		Ids:         payload.IDs,
		UserId:      payload.UserID,
		UpdatedAt:   payload.UpdatedAt,
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.FilesChannel, am.NewRawMessage(event.ID(), events.FilesPrivatedEvent, data))
}

func (h domainHandler[T]) onFilesPublished(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.FilesPublished)
	data, err := proto.Marshal(&api.FilesPublished{
		Ids:         payload.IDs,
		UserId:      payload.UserID,
		UpdatedAt:   payload.UpdatedAt,
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.FilesChannel, am.NewRawMessage(event.ID(), events.FilesPublishedEvent, data))
}