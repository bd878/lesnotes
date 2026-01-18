package model

import (
	"github.com/bd878/gallery/server/api"
)

func FileToProto(file *File) *api.File {
	return &api.File{
		Id:            file.ID,
		Name:          file.Name,
		CreateUtcNano: file.CreateUTCNano,
		Size:          file.Size,
		Error:         file.Error,
		Private:       file.Private,
		Mime:          file.Mime,
		Description:   file.Description,
	}
}

func FileFromProto(proto *api.File) *File {
	return &File{
		ID:                proto.Id,
		Name:              proto.Name,
		CreateUTCNano:     proto.CreateUtcNano,
		Error:             proto.Error,
		Size:              proto.Size,
		Private:           proto.Private,
		Mime:              proto.Mime,
		Description:       proto.Description,
	}
}

func MapFilesDictFromProto(mapper (func(*api.File) *File), files map[int64]*api.File) map[int64]*File {
	res := make(map[int64]*File, len(files))
	for id, file := range files {
		res[id] = mapper(file)
	}
	return res
}

func MapFilesDictToProto(mapper (func(*File) *api.File), files map[int64]*File) map[int64]*api.File {
	res := make(map[int64]*api.File, len(files))
	for id, file := range files {
		res[id] = mapper(file)
	}
	return res
}

func MapFilesFromProto(mapper (func(*api.File) *File), files []*api.File) []*File {
	res := make([]*File, len(files))
	for i, file := range files {
		res[i] = mapper(file)
	}
	return res
}

func MapFilesToProto(mapper (func(*File) *api.File), files []*File) []*api.File {
	res := make([]*api.File, len(files))
	for i, file := range files {
		res[i] = mapper(file)
	}
	return res
}