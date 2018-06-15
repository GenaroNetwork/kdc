package core

import (
	"math/big"
	"errors"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	Readwrite = 0
	Readonly  = 1
	Write     = 2
)

type CoinUnitT = big.Int
type AllowTableT = map[string]int
type MortgageTableT = map[string]CoinUnitT
type ModificationT struct {
	operation string
	value CoinUnitT
}

var UnImplementedErr = errors.New("unimplemented") // error usage https://medium.com/@sebdah/go-best-practices-error-handling-2d15e1f0c5ee
var InsufficientBalanceErr = errors.New("insufficient balance")
var NotOwnerErr = errors.New("insufficient privilege: not owner")
var NoPermissionErr = errors.New("user has no permission")
var UnSupportedOperationErr = errors.New("UnSupportedOperationErr")
var NoNegativeValueAllowedErr = errors.New("NoNegativeValueAllowedErr")

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
	permi, _ := getPermissionForFile(userId, fileId)
	if permi == Readwrite || permi == Write {
		// 2. check input
		if amount.Cmp(big.NewInt(0)) == -1 {
			return nil, NoNegativeValueAllowedErr
		}
		// 3. check balance
		bal, err := readValueDirect(fileId, userId)
		if err != nil {
			return nil, err
		}
		if bal.Cmp(amount) == -1 {
			return nil, InsufficientBalanceErr
		}
		// 4. insert modify table
		err = appendNewOperation(fileId, userId, "subtract", hexutil.EncodeBig(amount))
		if err != nil {
			return nil, err
		}
		// 5. return if success
		return bal.Sub(bal, amount), nil
	}
	return nil, NoPermissionErr
}

func singleOperation(operation string, lValue *CoinUnitT, rValue *CoinUnitT) (result *CoinUnitT, err error) {
	switch operation {
	case "init":
		return result.Add(lValue, rValue), nil
	case "subtract":
		return result.Sub(lValue, rValue), nil
	default:
		return nil, UnSupportedOperationErr
	}
}

func calculateAllValue(mods *[]ModificationT) (*CoinUnitT, error) {
	resultVal := big.NewInt(0)
	var err error
	for _, mod := range *mods {
		resultVal, err = singleOperation(mod.operation, resultVal, &mod.value)
		if err!=nil {
			return nil, err
		}
	}
	return resultVal, nil
}

func readValueDirect(fileId string, userId string) (*CoinUnitT, error){
	modifys, err := getOperationsForFile(fileId, userId)
	if err != nil {
		return nil, err
	}
	result, err := calculateAllValue(modifys)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ReadValue(readingUser string, fileId string, userId string) (*CoinUnitT, error) {
	// 1. check privilege
	permi, _ := getPermissionForFile(readingUser, fileId)
	if permi == Readwrite || permi == Readonly {
		// proceed to read
		return readValueDirect(fileId, userId)
	} else {
		return nil, NoPermissionErr
	}
}
