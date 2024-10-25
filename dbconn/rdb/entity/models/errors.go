package rdbmodel

import "errors"

var (
	// ErrNotImpl ...
	ErrNotImpl = errors.New("function no completed implmentation")
	// ErrNilDBConf ...
	ErrNilDBConf = errors.New("nil db config bing input")
	// ErrInitSQLxTX ...
	ErrInitSQLxTX = errors.New("init sqlx.tx driver get error")
	// ErrDBTXnotInit ...
	ErrDBTXnotInit = errors.New("sqlx.tx driver has not been inited yet")
	// ErrTXRollback ...
	ErrTXRolllback = errors.New("sqlx.tx rollback failed")
	// ErrTXCommit ..
	ErrTXCommit = errors.New("sqlx.tx commit get error")
)
