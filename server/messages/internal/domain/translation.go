package domain

import (
	"errors"
	"github.com/bd878/gallery/server/internal/ddd"
)

const (
	TranslationCreatedEvent  = "messages.TranslationCreated"
	TranslationDeletedEvent  = "messages.TranslationDeleted"
	TranslationUpdatedEvent  = "messages.TranslationUpdated"
)

var (
	ErrLangRequired = errors.New("lang is empty")
)

type TranslationCreated struct {
	MessageID     int64
	UserID        int64
	Lang          string
	Text          string
	Title         string
}

func (TranslationCreated) Key() string { return TranslationCreatedEvent }

func CreateTranslation(userID, messageID int64, lang string, title, text string) (ddd.Event, error) {
	if messageID == 0 {
		return nil, ErrIDRequired
	}
	if userID == 0 {
		return nil, ErrIDRequired
	}
	if lang == "" {
		return nil, ErrLangRequired
	}

	return ddd.NewEvent(TranslationCreatedEvent, &TranslationCreated{
		MessageID: messageID,
		UserID:    userID,
		Lang:      lang,
		Text:      text,
		Title:     title,
	}), nil
}

type TranslationDeleted struct {
	MessageID   int64
	Lang        string
}

func (TranslationDeleted) Key() string { return TranslationDeletedEvent }

func DeleteTranslation(messageID int64, lang string) (ddd.Event, error) {
	return ddd.NewEvent(TranslationDeletedEvent, &TranslationDeleted{
		MessageID:     messageID,
		Lang:          lang,
	}), nil
}

type TranslationUpdated struct {
	MessageID    int64
	Lang         string
	Text         *string
	Title        *string
}

func (TranslationUpdated) Key() string { return TranslationUpdatedEvent }

func UpdateTranslation(messageID int64, lang string, title, text *string) (ddd.Event, error) {
	if messageID == 0 {
		return nil, ErrIDRequired
	}
	if lang == "" {
		return nil, ErrLangRequired
	}

	return ddd.NewEvent(TranslationUpdatedEvent, &TranslationUpdated{
		MessageID:   messageID,
		Lang:        lang,
		Text:        text,
		Title:       title,
	}), nil
}