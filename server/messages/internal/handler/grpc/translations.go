package grpc

import (
	"context"
	"time"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/messages/internal/machine"
)

type TranslationsController interface {
	Apply(ctx context.Context, reqType machine.RequestType, cmd []byte, duration time.Duration) (err error)
	ReadTranslation(ctx context.Context, userID, messageID int64, lang string, name string) (result *api.Translation, err error)
	ListTranslations(ctx context.Context, userID, messageID int64, name string) (result []*api.Translation, err error)
}

type TranslationsHandler struct {
	api.UnimplementedTranslationsServer
	controller TranslationsController
}

func NewTranslationsHandler(ctrl TranslationsController) *TranslationsHandler {
	handler := &TranslationsHandler{controller: ctrl}

	return handler
}

func (h *TranslationsHandler) Apply(ctx context.Context, req *api.Command) (resp *api.CommandResponse, err error) {
	duration, err := time.ParseDuration(req.Duration)
	if err != nil {
		return nil, err
	}

	err = h.controller.Apply(ctx, machine.RequestType(req.ReqType), req.Cmd, duration)

	resp = &api.CommandResponse{}

	return
}

func (h *TranslationsHandler) ReadTranslation(ctx context.Context, req *api.ReadTranslationRequest) (resp *api.ReadTranslationResponse, err error) {
	var name string
	if req.Name != nil {
		name = *req.Name
	}

	translation, err := h.controller.ReadTranslation(ctx, req.UserId, req.Id, req.Lang, name)
	if err != nil {
		return nil, err
	}

	resp = &api.ReadTranslationResponse{
		Translation: translation,
	}

	return
}

func (h *TranslationsHandler) ListTranslations(ctx context.Context, req *api.ListTranslationsRequest) (resp *api.ListTranslationsResponse, err error) {
	translations, err := h.controller.ListTranslations(ctx, req.UserId, req.Id, req.Name)
	if err != nil {
		return nil, err
	}

	resp = &api.ListTranslationsResponse{Translations: translations}

	return
}
