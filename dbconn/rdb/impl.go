package rdb

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/hsuanshao/go-tools/ctx"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	rdbIfc "github.com/hsuanshao/go-tools/dbconn/rdb/entity/interface"
	rdbm "github.com/hsuanshao/go-tools/dbconn/rdb/entity/models"
)

var (
	ErrConnectFailed = errors.New("connect to mysql server return error")
	ErrPingFailed    = errors.New("rdb ping get error return")
)

func InitRDB(ctx ctx.CTX, dbConf *rdbm.Connect) (dbconn *sqlx.DB, err error) {
	// NOTE: set this to skip for unit test
	if flag.Lookup("test.v") != nil {
		testing.Short()
	}

	if dbConf == nil {
		ctx.Error("dbConf should  not to be nil pointer")
		return nil, rdbm.ErrNilDBConf
	}
	dbConnectStr := "{user}:{user_password}@{net_protocol}({server_host}:{server_port})/{db_name}"

	dbConnectStr = strings.Replace(dbConnectStr, "{user}", dbConf.User, 1)
	dbConnectStr = strings.Replace(dbConnectStr, "{user_password}", dbConf.Password, 1)
	dbConnectStr = strings.Replace(dbConnectStr, "{net_protocol}", dbConf.Network, 1)
	dbConnectStr = strings.Replace(dbConnectStr, "{server_host}", dbConf.Host, 1)
	dbConnectStr = strings.Replace(dbConnectStr, "{server_port}", dbConf.Port, 1)
	dbConnectStr = strings.Replace(dbConnectStr, "{db_name}", dbConf.DBName, 1)

	conn, err := sqlx.Connect(dbConf.Driver, dbConnectStr)
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err}).Error("connect to database failed")
		return nil, ErrConnectFailed
	}

	err = conn.Ping()
	if err != nil {
		ctx.WithField("err", err).Error("ping failed")
		return nil, ErrPingFailed
	}

	conn.SetConnMaxLifetime(10 * time.Second)
	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(100)

	return conn, nil
}

// PrepareDBTrasaction to init DB Transaction service,
func PrepareDBTrasaction(ctx ctx.CTX, dbConn *sqlx.DB) (txSrv rdbIfc.DBTransaction) {
	return &tximpl{
		dbConn: dbConn,
	}
}

type tximpl struct {
	dbConn   *sqlx.DB
	dbTxConn *sqlx.Tx
}

// GetDBTXDriver to get db sqlx.TX driver, please applied this mehotd to get sqlx.Tx, instead of get from sqlx.DB.Beginx(), it will leads you get unpexpected issue
func (ti *tximpl) GetDBTXDriver(ctx ctx.CTX) (txDriver *sqlx.Tx, err error) {
	if ti.dbTxConn == nil {
		ti.dbTxConn, err = ti.dbConn.Beginx()
		if err != nil {
			ctx.WithFields(logrus.Fields{"err": err}).Error("initail db sqlx.tx failed")
			return nil, rdbm.ErrInitSQLxTX
		}
	}
	return nil, rdbm.ErrNotImpl
}

// FetchTotalBatchStatementAmt to get latest tx batch db transaction db statement command amount
func (ti *tximpl) FetchTotalBatchStatementAmt(ctx ctx.CTX) (count uint16, err error) {
	if ti.dbTxConn == nil {
		ctx.Warn("request fetch tx statment status, but sqlx.tx driver is nil")
		return 0, rdbm.ErrDBTXnotInit
	}

	return 0, rdbm.ErrNotImpl
}

// CommitTX to commit db tx all statements
func (ti *tximpl) CommitTX(ctx ctx.CTX) (err error) {
	if ti.dbTxConn == nil {
		ctx.Warn("request commit tx statment, but sqlx.tx driver is nil")
		return rdbm.ErrDBTXnotInit
	}

	defer func() error {
		// handle some statement caused panic
		// prevent panic from sqlx.Tx original behavior if tx failed (panic error), we need to overwrite it
		if p := recover(); p != nil {
			ctx.WithFields(logrus.Fields{"err": p}).Error("sqlx.tx, db tx.commit panic error")

			switch pt := p.(type) {
			case error:
				err = pt
			default:
				err = fmt.Errorf("%s", pt)
			}
		}
		if err != nil {
			ctx.WithField("err", err).Warn("recover seems get error")
			if rollbackErr := ti.dbTxConn.Rollback(); rollbackErr != nil {
				ctx.WithField("err", rollbackErr).Error("sqlx tx rollback failed")
				return rdbm.ErrTXRolllback
			}
			return err
		}
		return nil
	}()

	err = ti.dbTxConn.Commit()
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err}).Error("sqlx tx commit failed")
		//
		return rdbm.ErrTXCommit
	}

	return
}
