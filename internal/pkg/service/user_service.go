package service

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"fmt"
	"reflect"
	"kdc/internal/pkg/core"
	"github.com/ethereum/go-ethereum/crypto"
	"encoding/hex"
	"strconv"
	"errors"
)

type jsonRpc struct {
	JsonRpc  string `json:"jsonrpc"`
	Method string `json:"method"`
	Id interface{} `json:"id"`
	Params *param `json:"params"`
}

type param struct {
	FileId string `json:"fileId,omitempty"`
	Data string `json:"data,omitempty"`
	Amount *hexutil.Big `json:"amount,omitempty"`
	Signature string `json:"signature"`
}

type jsonResponse struct {
	JsonRpc  string `json:"jsonrpc"`
	Id interface{} `json:"id"`
	Result interface{} `json:"result"`
}

var BadIdErr = errors.New("bad id")

func validJsonRpc2(rpc *jsonRpc) bool{
	// check version
	if rpc.JsonRpc != "2.0" {
		return false
	}
	// check id
	id := rpc.Id
	fmt.Print(reflect.TypeOf(id))
	switch id.(type) {
		case string: break
		case float64: break // go recognise number to float
		default: return false
	}
	// check method
	if rpc.Method != "subtract" && rpc.Method != "read" && rpc.Method != "terminate" {
		return false
	}
	// check param
	if rpc.Params == nil {
		return false
	}
	return true
}

func RunService() {
	// Echo instance
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Routes
	e.POST("/api", handle)
	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

func handle(c echo.Context) (err error) {
	j := new(jsonRpc)
	if err = c.Bind(j); err != nil {
		return
	}
	valid := validJsonRpc2(j)
	if !valid {
		err = echo.NewHTTPError(http.StatusBadRequest, "invalid json rpc 2.0")
		return
	}
	jResponse := new(jsonResponse)
	jResponse.Id = j.Id
	jResponse.JsonRpc = "2.0"
	switch j.Method {
		case "subtract":
			result := handleSubtract(j)
			jResponse.Result = result
		case "read":
			jResponse.Result = 1
		case "terminate":
			jResponse.Result = "terminate result"
		default:
			err = echo.NewHTTPError(http.StatusBadRequest, "method not supported")
			return
	}
	return c.JSON(http.StatusOK, j)
}

func idToStr(id interface{}) (string, error){
	switch id.(type) {
		case string:
			return id.(string), nil
		case float64:  // go recognise number to float
			return strconv.Itoa(int(id.(float64))), nil
		default:
			return "", BadIdErr
	}
}

func handleSubtract(json *jsonRpc) interface{} {
	pp := json.Params
	id := json.Id
	fmt.Print(reflect.TypeOf(id))
	reqId, err := idToStr(id)
	if err != nil {
		return 400
	}
	fileId := pp.FileId
	userId := pp.Data
	amount := pp.Amount
	sig, _ := hex.DecodeString(pp.Signature)
	fmt.Println(amount)
	// compose msg
	msg := json.JsonRpc + json.Method + reqId + fileId + userId + amount.String()
	// sha msg
	shaMsg := crypto.Keccak256([]byte(msg))
	recoveredPub2, _ := crypto.SigToPub(shaMsg, sig)
	// validate
	finalAddr := crypto.PubkeyToAddress(*recoveredPub2).Hex()
	fmt.Printf("calculated addr is %s\n", finalAddr)
	if finalAddr != userId {
		return 400 //invalid signature
	}
	// call core method
	_, err1 := core.SubtractValue(fileId, userId, amount.ToInt())
	if err1 != nil {
		return 500
	}
	return 0
}

func handleRead(json *jsonRpc) interface{} {

	return 0
}

func handleTerminate(json *jsonRpc) interface{} {

	return 0
}