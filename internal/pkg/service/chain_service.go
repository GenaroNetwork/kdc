package service

import (
	"kdc/internal/pkg/core"
	"io/ioutil"
	"net/http"
	"bytes"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type MortgageTab struct {
	FromAccount string 	`json:"fromAccount"`
	Terminate	bool			`json:"terminate"`
	Sidechain	*core.Mortgage `json:"sidechain"`
	FileID		string			`json:"fileID"`
}


type SpecialTxInput struct {
	Type    string    `json:"type"`
	SpecialTxTypeMortgageInit 	MortgageTab	`json:"specialTxTypeMortgageInit"`
}



type SendTxArgs struct {
	From     string  `json:"from"`
	To       string `json:"to"`
	Gas      string `json:"gas"`
	GasPrice string    `json:"gasPrice"`
	Value    string    `json:"value"`
	Data  *hexutil.Bytes `json:"data"`
	ExtraData  string      `json:"extraData"`
}


type FireSyncTransactionParameter struct {
	Jsonrpc  string 	`json:"jsonrpc"`
	Method 	 string		`json:"method"`
	Params	 []SendTxArgs	`json:"params"`
	Id		 int 		`json:"id"`
}

func FireSyncTransaction(isTerminate bool, fromAccount,fileId string, mortgage *core.Mortgage) bool {
	if"" == fileId || nil == mortgage {
		return false
	}
	unlock := UnlockAccount(SyncAccount,AccountPassword)
	if false == unlock {
		return false
	}
	mortgageTab := MortgageTab{
		FromAccount:fromAccount,
		Terminate:isTerminate,
		Sidechain:mortgage,
		FileID:fileId,
	}
	txInput := SpecialTxInput{
		Type:SyncTransactionType,
		SpecialTxTypeMortgageInit:mortgageTab,
	}
	sendTxArgs  := SendTxArgs{
		From:SyncAccount,
		To:SpecialAccount,
		Gas:GasVal,
		GasPrice:GasPriceVal,
	}
	extraData,_ := json.Marshal(txInput)
	sendTxArgs.ExtraData = string(extraData)
	parameter := FireSyncTransactionParameter{
		Jsonrpc: "2.0",
		Method: "eth_sendTransaction",
		Id: 1,
	}
	parameter.Params=append(parameter.Params,sendTxArgs)
	input,_ := json.Marshal(parameter)
	result := httpPost(input)
	if nil == result {
		return false
	}
	return true
}

type UnlockAccountParameter struct {
	Jsonrpc  string 	`json:"jsonrpc"`
	Method 	 string		`json:"method"`
	Params	 []string	`json:"params"`
	Id		 int 		`json:"id"`
}



func UnlockAccount(account,password string) bool {
	if "" == account || "" == password {
		return false
	}
	parameter := UnlockAccountParameter{
		Jsonrpc: "2.0",
		Method: "eth_getMortgageInitByBlockNumberRange",
		Id: 1,
	}
	parameter.Params = append(parameter.Params,account)
	parameter.Params = append(parameter.Params,password)
	input,_ := json.Marshal(parameter)
	result := httpPost(input)
	if nil == result {
		return false
	}
	return true
}



func GetInitFile(startNum string)  {
	if "" == startNum {
		return
	}
	mortgageInitResultArr := GetMortgageInitByBlockNumberRange(startNum)
	var AllowTableArr core.AllowTableT
	MortgageTableArr := make(core.MortgageTableT)
	for _,v := range mortgageInitResultArr{
		AllowTableArr = v.AuthorityTable
		for k,v := range v.MortgageTable {
			MortgageTableArr[k] = v.ToInt()
		}
		core.InitFile(v.FromAccount,v.FileID,&AllowTableArr,&MortgageTableArr,v.CreateTime,v.EndTime)
	}
}


type InitFile struct {
	MortgageTable	map[string] *hexutil.Big	`json:"mortgage"`
	AuthorityTable 	map[string]int	`json:"authority"`
	FileID 			string		`json:"fileID"`
	CreateTime  int64	`json:"createTime"`
	EndTime  int64	`json:"endTime"`
	FromAccount string 	`json:"fromAccount"`
}

type MortgageInitResult struct {
	Id 			int			`json:"id"`
	Jsonrpc		string		`json:"jsonrpc"`
	Result		[]InitFile		`json:"result"`
}

type MortgageInitParameter struct {
	Jsonrpc  string 	`json:"jsonrpc"`
	Method 	 string		`json:"method"`
	Params	 []string	`json:"params"`
	Id		 int 		`json:"id"`
}

func GetMortgageInitByBlockNumberRange(startNum string) []InitFile{
	if "" == startNum {
		return nil
	}
	parameter := MortgageInitParameter{
		Jsonrpc: "2.0",
		Method: "eth_getMortgageInitByBlockNumberRange",
		Id: 1,
	}
	parameter.Params = append(parameter.Params,startNum)
	endNum := GetBlockNumber()
	if "0x0" == endNum {
		return	nil
	}
	parameter.Params = append(parameter.Params,endNum)
	input,_ := json.Marshal(parameter)
	result := httpPost(input)
	if nil == result {
		return nil
	}
	var mortgageInitResultArr  MortgageInitResult
	json.Unmarshal(result, &mortgageInitResultArr)
	if nil != mortgageInitResultArr.Result {
		return mortgageInitResultArr.Result
	}
	return nil
}

type BlockNumberResult struct {
	Id 			int			`json:"id"`
	Jsonrpc		string		`json:"jsonrpc"`
	Result		string		`json:"result"`
}


func GetBlockNumber() string {
	var blockNumberResult BlockNumberResult
	result := httpPost([]byte(`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":83}`))
	if nil == result {
		return "0x0"
	}
	json.Unmarshal(result, &blockNumberResult)
	return blockNumberResult.Result
}

func httpPost(parameter []byte) []byte {
	if nil == parameter {
		return nil
	}
	client := &http.Client{}
	req_parameter := bytes.NewBuffer(parameter)
	request, _ := http.NewRequest("POST", ServeUrl, req_parameter)
	request.Header.Set("Content-type", "application/json")
	response, _ := client.Do(request)
	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return body
	}
	return nil
}