package machine

import (
	"context"
	"log/slog"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/api/users"
	"github.com/bd878/gallery/server/db/users/pkg/machine"
)

type UsersRepository interface {
	Save(ctx context.Context, id int64, login, hashedPassword string, metadata []byte, createdAt, updatedAt string) (err error)
	Delete(ctx context.Context, id int64) (err error)
	Update(ctx context.Context, id int64, login *string, metadata []byte, updatedAt string) (err error)
	MakePremium(ctx context.Context, id int64, invoiceID, createdAt, expiresAt string) (err error)
}

type UsersDumper interface {
	Open(ctx context.Context) (ch chan *users.UsersSnapshot, err error)
	Restore(ctx context.Context, user *users.UsersSnapshot) (err error)
	Close() (err error)
}

var _ raft.FSM = (*Machine)(nil)

type Machine struct {
	usersRepo   UsersRepository
	usersDumper UsersDumper
}

func New(usersRepo UsersRepository, usersDumper UsersDumper) *Machine {
	return &Machine{
		usersRepo:   usersRepo,
		usersDumper: usersDumper,
	}
}

func (f *Machine) Apply(record *raft.Log) interface{} {
	slog.Debug("record", slog.Uint64("index", record.Index),
		slog.Uint64("term", record.Term), slog.String("type", record.Type.String()), slog.Time("appended_at", record.AppendedAt))
	buf := record.Data
	reqType := machine.RequestType(buf[0])
	switch reqType {
	case machine.AppendRequest:
		return f.applyAppend(buf[1:])
	case machine.UpdateRequest:
		return f.applyUpdate(buf[1:])
	case machine.DeleteRequest:
		return f.applyDelete(buf[1:])
	case machine.MakePremiumRequest:
		return f.applyMakePremium(buf[1:])
	default:
		slog.Error("unknown request type", slog.Any("type", reqType))
	}
	return nil
}

func (f *Machine) applyAppend(raw []byte) interface{} {
	var cmd users.AppendCommand
	proto.Unmarshal(raw, &cmd)

	err := f.usersRepo.Save(context.TODO(), cmd.Id, cmd.Login, cmd.HashedPassword, cmd.Metadata, cmd.CreatedAt, cmd.UpdatedAt)
	if err != nil {
		slog.Debug("error", slog.String("error", err.Error()))
	}

	return err
}

func (f *Machine) applyUpdate(raw []byte) interface{} {
	var cmd users.UpdateCommand
	proto.Unmarshal(raw, &cmd)

	err := f.usersRepo.Update(context.TODO(), cmd.Id, cmd.Login, cmd.Metadata, cmd.UpdatedAt)
	if err != nil {
		slog.Debug("error", slog.String("error", err.Error()))
	}

	return err
}

func (f *Machine) applyDelete(raw []byte) interface{} {
	var cmd users.DeleteCommand
	proto.Unmarshal(raw, &cmd)

	err := f.usersRepo.Delete(context.TODO(), cmd.Id)
	if err != nil {
		slog.Debug("error", slog.String("error", err.Error()))
	}

	return err
}

func (f *Machine) applyMakePremium(raw []byte) interface{} {
	var cmd users.MakePremiumCommand
	proto.Unmarshal(raw, &cmd)

	err := f.usersRepo.MakePremium(context.TODO(), cmd.Id, cmd.InvoiceId, cmd.CreatedAt, cmd.ExpiresAt)
	if err != nil {
		slog.Debug("error", slog.String("error", err.Error()))
	}

	return err
}