package repository

import (
  "context"

  "github.com/bd878/gallery/server/pkg/model"
)

type Memory struct {
  records []string
}

func New() *Memory {
  return &Memory{
    records: make([]string, 0),
  }
}

func (m *Memory) Append(_ context.Context, str string) error {
  m.records = append(m.records, str)
  return nil
}

func (m *Memory) GetAll(_ context.Context) ([]model.Message, error) {
  msgs := make([]model.Message, len(m.records))
  for i, v := range m.records {
    msgs[i] = model.Message{Value: v}
  }
  return msgs, nil
}
