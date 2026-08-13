package service

import (
	"context"
	"github.com/bd878/gallery/server/messages/pkg/model"
)

type messagesControllerTx struct {
	MessagesController
}

func NewMessagesControllerTx(controller MessagesController) messagesControllerTx {
	return messagesControllerTx{
		MessagesController: controller,
	}
}

func (s messagesControllerTx) SaveMessage(ctx context.Context, id int64, text, title string, fileIDs []int64, threadID int64, userID int64, private bool, name string) (message *model.Message, err error) {
	return
}

func (s messagesControllerTx) DeleteUserMessages(ctx context.Context, userID int64) (err error) {
	return 
}

func (s messagesControllerTx) DeleteMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	return
}

func (s messagesControllerTx) PublishMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	return
}

func (s messagesControllerTx) PrivateMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	return
}

func (s messagesControllerTx) UpdateMessage(ctx context.Context, id int64, text, title, name *string, fileIDs []int64, userID int64) (err error) {
	return
}