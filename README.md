# wechat_spider
微信公众号爬虫  (只需设置代理, 一键可以爬取指定公众号的所有历史文章)

- 一个简单的Demo  [simple_server.go][1]


```
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
```



- 使用方法:
运行后, 设置手机的代理为 本机ip 8899端口,  打开微信客户端, 点击任一公众号查看历史文章按钮, 即可爬取该公众号的所有历史文章(已经支持自动翻页爬取)


- 自定义输出源,实现Processor接口的Output方法即可, [custom_output_server.go][2]


  [1]: https://github.com/sundy-li/wechat_spider/blob/master/examples/simple_server.go
  [2]: https://github.com/sundy-li/wechat_spider/blob/master/examples/custom_output_server.go
