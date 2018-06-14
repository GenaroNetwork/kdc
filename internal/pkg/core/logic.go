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
var NotOwnerErr = errors.New("insufficient privilege: not owner")

func InitFile(userId string, fileId string, allow *AllowTableT, mortgage *MortgageTableT) error{
	err := initNewFile(fileId, userId, "", allow, mortgage)
	return err
}

func Terminate(userId string, fileId string) (string, error){
	// 1. check privilege
	bOwner, _ := isOwner(fileId, userId)
	if !bOwner {
		return "", NotOwnerErr
	}
	// 2. update db.
	err := setFileTerminate(fileId)
	if err != nil {
		return "", err
	}
	// 3. send terminate transaction
	//txHash, err := service.FireSyncTransaction(true, "", nil)
	//if err != nil {
	//	return "", err
	//}
	return "", nil
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
