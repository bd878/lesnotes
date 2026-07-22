package service

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/api/translations"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/rpc"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/db/messages/pkg/machine"
	"github.com/bd878/gallery/server/db/messages/pkg/loadbalance"
	"github.com/bd878/gallery/server/messages/pkg/model"
	"github.com/bd878/gallery/server/messages/internal/domain"
)

type TranslationsConfig struct {
	RpcAddr string
}

type TranslationsController struct {
	conf   TranslationsConfig
	client translations.TranslationsClient
	conn   *grpc.ClientConn
	publisher    ddd.EventPublisher[ddd.Event]
}

func NewTranslationsController(conf TranslationsConfig, publisher ddd.EventPublisher[ddd.Event]) *TranslationsController {
	controller := &TranslationsController{
		conf: conf,
		publisher: publisher,
	}

	controller.setupConnection()

	return controller
}

func (s *TranslationsController) Close() {
	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			logger.Error(zap.Error(err))
		}
	}
}

func (s *TranslationsController) setupConnection() (err error) {
	s.Close()

	conn, err := rpc.NewClient(
		fmt.Sprintf(
			"%s:///%s",
			loadbalance.Name,
			s.conf.RpcAddr,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	client := translations.NewTranslationsClient(conn)

	s.conn = conn
	s.client = client

	return
}

func (s *TranslationsController) isConnFailed() bool {
	state := s.conn.GetState()
	if state == connectivity.Shutdown ||
		state == connectivity.TransientFailure ||
		state == connectivity.Connecting {
		logger.Debugw("translations conn failed", "state", state.String())
		return true
	}
	return false
}

func (s *TranslationsController) SaveTranslation(ctx context.Context, userID, messageID int64, lang, title, text string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("save translation", "user_id", userID, "message_id", messageID, "lang", lang, "title", title, "text", text)

	createdAt := time.Now().UTC().Format(time.RFC3339)
	updatedAt := time.Now().UTC().Format(time.RFC3339)

	cmd, err := proto.Marshal(&translations.AppendTranslationCommand{
		MessageId: messageID,
		Lang:      lang,
		Title:     title,
		Text:      text,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.AppendTranslationRequest),
		Cmd: cmd,
		Duration: "10s",
	})
	if err != nil {
		return err
	}

	event, err := domain.CreateTranslation(userID, messageID, lang, title, text, createdAt, updatedAt)
	if err != nil {
		return err
	}

	return s.publisher.Publish(ctx, event)
}

func (s *TranslationsController) UpdateTranslation(ctx context.Context, messageID int64, lang string, title, text *string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("update translation", "message_id", messageID, "lang", lang, "title", title, "text", text)

	updatedAt := time.Now().UTC().Format(time.RFC3339)

	cmd, err := proto.Marshal(&translations.UpdateTranslationCommand{
		MessageId: messageID,
		Lang:      lang,
		Title:     title,
		Text:      text,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.UpdateTranslationRequest),
		Cmd: cmd,
		Duration: "10s",
	})
	if err != nil {
		return err
	}

	event, err := domain.UpdateTranslation(messageID, lang, title, text, updatedAt)
	if err != nil {
		return err
	}

	return s.publisher.Publish(ctx, event)
}

func (s *TranslationsController) DeleteTranslation(ctx context.Context, messageID int64, lang string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("delete translation", "message_id", messageID, "lang", lang)

	cmd, err := proto.Marshal(&translations.DeleteTranslationCommand{
		MessageId: messageID,
		Lang:      lang,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.DeleteTranslationRequest),
		Cmd: cmd,
		Duration: "10s",
	})
	if err != nil {
		return err
	}

	event, err := domain.DeleteTranslation(messageID, lang)
	if err != nil {
		return err
	}

	return s.publisher.Publish(ctx, event)
}

func (s *TranslationsController) ReadTranslation(ctx context.Context, userID, messageID int64, lang string, name *string) (translation *model.Translation, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("read translation", "user_id", userID, "message_id", messageID, "lang", lang, "name", name)

	resp, err := s.client.ReadTranslation(ctx, &translations.ReadTranslationRequest{
		UserId:    userID,
		Id:        messageID,
		Lang:      lang,
		Name:      name,
	})
	if err != nil {
		return nil, err
	}

	translation = model.TranslationFromProto(resp.Translation)

	return
}

func (s *TranslationsController) ListTranslations(ctx context.Context, userID, messageID int64, name string) (list []*model.Translation, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("list translations", "user_id", userID, "message_id", messageID, "name", name)

	resp, err := s.client.ListTranslations(ctx, &translations.ListTranslationsRequest{
		UserId:    userID,
		Id:        messageID,
		Name:      name,
	})
	if err != nil {
		return nil, err
	}

	list = model.MapTranslationsFromProto(model.TranslationFromProto, resp.Translations)

	return
}
