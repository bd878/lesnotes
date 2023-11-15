package repository

import (
  "context"
  "log"

  "github.com/bd878/gallery/server/pkg/model"
)

type Repository struct {
  records [](*model.Message)
}

func New() *Repository {
  return &Repository{
    records: make([]string, 0),
  }
}

func (m *Repository) Put(_ context.Context, msg *model.Message) error {
  m.records = append(m.records, msg)
  return nil
}

func (m *Repository) GetAll(_ context.Context) ([]model.Message, error) {
  msgs := make([]model.Message, len(m.records))
  for i, m := range m.records {
    msgs[i] = *m
  }
  return msgs, nil
}

// adds user
// func (m *Repository) AddUser(_ context.Context, usr *model.User) error {}
