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
		wechat_spider.ProxyHandle(&CustomProcessor{}),
	)
	log.Println("server will at port:" + port)
	log.Fatal(http.ListenAndServe(":"+port, proxy))

}

//Just to implement Output Method of interface{} Processor
type CustomProcessor struct {
	wechat_spider.BaseProcessor
}

func (c *CustomProcessor) Output() {
	//Just print the length of result urls
	println("result urls size =>", len(c.Urls()))
	c.Urls()
}
