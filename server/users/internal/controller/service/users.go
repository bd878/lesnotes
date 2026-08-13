package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"

	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/api/users"
	sessions "github.com/bd878/gallery/server/db/sessions/pkg/model"
	"github.com/bd878/gallery/server/db/users/pkg/machine"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/di"
	"github.com/bd878/gallery/server/users/internal/controller"
	"github.com/bd878/gallery/server/users/internal/domain"
	"github.com/bd878/gallery/server/users/pkg/model"
)

type SessionsGateway interface {
	GetSession(ctx context.Context, token string) (session *sessions.Session, err error)
	ListUserSessions(ctx context.Context, userID int64) (sessions []*sessions.Session, err error)
	RemoveAllUserSessions(ctx context.Context, userID int64) (err error)
	CreateSession(ctx context.Context, userID int64) (session *sessions.Session, err error)
	RemoveSession(ctx context.Context, token string) (err error)
}

type Controller struct {
	client    users.UsersClient
	sessions  SessionsGateway
	publisher ddd.EventPublisher[ddd.Event]
}

func New(container di.Container, sessions SessionsGateway, publisher ddd.EventPublisher[ddd.Event]) Controller {
	client := container.Get("usersClient").(users.UsersClient)

	return Controller{
		client: client,
		sessions:  sessions,
		publisher: publisher,
	}
}

func (s Controller) CreateUser(ctx context.Context, id int64, login, password string) (user *model.User, err error) {
	slog.Debug("create user",
		slog.Int64("id", id),
		slog.String("login", login),
		slog.Int("len(password)", len(password)),
	)

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	cmd, err := proto.Marshal(&users.AppendCommand{
		Id:             id,
		Login:          login,
		HashedPassword: string(hashed),
		CreatedAt:      time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:      time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return nil, err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.AppendRequest),
		Cmd:      cmd,
		Duration: "10s",
	})
	if err != nil {
		return
	}

	slog.Debug("user created")
	slog.Debug("create session")

	session, err := s.sessions.CreateSession(ctx, id)
	if err != nil {
		return nil, err
	}

	slog.Debug("session created")

	user = &model.User{
		ID:        id,
		Login:     login,
		Token:     session.Token,
		ExpiresAt: session.ExpiresAt,
	}

	return
}

func (s Controller) FindUser(ctx context.Context, id int64, login, token string) (user *model.User, err error) {
	slog.Debug("find user",
		slog.Int64("id", id),
		slog.String("login", login),
		slog.String("token", token),
	)

	var userProto *users.User
	if token != "" {
		session, err := s.sessions.GetSession(ctx, token)
		if err != nil {
			return nil, err
		}

		userProto, err = s.client.GetUser(ctx, &users.GetUserRequest{
			Id: int64(session.UserID),
		})
	} else {
		userProto, err = s.client.FindUser(ctx, &users.FindUserRequest{
			Login: login,
		})
	}
	if err != nil {
		return
	}

	user = model.UserFromProto(userProto)

	return
}

func (s Controller) AuthUser(ctx context.Context, token string) (user *model.User, err error) {
	slog.Debug("auth user", slog.String("token", token))

	session, err := s.sessions.GetSession(ctx, token)
	if err != nil {
		return nil, err
	}

	expiresAt, err := time.Parse(time.RFC3339, session.ExpiresAt)
	if err != nil {
		return nil, err
	}

	if expiresAt.Unix() <= time.Now().Unix() {
		return nil, controller.ErrTokenExpired
	}

	var userProto *users.User
	userProto, err = s.client.GetUser(ctx, &users.GetUserRequest{
		Id: int64(session.UserID),
	})
	if err != nil {
		return
	}

	user = model.UserFromProto(userProto)

	user.Token = session.Token
	user.ExpiresAt = session.ExpiresAt

	return
}

func (s Controller) GetUser(ctx context.Context, id int64) (user *model.User, err error) {
	slog.Debug("get user", slog.Int64("id", id))

	userProto, err := s.client.GetUser(ctx, &users.GetUserRequest{Id: id})
	if err != nil {
		return nil, err
	}

	user = model.UserFromProto(userProto)

	return
}

func (s Controller) UpdateUser(ctx context.Context, id int64, login *string, metadata []byte) (err error) {
	slog.Debug("update user",
		slog.Int64("id", id),
		slog.String("login", *login),
		slog.String("metadata", fmt.Sprintf("%v", metadata)),
	)

	cmd, err := proto.Marshal(&users.UpdateCommand{
		Id:        id,
		Login:     login,
		Metadata:  metadata,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.UpdateRequest),
		Cmd:      cmd,
		Duration: "10s",
	})

	return
}

func (s Controller) MakePremium(ctx context.Context, id int64, invoiceID, createdAt, expiresAt string) (err error) {
	slog.Debug("make premium",
		slog.Int64("id", id),
		slog.String("invoice_id", invoiceID),
		slog.String("created_at", createdAt),
		slog.String("expires_at", expiresAt),
	)

	cmd, err := proto.Marshal(&users.MakePremiumCommand{
		InvoiceId: invoiceID,
		Id:        id,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.MakePremiumRequest),
		Cmd:      cmd,
		Duration: "10s",
	})

	return
}

func (s Controller) LoginUser(ctx context.Context, login, password string) (session *sessions.Session, err error) {
	slog.Debug("login user",
		slog.String("login", login),
		slog.Int("len(password)", len(password)),
	)

	user, err := s.client.FindUser(ctx, &users.FindUserRequest{
		Login: login,
	})
	if err != nil {
		return nil, err
	}

	slog.Debug("user found", slog.Int64("user_id", user.Id))

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		slog.Error(err.Error())
		return nil, controller.ErrWrongPassword
	}

	session, err = s.sessions.CreateSession(ctx, int64(user.Id))

	slog.Debug("session created", slog.String("token", session.Token))

	return
}

func (s Controller) DeleteUser(ctx context.Context, id int64) (err error) {
	slog.Debug("delete user", slog.Int64("id", id))

	event, err := domain.DeleteUser(id)
	if err != nil {
		return err
	}

	err = s.publisher.Publish(ctx, event)
	if err != nil {
		return
	}

	// TODO: emit event, not call
	err = s.sessions.RemoveAllUserSessions(ctx, id)
	if err != nil {
		return
	}

	cmd, err := proto.Marshal(&users.DeleteCommand{
		Id: id,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.DeleteRequest),
		Cmd:      cmd,
		Duration: "10s",
	})

	return 
}

func (s Controller) LogoutUser(ctx context.Context, token string) (err error) {
	slog.Debug("logout user", slog.String("token", token))

	err = s.sessions.RemoveSession(ctx, token)

	return
}
