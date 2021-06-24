package main

import (
	"fmt"
	"testing"
)

func TestAddNewAddresses(t *testing.T) {
	LoadConf()
	fmt.Println("")
	InitDB(GlobalConfig.DbConfig.DbType, GlobalConfig.DbConfig.DbSource)
	GlobalDBMgr.TblAddressMgr.AddNewAddresses([]string{"13K4uYefwJ19t4NgYDgRyHfQfnwh5qULka",
		"14K4uYefwJ19t4NgYDgRyHfQfnwh5qULka"})
}

func TestListAddrUtxos(t *testing.T) {
	LoadConf()
	fmt.Println("")
	InitDB(GlobalConfig.DbConfig.DbType, GlobalConfig.DbConfig.DbSource)
	utxos, _ := GlobalDBMgr.TblUtxoMgr.ListAddrUtxos("13K4uYefwJ19t4NgYDgRyHfQfnwh5qULka")
	fmt.Println("utxos:", utxos)
}
