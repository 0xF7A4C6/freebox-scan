package freebox

import (
	"github.com/dop251/goja"
	"github.com/valyala/fasthttp"
)

type Client struct {
	http    *fasthttp.Client
	address string
	headers map[string]string
	vm      *goja.Runtime
}

type Challenge struct {
	Success bool `json:"success"`
	Result  struct {
		LoggedIn     bool     `json:"logged_in"`
		Challenge    []string `json:"challenge"`
		PasswordSalt string   `json:"password_salt"`
		PasswordSet  bool     `json:"password_set"`
	} `json:"result"`
}

type LoginResponse struct {
	UID     string `json:"uid"`
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Result  struct {
		PasswordSalt string   `json:"password_salt"`
		Challenge    []string `json:"challenge"`
	} `json:"result"`
	ErrorCode string `json:"error_code"`
}
