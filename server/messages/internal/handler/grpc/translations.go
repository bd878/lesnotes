package grpc

import (
	"context"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/messages/pkg/model"
)

type TranslationsController interface {
	SaveTranslation(ctx context.Context, messageID int64, lang, title, text string) (err error)
	UpdateTranslation(ctx context.Context, messageID int64, lang string, title, text *string) (err error)
	DeleteTranslation(ctx context.Context, messageID int64, lang string) (err error)
	ReadTranslation(ctx context.Context, messageID int64, lang string) (result *model.Translation, err error)
}

type TranslationsHandler struct {
	api.UnimplementedTranslationsServer
	controller TranslationsController
}

func NewTranslationsHandler(ctrl TranslationsController) *TranslationsHandler {
	handler := &TranslationsHandler{controller: ctrl}

	return handler
}

func (h *TranslationsHandler) SaveTranslation(ctx context.Context, req *api.SaveTranslationRequest) (resp *api.SaveTranslationResponse, err error) {
	err = h.controller.SaveTranslation(ctx, req.MessageId, req.Lang, req.Title, req.Text)

	resp = &api.SaveTranslationResponse{}

	return
}

func (h *TranslationsHandler) UpdateTranslation(ctx context.Context, req *api.UpdateTranslationRequest) (resp *api.UpdateTranslationResponse, err error) {
	err = h.controller.UpdateTranslation(ctx, req.MessageId, req.Lang, req.Title, req.Text)

	resp = &api.UpdateTranslationResponse{}

	return
}

func (h *TranslationsHandler) DeleteTranslation(ctx context.Context, req *api.DeleteTranslationRequest) (resp *api.DeleteTranslationResponse, err error) {
	err = h.controller.DeleteTranslation(ctx, req.MessageId, req.Lang)

	resp = &api.DeleteTranslationResponse{}

	return
}

func (h *TranslationsHandler) ReadTranslation(ctx context.Context, req *api.ReadTranslationRequest) (resp *api.ReadTranslationResponse, err error) {
	translation, err := h.controller.ReadTranslation(ctx, req.MessageId, req.Lang)
	if err != nil {
		return nil, err
	}

	resp = &api.ReadTranslationResponse{}
	resp.Translation = model.TranslationToProto(translation)

	return
}