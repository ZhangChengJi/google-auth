package main

import (
	"context"
	"fmt"
	qrcode "github.com/skip2/go-qrcode"
	"google-auth/returncode"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type GoogleCode struct {
	Issuer string `json:"issuer,omitempty"`
	Code   string `json:"code,omitempty"`
}
type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

const kvSuffix = "google_authenticator_"

var (
	codes   = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	codeLen = len(codes)
)

func RandNewStr(len int) string {
	data := make([]byte, len)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < len; i++ {
		idx := rand.Intn(codeLen)
		data[i] = byte(codes[idx])
	}

	return string(data)
}
func (g *GoogleCode) CreateCode() *Result {
	ctx := context.Background()
	secretId := RandNewStr(32)
	if v := rdb.Get(ctx, kvSuffix+g.Issuer).Val(); v == "" {
		rdb.Set(ctx, kvSuffix+g.Issuer, secretId, 5000*time.Hour)
	} else {
		log.Printf("%v 用户已经注册,秘钥 %v", g.Issuer, v)
		return &Result{
			-1,
			fmt.Sprintf("%v 用户已经注册,秘钥 %v", g.Issuer, v),
		}
	}
	url := "otpauth://totp/zabbix?secret=" + secretId + "&issuer=" + g.Issuer
	error := qrcode.WriteFile(url, qrcode.Medium, 256, "./image/"+g.Issuer+".png")
	if error != nil {
		fmt.Println("write error")
	}
	return &Result{
		1,
		fmt.Sprintf("%v 用户注册成功", g.Issuer),
	}

}
func (g *GoogleCode) VerifyCode() *Result {
	ctx := context.Background()
	if v := rdb.Get(ctx, kvSuffix+g.Issuer).Val(); v == "" {
		return &Result{
			-1,
			fmt.Sprintf("用户%v未在Google Authenticator注册", g.Issuer),
		}
	} else {
		re := returncode.ReturnCode(v)
		if g.Code == strconv.Itoa(int(re)) {
			return &Result{
				0,
				fmt.Sprintf("用户:%v Google验证码验证成功", g.Issuer),
			}
		} else {
			return &Result{
				1,
				"Google验证码错误",
			}
		}
	}
}
