package sqlitestore

import (
	"github.com/unprofession-al/gerty/store"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	nodes NodeStore
	roles RoleStore
)

func init() {
	store.Register("sqlitestore", Setup)
}

func Setup(config string) (*store.Store, error) {
	db, err := sqlx.Open("sqlite3", config)
	if err != nil {
		return nil, err
	}

	db.MustExec(node_schema)
	db.MustExec(role_schema)

	nodes = NodeStore{db: db}
	roles = RoleStore{db: db}
	s := &store.Store{
		Nodes: nodes,
		Roles: roles,
	}
	return s, nil
}
