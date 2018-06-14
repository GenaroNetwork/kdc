package service

import (
	"kdc/internal/pkg/core"
	"io/ioutil"
	"net/http"
	"bytes"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type FileIDArr struct {
	FromAccount string 	`json:"fromAccount"`
	Terminate	bool			`json:"terminate"`
	Sidechain	map[string] *hexutil.Big `json:"sidechain"`
}


type SpecialTxInput struct {
	Type    *hexutil.Big    `json:"type"`
	SpecialTxTypeMortgageInit 	FileIDArr	`json:"specialTxTypeMortgageInit"`
}



type SendTxArgs struct {
	From     string  `json:"from"`
	To       string `json:"to"`
	Gas      *hexutil.Uint64 `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Value    *hexutil.Big    `json:"value"`
	Nonce    *hexutil.Uint64 `json:"nonce"`
	// We accept "data" and "input" for backwards-compatibility reasons. "input" is the
	// newer name and should be preferred by clients.
	Data  *hexutil.Bytes `json:"data"`
	ExtraData  SpecialTxInput      `json:"extraData"`
}


func FireSyncTransaction(isTerminate bool, fileId string, mortgage *core.MortgageTableT) (string, error){
	UnlockAccount(SyncAccount,AccountPassword)
	return "0x", nil
}

type UnlockAccountParameter struct {
	Jsonrpc  string 	`json:"jsonrpc"`
	Method 	 string		`json:"method"`
	Params	 []string	`json:"params"`
	Id		 int 		`json:"id"`
}



func UnlockAccount(Account,Password string) bool {
	parameter := UnlockAccountParameter{
		Jsonrpc: "2.0",
		Method: "eth_getMortgageInitByBlockNumberRange",
		Id: 1,
	}
	parameter.Params = append(parameter.Params,Account)
	parameter.Params = append(parameter.Params,Password)
	input,_ := json.Marshal(parameter)
	result := httpPost(input)
	if nil == result {
		return false
	}
	return true
}



func GetInitFile(startNum string)  {
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