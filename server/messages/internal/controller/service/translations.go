package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/messages/pkg/model"
	"github.com/bd878/gallery/server/messages/pkg/loadbalance"
)

type TranslationsConfig struct {
	RpcAddr string
}

type TranslationsController struct {
	conf       TranslationsConfig
	client     api.TranslationsClient
	conn       *grpc.ClientConn
}

func NewTranslationsController(conf TranslationsConfig) *TranslationsController {
	controller := &TranslationsController{
		conf: conf,
	}

	controller.setupConnection()

	return controller
}

func (s *TranslationsController) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *TranslationsController) setupConnection() (err error) {
	conn, err := grpc.NewClient(
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

	client := api.NewTranslationsClient(conn)

	s.conn = conn
	s.client = client

	return
}

func (s *TranslationsController) isConnFailed() bool {
	state := s.conn.GetState()
	if state == connectivity.Shutdown || state == connectivity.TransientFailure {
		logger.Debugln("connection failed")
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

	_, err = s.client.SaveTranslation(ctx, &api.SaveTranslationRequest{
		MessageId:    messageID,
		UserId:       userID,
		Lang:         lang,
		Title:        title,
		Text:         text,
	})

	return
}

func (s *TranslationsController) UpdateTranslation(ctx context.Context, messageID int64, lang string, title, text *string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("update translation", "message_id", messageID, "lang", lang, "title", title, "text", text)

	_, err = s.client.UpdateTranslation(ctx, &api.UpdateTranslationRequest{
		MessageId:   messageID,
		Lang:        lang,
		Title:       title,
		Text:        text,
	})

	return
}

func (s *TranslationsController) DeleteTranslation(ctx context.Context, messageID int64, lang string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("delete translation", "message_id", messageID, "lang", lang)

	_, err = s.client.DeleteTranslation(ctx, &api.DeleteTranslationRequest{
		MessageId:      messageID,
		Lang:           lang,
	})

	return
}

func (s *TranslationsController) ReadTranslation(ctx context.Context, userID, messageID int64, lang string, name *string) (translation *model.Translation, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("read translation", "user_id", userID, "message_id", messageID, "lang", lang, "name", name)

	resp, err := s.client.ReadTranslation(ctx, &api.ReadTranslationRequest{
		UserId:       userID,
		MessageId:    messageID,
		Lang:         lang,
		Name:         name,
	})
	if err != nil {
		return nil, err
	}

	translation = model.TranslationFromProto(resp.Translation)

	return
}

func (s *TranslationsController) ListTranslations(ctx context.Context, userID, messageID int64, name string) (translations []*model.Translation, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("list translations", "user_id", userID, "message_id", messageID, "name", name)

	resp, err := s.client.ListTranslations(ctx, &api.ListTranslationsRequest{
		UserId:    userID,
		MessageId: messageID,
		Name:      name,
	})
	if err != nil {
		return nil, err
	}

	translations = model.MapTranslationsFromProto(model.TranslationFromProto, resp.Translations)

	return
}