package model

import "github.com/bd878/gallery/server/api/files"

func FileToProto(file *File) *files.File {
	return &files.File{
		Id:            file.ID,
		Name:          file.Name,
		CreatedAt:     file.CreatedAt,
		UpdatedAt:     file.UpdatedAt,
		Size:          file.Size,
		Error:         file.Error,
		Private:       file.Private,
		Mime:          file.Mime,
		Oid:           file.OID,
		Description:   file.Description,
	}
}

func FileFromProto(proto *files.File) *File {
	return &File{
		ID:                proto.Id,
		Name:              proto.Name,
		CreatedAt:         proto.CreatedAt,
		UpdatedAt:         proto.UpdatedAt,
		Error:             proto.Error,
		Size:              proto.Size,
		Private:           proto.Private,
		Mime:              proto.Mime,
		OID:               proto.Oid,
		Description:       proto.Description,
	}
}

func MapFilesDictFromProto(mapper (func(*files.File) *File), list map[int64]*files.File) map[int64]*File {
	res := make(map[int64]*File, len(list))
	for id, file := range list {
		res[id] = mapper(file)
	}
	return res
}

func MapFilesDictToProto(mapper (func(*File) *files.File), list map[int64]*File) map[int64]*files.File {
	res := make(map[int64]*files.File, len(list))
	for id, file := range list {
		res[id] = mapper(file)
	}
	return res
}

func MapFilesFromProto(mapper (func(*files.File) *File), files []*files.File) []*File {
	res := make([]*File, len(files))
	for i, file := range files {
		res[i] = mapper(file)
	}
	return res
}

func MapFilesToProto(mapper (func(*File) *files.File), list []*File) []*files.File {
	res := make([]*files.File, len(list))
	for i, file := range list {
		res[i] = mapper(file)
	}
	return res
}