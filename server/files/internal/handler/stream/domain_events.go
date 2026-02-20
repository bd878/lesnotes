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
	publisher am.MessagePublisher[am.Message]
}

var _ ddd.EventHandler[ddd.Event] = (*domainHandler[ddd.Event])(nil)

func NewDomainEventHandlers(publisher am.MessagePublisher[am.Message]) *domainHandler[ddd.Event] {
	return &domainHandler[ddd.Event]{publisher}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handler ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handler,
		domain.FileUploadedEvent,
		domain.FileDeletedEvent,
		domain.FilePublishedEvent,
		domain.FilePrivatedEvent,
	)
}

func (h domainHandler[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.FileUploadedEvent:
		return h.onFileUploaded(ctx, event)
	case domain.FileDeletedEvent:
		return h.onFileDeleted(ctx, event)
	case domain.FilePrivatedEvent:
		return h.onFilePrivated(ctx, event)
	case domain.FilePublishedEvent:
		return h.onFilePublished(ctx, event)
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
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.FilesChannel, am.NewRawMessage(event.ID(), events.FileUploadedEvent, data))
}

func (h domainHandler[T]) onFileDeleted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.FileDeleted)
	data, err := proto.Marshal(&api.FileDeleted{
		Id:     payload.ID,
		UserId: payload.UserID, 
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.FilesChannel, am.NewRawMessage(event.ID(), events.FileDeletedEvent, data))
}

func (h domainHandler[T]) onFilePrivated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.FilePrivated)
	data, err := proto.Marshal(&api.FilePrivated{
		Id:      payload.ID,
		UserId:  payload.UserID,
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.FilesChannel, am.NewRawMessage(event.ID(), events.FilePrivatedEvent, data))
}

func (h domainHandler[T]) onFilePublished(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.FilePublished)
	data, err := proto.Marshal(&api.FilePublished{
		Id:       payload.ID,
		UserId:   payload.UserID,
	})
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, events.FilesChannel, am.NewRawMessage(event.ID(), events.FilePublishedEvent, data))
}