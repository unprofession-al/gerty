// Package sqlitestore implements the store interfaces in order to1
// persist data to an sqlite database.
package sqlitestore

import (
	"database/sql"

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

	db.MustExec(roleSchema)
	db.MustExec(nodeSchema)
	db.MustExec(nodeRoleSchema)
	db.MustExec(dbSettings)

	nodes = NodeStore{db: db}
	roles = RoleStore{db: db}
	s := &store.Store{
		Nodes: nodes,
		Roles: roles,
	}
	return s, nil
}

var dbSettings = `PRAGMA foreign_keys = on;`

// Role defines the structure of the database entries.
type Role struct {
	ID         int           `db:"id"`
	Name       string        `db:"name"`
	Vars       string        `db:"vars"`
	Parent     sql.NullInt64 `db:"parent"`
	ParentName string        `db:"parent_name"`
}

var roleSchema = `CREATE TABLE IF NOT EXISTS role (
	id     INTEGER PRIMARY KEY AUTOINCREMENT,
	name   TEXT UNIQUE NOT NULL,
	vars   TEXT,
	parent INTEGER REFERENCES role(id) ON DELETE RESTRICT
);`

// Node defines the structure of the database entries.
type Node struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Vars string `db:"vars"`
}

var nodeSchema = `CREATE TABLE IF NOT EXISTS node (
	id     INTEGER PRIMARY KEY AUTOINCREMENT,
	name   TEXT UNIQUE NOT NULL,
	vars   TEXT
);`

// NodeRole represents the n-m mapping of nodes and roles

var nodeRoleSchema = `CREATE TABLE IF NOT EXISTS node_role (
	id      STRING UNIQUE NOT NULL,
	node_id INTEGER NOT NULL REFERENCES node(id) ON DELETE CASCADE,
 	role_id INTEGER NOT NULL REFERENCES role(id) ON DELETE RESTRICT
);`
