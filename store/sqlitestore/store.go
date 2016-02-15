// Package sqlitestore implements the store interfaces in order to1
// persist data to an sqlite database.
package sqlitestore

import (
	"github.com/unprofession-al/gerty/store"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // import is blank since registration to sql/sqlx is done via init func
)

var (
	nodes NodeStore
	roles RoleStore
)

func init() {
	store.Register("sqlitestore", Setup)
}

// Setup creates the DB if it does not exist yet. A connectien is established.
// The configuration string must contain the file path of the database file.
func Setup(config string) (*store.Store, error) {
	db, err := sqlx.Open("sqlite3", config)
	if err != nil {
		return nil, err
	}

	db.MustExec(nodeSchema)
	db.MustExec(roleSchema)

	nodes = NodeStore{db: db}
	roles = RoleStore{db: db}
	s := &store.Store{
		Nodes: nodes,
		Roles: roles,
	}
	return s, nil
}
