package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func createCode(w http.ResponseWriter, r *http.Request) {
	gc := &GoogleCode{
		Issuer: r.FormValue("issuer"),
	}
	re := gc.CreateCode()
	b, err := json.Marshal(re)
	if err != nil {
		log.Println("json转化错误")
	}
	w.Write(b)

}
func verifyCode(w http.ResponseWriter, r *http.Request) {
	gc := &GoogleCode{
		Issuer: r.FormValue("issuer"),
		Code:   r.FormValue("code"),
	}
	re := gc.VerifyCode()
	b, err := json.Marshal(re)
	if err != nil {
		log.Println("json转化错误")
	}
	w.Write(b)
}

func main() {
	initClient()
	http.HandleFunc("/createCode", createCode)
	http.HandleFunc("/verifyCode", verifyCode)
	log.Println("google-authenticator启动 端口：8082")
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Panicln("服务启动错误")
		return
	}
	log.Println("google-authenticator启动 端口：8082")

}
