package service

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"fmt"
	"reflect"
	"kdc/internal/pkg/core"
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

func handleSubtract(jsonRpc *jsonRpc) interface{} {
	pp := jsonRpc.Params

	fileId := pp.FileId
	userId := pp.Data
	amount := pp.Amount.ToInt()
	fmt.Println(amount)
	// check signature

	// call core method
	core.SubtractValue(fileId, userId, amount)
	return 1
}
