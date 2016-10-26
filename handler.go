package wechat_spider

import (
	"log"
	"net/http"
	"os"
	"reflect"
	"spiderx/utils"
	"strings"

	"github.com/elazarl/goproxy"
)

var (
	Verbose = false
	Logger  = log.New(os.Stderr, "", log.LstdFlags)
)

func ProxyHandle(proc Processor) func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	return func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if ctx.Req.URL.Path == `/mp/getmasssendmsg` && !strings.Contains(ctx.Req.URL.RawQuery, `f=json`) {
			var data []byte
			var err error
			data, resp.Body, err = utils.CopyReader(resp.Body)
			if err != nil {
				return resp
			}
			t := reflect.TypeOf(proc)
			v := reflect.New(t.Elem())
			p := v.Interface().(Processor)
			go func() {
				err = p.Process(ctx.Req, data)
				if err != nil {
					Logger.Println(err.Error())
				}
				p.Output()
			}()
		}
		return resp
	}

}
