package rdbifc

import (
	"github.com/hsuanshao/go-tools/ctx"
	"github.com/jmoiron/sqlx"
)

// DBTransaction is based on sqlx.TX to provide a method that can bypass DB level process (aka repo), can as parameter of Business method param, and could be understand from method input/output
type DBTransaction interface {
	// GetDBTXDriver to get db sqlx.TX driver, please applied this mehotd to get sqlx.Tx, instead of get from sqlx.DB.Beginx(), it will leads you get unpexpected issue
	GetDBTXDriver(ctx ctx.CTX) (txDriver *sqlx.Tx, err error)

	// FetchTotalBatchStatementAmt to get latest tx batch db transaction db statement command amount
	FetchTotalBatchStatementAmt(ctx ctx.CTX) (count uint16, err error)

	// CommitTX to commit db tx all statements
	CommitTX(ctx ctx.CTX) (err error)
}
