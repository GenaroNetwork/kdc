package service

import (
	"testing"
	"kdc/internal/pkg/core"
	"fmt"
)

func TestGetBlockNumber(t *testing.T) {
	GetBlockNumber()
}

func TestGetMortgageInitByBlockNumberRange(t *testing.T) {
	GetMortgageInitByBlockNumberRange("0x0")
}

func TestGetInitFile(t *testing.T) {
	GetInitFile("0x0")
}

func TestUnlockAccount(t *testing.T) {
	UnlockAccount(SyncAccount,AccountPassword)
}

func TestFireSyncTransaction(t *testing.T)  {
	mortgageTable := make(core.Mortgage)
	mortgageTable["0xc4f27fe3af8b76c3120c7e8b27beac0b2927abf2"] = "0x22222244"
	mortgageTable["0x92fb6a50a6817d19b1cb47bdc55a687add4ea21a"] = "0x777777"
	FireSyncTransaction(true,"0xaf7a12de8dc1de25c0541966695498074f52a1cc","c2bb8976c35037f73b594425b5ee77f6931bae7e3c6fd91cec52a570a804e6f8",&mortgageTable)
}


/*
`{
    "jsonrpc": "2.0",
    "method": "eth_getLogSwitchByAddressAndFileID",
    "params": ["{\"0xaf7a12de8dc1de25c0541966695498074f52a1cc\":[\"cb88764da463d1c4497e159fd6668852d451864e677d4fa1c1a6fdd958f2aed0\",\"092a11e499aa4f731576552829410b2b0e1e562f65a1307d47ff06d5a0212a47\",\"092a11e499aa4f731576552829410b2b0e1e562f65a1307d47ff06d5a0212a47\"],\"0xaf7a12de8dc1de25c0541966695498074f52a1cx\":[\"0628fefe0ce2f126e99d2aa98c0220d4092f4ed12b4f4c0d98beb24b984cfc2a\",\"092a11e499aa4f731576552829410b2b0e1e562f65a1307d47ff06d5a0212a47\",\"092a11e499aa4f731576552829410b2b0e1e562f65a1307d47ff06d5a0212a47\"]}"],
    "id": 1
}`*/
func TestGetLogSwitchByAddressAndFileID(t *testing.T)  {
	addressAndFileID := make(map[string]FileID)
	var fileIDArr FileID
	fileIDArr = append(fileIDArr,"cb88764da463d1c4497e159fd6668852d451864e677d4fa1c1a6fdd958f2aed0")
	fileIDArr = append(fileIDArr,"cb88764da463d1c4497e159fd6668852d451864e677d4fa1c1a6fdd958f2aed0")
	addressAndFileID["0xaf7a12de8dc1de25c0541966695498074f52a1cc"] = fileIDArr
	resultArr := GetLogSwitchByAddressAndFileID(addressAndFileID)
	fmt.Println(resultArr)
}