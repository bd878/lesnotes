package discovery_test

import (
  "testing"
  "fmt"
  "time"
  discovery "github.com/bd878/gallery/server/messages/internal/discovery/serf"
)

func TestMembership(t *testing.T) {
  configs := make([]discovery.Config, 3)
  for i := 0; i < len(configs); i++ {
    addr := fmt.Sprintf("%s:%d", "127.0.0.1", 8000 + i)
    configs[i] = discovery.Config{
      NodeName: fmt.Sprintf("%d", i),
      BindAddr: addr,
      Tags: map[string]string{
        "rpc_addr": addr,
      },
    }
    if i > 0 {
      configs[i].StartJoinAddrs = []string{
        configs[0].BindAddr,
      }
    }
  }

  members := make([]*discovery.Membership, 0)
  handlers := make([]*handler, 0)
  for _, c := range configs {
    h := &handler{
      joins: make(chan map[string]string, 3),
      leaves: make(chan string, 3),
    }
    m, err := discovery.New(h, c)
    if err != nil {
      t.Fatal(err)
    }
    time.Sleep(250*time.Millisecond)

    members = append(members, m)
    handlers = append(handlers, h)
  }

  if len(handlers[0].joins) != 2 {
    t.Errorf("joins != 2, got %d\n", len(handlers[0].joins))
  }

  if len(handlers[0].leaves) != 0 {
    t.Errorf("leaves != 0, got %d\n", len(handlers[0].leaves))
  }

  err := members[0].Leave()
  if err != nil {
    t.Fatal(err)
  }

  if len(handlers[1].leaves) != 1 {
    t.Errorf("leaves != 1, got %d\n", len(handlers[1].leaves))
  }
}

type handler struct {
  joins chan map[string]string
  leaves chan string
}

func (h *handler) Join(id, addr string) error {
  if h.joins != nil {
    h.joins <- map[string]string{
      "id": id,
      "addr": addr,
    }
  }
  return nil
}

func (h *handler) Leave(id string) error {
  if h.leaves != nil {
    h.leaves <- id
  }
  return nil
}