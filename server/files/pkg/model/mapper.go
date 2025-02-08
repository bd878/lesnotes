package model

import (
  "github.com/bd878/gallery/server/api"
)

func FileToProto(file *File) *api.File {
  return &api.File{
    Id:            file.ID,
    Name:          file.Name,
    CreateUtcNano: file.CreateUTCNano,
    Error:         file.Error,
  }
}

func FileFromProto(proto *api.File) *File {
  return &File{
    ID:                proto.Id,
    Name:              proto.Name,
    CreateUTCNano:     proto.CreateUtcNano,
    Error:             proto.Error,
  }
}

func MapFilesFromProto(mapper (func(*api.File) *File), files []*api.File) []*File {
  res := make([]*File, len(files))
  for i, msg := range files {
    res[i] = mapper(msg)
  }
  return res
}

func MapFilesToProto(mapper (func(*File) *api.File), files []*File) []*api.File {
  res := make([]*api.File, len(files))
  for i, msg := range files {
    res[i] = mapper(msg)
  }
  return res
}