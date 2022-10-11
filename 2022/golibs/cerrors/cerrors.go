package cerrors

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Errors struct {
	Code    int    `json:"code" example:"4000000"`
	Message string `json:"message"`
}

var ErrorCodeList = map[string]Errors{
	"parameter_error":     {4000000, "参数错误[ %s ]"},
	"unauthorized_access": {4000001, "未授权访问"},
	"forbidden_access":    {4000003, "禁止访问"},

	"internal_service_error":   {5000000, "服务内部错误"},
	"server_unavailable_error": {5000003, "服务暂不可用"},
	"api_undefined_error":      {5000101, "未定义错误[ %s ]"},
}

func Error(key string) error {
	ret, _ := json.Marshal(ErrorCodeList[key])
	return errors.New(string(ret))
}

func Errorf(key string, msg interface{}) error {
	e := ErrorCodeList[key]
	e.SetMessage(fmt.Sprintf(e.GetMessage(), msg))
	ret, _ := json.Marshal(e)
	return errors.New(string(ret))
}

func (e *Errors) Error() string {
	ret, _ := json.Marshal(e)
	return string(ret)
}

func (e *Errors) SetCode(code int) {
	e.Code = code
}

func (e *Errors) SetMessage(message string) {
	e.Message = message
}

func (e *Errors) GetCode() int {
	return e.Code
}

func (e *Errors) GetMessage() string {
	return e.Message
}
