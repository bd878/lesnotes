package distributed

import (
	"time"
	"context"
	"bytes"

	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/sessions/pkg/model"
)

func (m *DistributedSessions) apply(ctx context.Context, reqType RequestType, cmd []byte) (res interface{}, err error) {
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

func (m *DistributedSessions) CreateSession(ctx context.Context, userID int64) (session *model.Session, err error) {
	logger.Debugw("create session", "userID", userID)

	token := utils.RandomString(10)
	expiresUTCNano := time.Now().Add(time.Hour * 24 * 5).UnixNano()

	cmd, err := proto.Marshal(&AppendCommand{
		UserId: userID,
		Token:  token,
		ExpiresUtcNano: expiresUTCNano,
	})
	if err != nil {
		return nil, err
	}

	_, err = m.apply(ctx, AppendRequest, cmd)

	session = &model.Session{
		UserID:         userID,
		Token:          token,
		ExpiresUTCNano: expiresUTCNano,
	}

	return
}

func (m *DistributedSessions) GetSession(ctx context.Context, token string) (*model.Session, error) {
	logger.Debugw("get session", "token", token)

	return m.repo.Get(ctx, token)
}

func (m *DistributedSessions) ListUserSessions(ctx context.Context, userID int64) ([]*model.Session, error) {
	logger.Debugw("list user sessions", "userID", userID)

	return m.repo.List(ctx, userID)
}

func (m *DistributedSessions) RemoveSession(ctx context.Context, token string) error {
	logger.Debugw("remove session", "token", token)

	cmd, err := proto.Marshal(&DeleteCommand{
		Token: token,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, DeleteRequest, cmd)

	return err
}

func (m *DistributedSessions) RemoveUserSessions(ctx context.Context, userID int64) error {
	logger.Debugw("remove user sessions", "userID", userID)

	cmd, err := proto.Marshal(&DeleteUserSessionsCommand{
		UserId: userID,
	})
	if err != nil {
		return err
	}

	_, err = m.apply(ctx, DeleteUserSessionsRequest, cmd)

	return err
}
