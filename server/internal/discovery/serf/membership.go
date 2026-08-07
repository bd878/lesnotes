package serf

import (
	"context"
	"fmt"
	"net"
	"os"
	"log/slog"

	"github.com/hashicorp/serf/serf"
)

type Membership struct {
	Config
	handler Handler
	serf *serf.Serf
	events chan serf.Event
}

type Config struct {
	NodeName       string
	BindAddr       string
	Tags           map[string]string
	SerfJoinAddrs  []string
}

func New(config Config, handler Handler) (*Membership, error) {
	c := &Membership{
		Config: config,
		handler: handler,
	}
	if err := c.setupSerf(); err != nil {
		return nil, err
	}
	return c, nil
}

func (m *Membership) setupSerf() error {
	addr, err := net.ResolveTCPAddr("tcp", m.BindAddr)
	if err != nil {
		return err
	}
	config := serf.DefaultConfig()
	config.Init()
	config.MemberlistConfig.BindAddr = addr.IP.String()
	config.MemberlistConfig.BindPort = addr.Port
	m.events = make(chan serf.Event)
	config.EventCh = m.events
	config.Tags = m.Tags
	config.NodeName = m.Config.NodeName

	m.serf, err = serf.Create(config)
	if err != nil {
		return err
	}

	if m.SerfJoinAddrs != nil {
		_, err = m.serf.Join(m.SerfJoinAddrs, true)
		if err != nil {
			return err
		}
	}

	return nil
}

type Handler interface {
	Join(name, addr string) error
	Leave(name string) error
	Snapshot() error
	Restore() error
	ShowLeader() error
}

func (m *Membership) Run(ctx context.Context) {
	defer fmt.Fprintln(os.Stdout, "leaving membership")
	fmt.Fprintf(os.Stdout, "membership started %s\n", m.BindAddr)
	for {
		select {
		case <-ctx.Done():
			if err := m.Leave(); err != nil {
				slog.Error("membership failed to leave", slog.String("error", err.Error()))
			}
			return
		case e := <- m.events:
			slog.Debug("serf new event", slog.String("type", e.EventType().String()))
			switch e.EventType() {
			case serf.EventMemberJoin:
				for _, member := range e.(serf.MemberEvent).Members {
					if m.isLocal(member) {
						continue
					}
					m.handleJoin(member)
				}

			case serf.EventMemberLeave, serf.EventMemberFailed:
				for _, member := range e.(serf.MemberEvent).Members {
					if m.isLocal(member) {
						return
					}
					m.handleLeave(member)
				}

			case serf.EventQuery:
				switch e.String() {
				case "query: snapshot":
					slog.Debug("performing snapshot")
					err := m.handler.Snapshot()
					if err != nil {
						slog.Debug("snapshot returned error", slog.String("error", err.Error()))
					}
					slog.Debug("snapshot finished")

				case "query: restore":
					slog.Debug("performing restore")
					err := m.handler.Restore()
					if err != nil {
						slog.Debug("restore returned error", slog.String("error", err.Error()))
					}
					slog.Debug("restore finished")

				case "query: leader":
					slog.Debug("who is leader")
					err := m.handler.ShowLeader()
					if err != nil {
						slog.Debug("failed to show roles", slog.String("error", err.Error()))
					}
					slog.Debug("show leader finished")

				default:
					slog.Error("unknown event payload", slog.String("payload", e.String()))
				}

			default:
				slog.Warn("unknown event", slog.String("event", e.String()))
			}
		}
	}
}

func (m *Membership) isLocal(member serf.Member) bool {
	return m.serf.LocalMember().Name == member.Name
}

func (m *Membership) Members() []serf.Member {
	return m.serf.Members()
}

func (m *Membership) Leave() error {
	return m.serf.Leave()
}

func (m *Membership) handleJoin(member serf.Member) {
	m.handler.Join(member.Name, member.Tags["raft_addr"])
}

func (m *Membership) handleLeave(member serf.Member) {
	m.handler.Leave(member.Name)
}