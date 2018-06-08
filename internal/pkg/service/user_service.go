package service

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"fmt"
	"reflect"
)

type jsonRpc2 struct {
	Jsonrpc  string `json:"jsonrpc"`
	Method string `json:"method"`
	Id interface{} `json:"id"`
	Params *writeParam `json:"params"`
}

type writeParam struct {
	FileId string `json:"fileId"`
	Data string `json:"data"`
	Amount *hexutil.Big `json:"amount"`
	Signature string `json:"signature"`
}

type jsonResponse struct {
	Jsonrpc  string `json:"jsonrpc"`
	Id interface{} `json:"id"`
	Result interface{} `json:"result"`
}

func validJsonRpc2(rpc *jsonRpc2) bool{
	// check version
	if rpc.Jsonrpc != "2.0" {
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
	e.GET("/", hello)
	e.POST("/api", handle)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}


func handle(c echo.Context) (err error) {
	j := new(jsonRpc2)
	if err = c.Bind(j); err != nil {
		return
	}
	valid := validJsonRpc2(j)
	if !valid {
		err = echo.NewHTTPError(http.StatusBadRequest, "invalid json rpc 2.0")
		return
	}
	jresponse := new(jsonResponse)
	jresponse.Id = j.Id
	jresponse.Jsonrpc = "2.0"
	switch j.Method {
		case "subtract":
			jresponse.Result = "subtract result"
		case "read":
			jresponse.Result = 1
		case "terminate":
			jresponse.Result = "terminate result"
		default:
			err = echo.NewHTTPError(http.StatusBadRequest, "method not supported")
			return
	}
	return c.JSON(http.StatusOK, jresponse)
}

func handleSubtract(jsonrpc *jsonRpc2) interface{} {
	// check signature

	// call core method
	return 1
}
