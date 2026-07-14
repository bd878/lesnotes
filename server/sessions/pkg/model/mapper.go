package model

import "github.com/bd878/gallery/server/api/sessions"

func SessionFromProto(proto *sessions.Session) *Session {
	return &Session{
		UserID:             proto.UserId,
		Token:              proto.Token,
		ExpiresAt:          proto.ExpiresAt,
		CreatedAt:          proto.CreatedAt,
	}
}

func SessionToProto(session *Session) *sessions.Session {
	return &sessions.Session{
		UserId:             session.UserID,
		Token:              session.Token,
		ExpiresAt:          session.ExpiresAt,
		CreatedAt:          session.CreatedAt,
	}
}

func MapSessionsToProto(mapper (func(*Session) *sessions.Session), list []*Session) []*sessions.Session {
	res := make([]*sessions.Session, len(list))
	for i, session := range list {
		res[i] = mapper(session)
	}
	return res
}

func MapSessionsFromProto(mapper (func(*sessions.Session) *Session), list []*sessions.Session) []*Session {
	res := make([]*Session, len(list))
	for i, session := range list {
		res[i] = mapper(session)
	}
	return res
}