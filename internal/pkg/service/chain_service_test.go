package service

import (
	"testing"
	"kdc/internal/pkg/core"
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