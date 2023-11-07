package repository

import (
  "context"
  "log"

  "github.com/bd878/gallery/server/pkg/model"
)

type Repository struct {
  records []string
}

func New() *Repository {
  return &Repository{
    records: make([]string, 0),
  }
}

func (m *Repository) Put(_ context.Context, str string) error {
  log.Println("append string =", str)
  m.records = append(m.records, str)
  return nil
}

func (m *Repository) GetAll(_ context.Context) ([]model.Message, error) {
  log.Println("get all records")
  msgs := make([]model.Message, len(m.records))
  for i, v := range m.records {
    msgs[i] = model.Message{Value: v}
  }
  return msgs, nil
}
