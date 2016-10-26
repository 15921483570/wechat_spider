package main

import (
	"log"
	"net/http"

	"github.com/sundy-li/wechat_spider"

	"github.com/elazarl/goproxy"
)

func main() {
	var port = "8899"
	proxy := goproxy.NewProxyHttpServer()
	//open it see detail logs
	// wechat.Verbose = true
	proxy.OnResponse().DoFunc(
		wechat_spider.ProxyHandle(wechat_spider.NewBaseProcessor()),
	)
	log.Println("server will at port:" + port)
	log.Fatal(http.ListenAndServe(":"+port, proxy))

}
