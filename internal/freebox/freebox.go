package freebox

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dop251/goja"
	"github.com/valyala/fasthttp"
)


func NewFreeboxClient(address string, Dial fasthttp.DialFunc) *Client {
	c := Client{
		http:    &fasthttp.Client{
			Dial: Dial,
		},
		address: address,
		vm:      goja.New(),
		headers: map[string]string{},
	}

	for k, v := range map[string]string{
		`authority`:          address,
		`accept`:             `application/json, text/javascript, */*; q=0.01`,
		`accept-language`:    `fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7`,
		`content-type`:       `application/x-www-form-urlencoded; charset=UTF-8`,
		`origin`:             fmt.Sprintf(`%s`, address),
		`referer`:            fmt.Sprintf(`%s/login.php`, address),
		`sec-ch-ua`:          `"Not.A/Brand";v="8", "Chromium";v="114", "Google Chrome";v="114"`,
		`sec-ch-ua-mobile`:   `?0`,
		`sec-ch-ua-platform`: `"Windows"`,
		`sec-fetch-dest`:     `empty`,
		`sec-fetch-mode`:     `cors`,
		`sec-fetch-site`:     `same-origin`,
		`user-agent`:         `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36`,
		`x-fbx-freebox0s`:    `1`,
		`x-requested-with`:   `XMLHttpRequest`,
	} {
		c.headers[k] = v
	}

	return &c
}

func (c *Client) GetHash(password, salt, challenge string) string {
	hPass := sha1.Sum([]byte(salt + password))
	hPassStr := hex.EncodeToString(hPass[:])

	hash := hmac.New(sha1.New, []byte(hPassStr))
	hash.Write([]byte(challenge))

	return hex.EncodeToString(hash.Sum(nil))
}

func (c *Client) SolveChallenge(challenge *Challenge) string {
	var s string

	for _, value := range challenge.Result.Challenge {
		result, err := c.vm.RunString(value)
		if err != nil {
			panic(err)
		}

		s += fmt.Sprintf("%v", result)
	}
	return s
}

func (c *Client) Login(password string) (*LoginResponse, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(fmt.Sprintf("%s/api/latest/login/", c.address))
	req.Header.SetMethod("POST")

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	req.SetBodyString(fmt.Sprintf("password=%s", password))

	resp := fasthttp.AcquireResponse()
	err := c.http.Do(req, resp)
	if err != nil {
		return nil, err
	}

	var login LoginResponse
	if err := json.Unmarshal(resp.Body(), &login); err != nil {
		return nil, err
	}

	return &login, nil
}

func (c *Client) GetChallenge() (*Challenge, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(fmt.Sprintf("%s/api/latest/login/?_=%v", c.address, time.Now().Unix()))
	req.Header.SetMethod("GET")

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	resp := fasthttp.AcquireResponse()
	err := c.http.Do(req, resp)
	if err != nil {
		return nil, err
	}

	var chall Challenge
	if err := json.Unmarshal(resp.Body(), &chall); err != nil {
		return nil, err
	}

	return &chall, nil
}
