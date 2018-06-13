package core

import (
	"math/big"
	"errors"
)

const (
	Readwrite = 0
	Readonly  = 1
	Write     = 2
)

type CoinUnitT = big.Int
type AllowTableT = map[string]int
type MortgageTableT = map[string]CoinUnitT

var UnImplementedErr = errors.New("unimplemented") // error usage https://medium.com/@sebdah/go-best-practices-error-handling-2d15e1f0c5ee
var InsufficientBalanceErr = errors.New("insufficient balance")

func InitFile(userId string, fileId string, allow *AllowTableT, mortgage *MortgageTableT) error{
	// insert things into db.
	return UnImplementedErr
}



func Terminate(userId string, fileId string) error{
	// 1. update db.
	// 2. send terminate transaction
	return UnImplementedErr
}

func SubtractValue(userId string, fileId string, amount *CoinUnitT) (*CoinUnitT, error) {
	// 1. check privilege
	// 2. insert modify table
	// 3. update state
	// 4. return if success
	return nil, UnImplementedErr
}

func ReadValue(readingUser string, fileId string, userId string) (*CoinUnitT, error) {
	// 1. check privilege
	// 2. read state
	return nil, UnImplementedErr
}
