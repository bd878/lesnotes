package domain

import (
	"errors"
	"github.com/bd878/gallery/server/internal/ddd"
)

const (
	UserDeletedEvent = "users.UserDeleted"
)

var (
	ErrIDRequired = errors.New("id is 0")
)

type UserDeleted struct {
	UserID int64
}

func (UserDeleted) Key() string { return UserDeletedEvent }

func DeleteUser(userID int64) (ddd.Event, error) {
	if userID == 0 {
		return nil, ErrIDRequired
	}

	return ddd.NewEvent(UserDeletedEvent, &UserDeleted{
		UserID: userID,
	}), nil
}