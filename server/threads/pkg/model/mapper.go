package model

import (
	"github.com/bd878/gallery/server/api"
)

func ThreadFromProto(proto *api.Thread) *Thread {
	return &Thread{
		ID:             proto.Id,
		UserID:         proto.UserId,
		ParentID:       proto.ParentId,
		Private:        proto.Private,
		Name:           proto.Name,
		Count:          proto.Count,
		NextID:         proto.NextId,
		PrevID:         proto.PrevId,
		Description:    proto.Description,
	}
}

func ThreadToProto(msg *Thread) *api.Thread {
	return &api.Thread{
		Id:             msg.ID,
		UserId:         msg.UserID,
		ParentId:       msg.ParentID,
		Private:        msg.Private,
		Count:          msg.Count,
		Name:           msg.Name,
		NextId:         msg.NextID,
		PrevId:         msg.PrevID,
		Description:    msg.Description,
	}
}

func MapThreadsToProto(mapper (func(*Thread) *api.Thread), msgs []*Thread) []*api.Thread {
	res := make([]*api.Thread, len(msgs))
	for i, msg := range msgs {
		res[i] = mapper(msg)
	}
	return res
}

func MapThreadsFromProto(mapper (func(*api.Thread) *Thread), msgs []*api.Thread) []*Thread {
	res := make([]*Thread, len(msgs))
	for i, msg := range msgs {
		res[i] = mapper(msg)
	}
	return res
}