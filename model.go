package main

import (
	"errors"
	"sync"
	"time"
)

type address struct {
	Id         int       `xorm:"pk INTEGER autoincr"`
	Address    string    `xorm:"VARCHAR(128) NOT NULL"`
	Extra      int       `xorm:"INT NULL"`
	Created_at time.Time `xorm:"created"`
	Updated_at time.Time `xorm:"DATETIME"`
}

type tblAddressMgr struct {
	TableName string
	Mutex     *sync.Mutex
}

func (t *tblAddressMgr) Init() {
	t.TableName = "address"
	t.Mutex = new(sync.Mutex)
}

func (t *tblAddressMgr) AddNewAddresses(addrs []string) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	for _, addr := range addrs {
		var addressRes address
		count, err := GetDBEngine().Where("address=?", addr).Count(addressRes)
		if err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		addressRes.Address = addr
		_, err = GetDBEngine().Cols("address").InsertOne(addressRes)
		if err != nil {
			return err
		}
	}
	return nil
}

type utxo struct {
	Id           int       `xorm:"pk INTEGER autoincr"`
	Txid         string    `xorm:"VARCHAR(128) NOT NULL"`
	Vout         int       `xorm:"INT NOT NULL"`
	Amount       string    `xorm:"VARCHAR(128) NOT NULL"`
	Used         int       `xorm:"INT NOT NULL"`
	Address      string    `xorm:"VARCHAR(128) NOT NULL"`
	Scriptpubkey string    `xorm:"VARCHAR(128) NOT NULL"`
	Coin_symbol  string    `xorm:"VARCHAR(128) NOT NULL"`
	Created_at   time.Time `xorm:"created"`
	Updated_at   time.Time `xorm:"DATETIME"`
	Pending      int       `xorm:"INT NOT NULL"`
}

type tblUtxoMgr struct {
	TableName string
	Mutex     *sync.Mutex
}

func (t *tblUtxoMgr) Init() {
	t.TableName = "utxo"
	t.Mutex = new(sync.Mutex)
}

func (t *tblUtxoMgr) ListAddrUtxos(addr string) ([]utxo, error) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	utxos := make([]utxo, 0)
	err := GetDBEngine().Cols("*").Where("address=? and used=0 and pending=0", addr).Find(&utxos)
	if err != nil {
		return utxos, err
	}
	return utxos, nil
}

func (t *tblUtxoMgr) UpdateUtxoPendingState(txId string, vout int, pending int) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	var u utxo
	u.Txid = txId
	u.Vout = vout
	exist, err := GetDBEngine().Where("txid=?", txId).And("vout=?", vout).Get(&u)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("key not found!")
	}

	u.Pending = pending
	_, err = GetDBEngine().Where("txid=?", txId).And("vout=?", vout).Cols("pending").Update(&u)
	return err
}
