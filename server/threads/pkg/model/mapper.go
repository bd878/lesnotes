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
		CreatedAt:      proto.CreatedAt,
		UpdatedAt:      proto.UpdatedAt,
		Title:          proto.Title,
		PrivateMessage: proto.PrivateMessage,
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
		CreatedAt:      msg.CreatedAt,
		UpdatedAt:      msg.UpdatedAt,
		Title:          msg.Title,
		PrivateMessage: msg.PrivateMessage,
	}
}

func PathFromProto(path []*api.PathStep) []*PathStep {
	res := make([]*PathStep, len(path))
	for i, step := range path {
		res[i] = &PathStep{
			ID: step.Id,
			Name: step.Name,
			Private: step.Private,
		}
	}
	return res
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