package model

import (
	"github.com/bd878/gallery/server/api"
)

func SessionFromProto(proto *api.Session) *Session {
	return &Session{
		UserID:             proto.UserId,
		Token:              proto.Token,
		ExpiresUTCNano:     proto.ExpiresUtcNano,
	}
}

func SessionToProto(session *Session) *api.Session {
	return &api.Session{
		UserId:             session.UserID,
		Token:              session.Token,
		ExpiresUtcNano:     session.ExpiresUTCNano,
	}
}

func MapSessionsToProto(mapper (func(*Session) *api.Session), sessions []*Session) []*api.Session {
	res := make([]*api.Session, len(sessions))
	for i, session := range sessions {
		res[i] = mapper(session)
	}
	return res
}

func MapSessionsFromProto(mapper (func(*api.Session) *Session), sessions []*api.Session) []*Session {
	res := make([]*Session, len(sessions))
	for i, session := range sessions {
		res[i] = mapper(session)
	}
	return res
}