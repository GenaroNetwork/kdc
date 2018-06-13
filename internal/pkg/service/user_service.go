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
	Result interface{} `json:"result,omitempty"`
	Id interface{} `json:"id"`
	Error jsonErr `json:"error,omitempty"`
}

type jsonErr struct {
	Code int `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data interface{} `json:"data,omitempty"`
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
	var jResponse *jsonResponse
	switch j.Method {
		case "subtract":
			jResponse = handleSubtract(j)
			return c.JSON(http.StatusOK, jResponse)
		case "read":
			jResponse = handleRead(j)
			return c.JSON(http.StatusOK, jResponse)
		case "terminate":
			jResponse = handleTerminate(j)
			return c.JSON(http.StatusOK, jResponse)
		default:
			err = echo.NewHTTPError(http.StatusBadRequest, "method not supported")
			return
	}
	err = echo.NewHTTPError(http.StatusInternalServerError, "unknown error")
	return
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

func initJResponse(json *jsonRpc) *jsonResponse {
	jResponse := new(jsonResponse)
	jResponse.Id = json.Id
	jResponse.JsonRpc = "2.0"
	return jResponse
}

func makeJsonError(code int, message string) *jsonErr {
	je := new(jsonErr)
	je.Code = code
	je.Message = message
	return je
}

func handleSubtract(json *jsonRpc) *jsonResponse {
	jResponse := initJResponse(json)
	pp := json.Params
	id := json.Id
	fmt.Print(reflect.TypeOf(id))
	reqId, err := idToStr(id)
	if err != nil {
		jResponse.Error = *makeJsonError(400, err.Error())
		return jResponse
	}
	fileId := pp.FileId
	userId := pp.Data
	amount := pp.Amount
	sig, err1 := hex.DecodeString(pp.Signature)
	if err1 != nil {
		jResponse.Error = *makeJsonError(400, "bad signature")
		return jResponse
	}
	fmt.Println(amount)
	// compose msg
	msg := json.JsonRpc + json.Method + reqId + fileId + userId + amount.String()
	// sha msg
	shaMsg := crypto.Keccak256([]byte(msg))
	recoveredPub2, err3 := crypto.SigToPub(shaMsg, sig)
	if err3 != nil {
		jResponse.Error = *makeJsonError(400, "unable to recover public key")
		return jResponse
	}
	// validate
	finalAddr := crypto.PubkeyToAddress(*recoveredPub2).Hex()
	fmt.Printf("calculated addr is %s\n", finalAddr)
	if finalAddr != userId {
		jResponse.Error = *makeJsonError(400, "invalid signature")
		return jResponse
	}
	// call core method
	_, err2 := core.SubtractValue(userId, fileId, amount.ToInt())
	if err2 != nil {
		jResponse.Error = *makeJsonError(400, err2.Error())
		return jResponse
	}
	jResponse.Result = 1
	return jResponse
}

func handleRead(json *jsonRpc) *jsonResponse {
	jResponse := initJResponse(json)
	pp := json.Params
	id := json.Id
	fmt.Print(reflect.TypeOf(id))
	reqId, err := idToStr(id)
	if err != nil {
		jResponse.Error = *makeJsonError(400, err.Error())
		return jResponse
	}
	fileId := pp.FileId
	userId := pp.Data
	sig, err1 := hex.DecodeString(pp.Signature)
	if err1 != nil {
		jResponse.Error = *makeJsonError(400, "bad signature")
		return jResponse
	}
	// compose msg
	msg := json.JsonRpc + json.Method + reqId + fileId + userId
	// sha msg
	shaMsg := crypto.Keccak256([]byte(msg))
	recoveredPub2, err3 := crypto.SigToPub(shaMsg, sig)
	if err3 != nil {
		jResponse.Error = *makeJsonError(400, "unable to recover public key")
		return jResponse
	}
	// validate
	readingUser := crypto.PubkeyToAddress(*recoveredPub2).Hex()
	fmt.Printf("reading user addr is %s\n", readingUser)
	// call core method
	balance, err2 := core.ReadValue(readingUser, fileId, userId)
	if err2 != nil {
		jResponse.Error = *makeJsonError(400, err2.Error())
		return jResponse
	}
	jResponse.Result = hexutil.EncodeBig(balance)
	return jResponse
}

func handleTerminate(json *jsonRpc) *jsonResponse {
	jResponse := initJResponse(json)
	pp := json.Params
	id := json.Id
	fmt.Print(reflect.TypeOf(id))
	reqId, err := idToStr(id)
	if err != nil {
		jResponse.Error = *makeJsonError(400, err.Error())
		return jResponse
	}
	fileId := pp.FileId
	sig, err1 := hex.DecodeString(pp.Signature)
	if err1 != nil {
		jResponse.Error = *makeJsonError(400, "bad signature")
		return jResponse
	}
	// compose msg
	msg := json.JsonRpc + json.Method + reqId + fileId
	// sha msg
	shaMsg := crypto.Keccak256([]byte(msg))
	recoveredPub2, err3 := crypto.SigToPub(shaMsg, sig)
	if err3 != nil {
		jResponse.Error = *makeJsonError(400, "unable to recover public key")
		return jResponse
	}
	// validate
	readingUser := crypto.PubkeyToAddress(*recoveredPub2).Hex()
	fmt.Printf("reading user addr is %s\n", readingUser)
	// call core method
	err2 := core.Terminate(readingUser, fileId)
	if err2 != nil {
		jResponse.Error = *makeJsonError(400, err2.Error())
		return jResponse
	}
	jResponse.Result = 0
	return jResponse
}