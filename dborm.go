package main

import (
	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

type DBMgr struct {
	DBEngine      *xorm.Engine
	TblAddressMgr *tblAddressMgr
	TblUtxoMgr    *tblUtxoMgr
}

var GlobalDBMgr *DBMgr

func GetDBEngine() *xorm.Engine {
	return GlobalDBMgr.DBEngine
}

func InitDB(dbType string, dbSource string) error {
	var err error
	GlobalDBMgr = new(DBMgr)
	GlobalDBMgr.DBEngine, err = xorm.NewEngine(dbType, dbSource)
	if err != nil {
		return err
	}
	GlobalDBMgr.DBEngine.SetTableMapper(core.SnakeMapper{})
	GlobalDBMgr.DBEngine.SetColumnMapper(core.SnakeMapper{})

	GlobalDBMgr.TblAddressMgr = new(tblAddressMgr)
	GlobalDBMgr.TblAddressMgr.Init()

	GlobalDBMgr.TblUtxoMgr = new(tblUtxoMgr)
	GlobalDBMgr.TblUtxoMgr.Init()

	return nil
}
