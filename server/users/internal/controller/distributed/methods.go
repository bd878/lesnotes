package distributed

import (
	"time"
	"context"
	"bytes"
	"golang.org/x/crypto/bcrypt"

	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/users/pkg/model"
)

func (m *DistributedUsers) apply(ctx context.Context, reqType RequestType, cmd []byte) (res interface{}, err error) {
	var buf bytes.Buffer
	_, err = buf.Write([]byte{byte(reqType)})
	if err != nil {
		return
	}

	_, err = buf.Write(cmd)
	if err != nil {
		return
	}

	timeout := 10*time.Second
	/* fsm.Apply() */
	future := m.raft.Apply(buf.Bytes(), timeout)
	if future.Error() != nil {
		return nil, future.Error()
	}

	res = future.Response()
	if err, ok := res.(error); ok {
		return nil, err
	}

	return
}

func (m *DistributedUsers) CreateUser(ctx context.Context, id int64, login, password string) (user *model.User, err error) {
	logger.Debugw("create user", "id", id, "login", login, "password", password)


	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	cmd, err := proto.Marshal(&AppendCommand{
		Id:             id,
		Login:          login,
		HashedPassword: string(hashed),
		FontSize:       0,
		Theme:          "light",
		Lang:           "",
	})
	if err != nil {
		return nil, err
	}

	_, err = m.apply(ctx, AppendRequest, cmd)

	user = &model.User{
		ID:           id,
		Login:        login,
		Theme:        "light",
		Lang:         "",
		HashedPassword: string(hashed),
		FontSize:      0,
	}

	return
}

func (m *DistributedUsers) DeleteUser(ctx context.Context, id int64) (err error) {
	logger.Debugw("delete user", "id", id)

	cmd, err := proto.Marshal(&DeleteCommand{
		Id: id,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, DeleteRequest, cmd)
	if err != nil {
		return
	}

	return nil
}

func (m *DistributedUsers) UpdateUser(ctx context.Context, id int64, login, theme, lang string, fontSize int32) (err error) {
	logger.Debugw("update user", "id", id, "login", login, "theme", theme, "lang", lang, "font_size", fontSize)

	cmd, err := proto.Marshal(&UpdateCommand{
		Id:    id,
		Login: login,
		Lang:  lang,
		Theme: theme,
		FontSize: fontSize,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, UpdateRequest, cmd)

	return
}

func (m *DistributedUsers) FindUser(ctx context.Context, login string) (user *model.User, err error) {
	logger.Debugw("find user", "login", login)
	return m.repo.Find(ctx, 0, login)
}

func (m *DistributedUsers) GetUser(ctx context.Context, id int64) (user *model.User, err error) {
	logger.Debugw("get user", "id", id)
	return m.repo.Find(ctx, id, "")
}