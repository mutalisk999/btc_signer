package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
)

var GlobalError map[int]string

type Err struct {
	ErrCode int    `json:"code"`
	ErrMsg  string `json:"message"`
}

func SetInternalError(ctx iris.Context, errorStr string) {
	ctx.Values().Set("error", errorStr)
	ctx.StatusCode(iris.StatusInternalServerError)
}

func FormatSysError(err error) *Err {
	sysError := new(Err)
	if err == nil {
		sysError.ErrCode = 0
		sysError.ErrMsg = ""
	} else {
		sysError.ErrCode = -1
		sysError.ErrMsg = err.Error()
	}
	return sysError
}

func GetErrorString(err *Err) string {
	if err == nil {
		return "no error"
	}
	return fmt.Sprintf("errcode: %d, error msg: %s", err.ErrCode, err.ErrMsg)
}

func MakeError(errCode int, errMsg string) *Err {
	err := new(Err)
	err.ErrCode = errCode
	err.ErrMsg = errMsg

	return err
}
