package main

import (
	"context"
	"fmt"
	"flag"
	"time"
	"net"
	"io"
	"os"
	"bytes"
	"database/sql"
	"path/filepath"
	"github.com/soheilhy/cmux"
	"github.com/pressly/goose/v3"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/bd878/gallery/server/db/users/migrations"
	"github.com/bd878/gallery/server/db/users/internal/machine"
	"github.com/bd878/gallery/server/db/users/internal/repository/postgres"
)

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s myNodename myAddr:myPort nodename2 addr2:port2 nodename3 addr3:port3\n", os.Args[0])
	}
}

func main() {
	flag.Parse()

	if flag.NArg() != 6 {
		flag.Usage()
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "myNodename = %s, myAddr:myPort = %s, nodename2 = %s, addr2:port2 = %s, nodename3 = %s, addr3:port3 = %s\n",
		flag.Arg(0), flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4), flag.Arg(5))

	myAddr := flag.Arg(1)
	addrs := [][]string{[]string{flag.Arg(2), flag.Arg(3)}, []string{flag.Arg(4), flag.Arg(5)}}

	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(flag.Arg(0))
	config.LogLevel = hclog.Debug.String()

	var configuration raft.Configuration
	configuration.Servers = []raft.Server{{
		ID:      raft.ServerID(flag.Arg(0)),
		Address: raft.ServerAddress(myAddr),
	}}
	for _, pair := range addrs {
		nodename, addr := pair[0], pair[1]
		configuration.Servers = append(configuration.Servers, raft.Server{
			ID:      raft.ServerID(nodename),
			Address: raft.ServerAddress(addr),
			Suffrage: raft.Voter,
		})
	}

	// transport
	listener, err := net.Listen("tcp4", myAddr)
	if err != nil {
		panic(err)
	}

	mux := cmux.New(listener)

	raftListener := mux.Match(func(r io.Reader) bool {
		b := make([]byte, 1)
		if _, err := r.Read(b); err != nil {
			return false
		}
		return bytes.Compare(b, []byte{byte(RaftRPC)}) == 0
	})

	transport := raft.NewNetworkTransport(NewStreamLayer(raftListener), 5, 10 * time.Second, os.Stderr)
	// END transport

	// fsm

	pool, err := pgxpool.New(context.TODO(), os.Getenv("PG_CONN"))
	if err != nil {
		panic(err)
	}

	usersRepo := postgres.NewUsersRepository(pool, "users.users", "users.premiums")
	usersDumper := postgres.NewUsersDumper(pool, "users.users", "users.premiums")

	fsm := machine.New(usersRepo, usersDumper)

	// END fsm

	// migrate

	goose.SetVerbose(true)
	goose.SetTableName("users_goose_db_version")

	db, err := sql.Open("pgx", os.Getenv("PG_CONN"))
	if err != nil {
		panic(err)
	}

	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	if err := goose.Up(db, "."); err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err := goose.Reset(db, ".")
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to reset db", err)
		}

		err = db.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to close db", err)
		}
	}(db)

	// END migrate

	// store

	logStore, err := raftboltdb.NewBoltStore(filepath.Join("./data.users/raft", "log"))
	if err != nil {
		panic(err)
	}
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join("./data.users/raft", "stable"))
	if err != nil {
		panic(err)
	}

	snapshotStore, err := raft.NewFileSnapshotStore(filepath.Join("./data.users/raft", "snapshot"), 1, os.Stderr)
	if err != nil {
		panic(err)
	}

	// END store

	err = raft.RecoverCluster(config, fsm, logStore, stableStore, snapshotStore, transport, configuration)
	if err != nil {
		panic(err)
	}
}