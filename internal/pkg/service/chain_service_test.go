package service

import (
	"testing"
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
