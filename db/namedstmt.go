package db

import "database/sql"

type NamedStmt struct {
	Params      []string
	QueryString string
	Stmt        *sql.Stmt
}
