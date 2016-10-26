package wechat_spider

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"time"

	"net/http"

	"github.com/palantir/stacktrace"
)

type Processor interface {
	//Core method
	Process(req *http.Request, data []byte) error
	//Result urls
	Urls() []string
	//Output
	Output()
	//Sleep method to avoid the req control of wechat
	Sleep()
}

type BaseProcessor struct {
	req    *http.Request
	lastId string
	data   []byte
	result []string
}

var (
	replacer = strings.NewReplacer(
		"\t", "", " ", "",
		"&quot;", `"`, "&nbsp;", "",
		`\\`, "", "&amp;amp;", "&",
		"&amp;", "&", `\`, "",
	)

	urlRegex    = regexp.MustCompile("http://mp.weixin.qq.com/s?[^#]*")
	idRegex     = regexp.MustCompile(`"id":(\d+)`)
	MsgNotFound = errors.New("MsgLists not found")
)

func NewBaseProcessor() *BaseProcessor {
	return &BaseProcessor{}
}

func (p *BaseProcessor) init(req *http.Request, data []byte) (err error) {
	p.req = req
	p.data = data
	fmt.Println("Running a new wechat processor, please wait...")
	return nil
}
func (p *BaseProcessor) Process(req *http.Request, data []byte) error {
	if err := p.init(req, data); err != nil {
		return err
	}

	if err := p.processMain(); err != nil {
		return err
	}
	if err := p.processPages(); err != nil {
		return err
	}
	return nil
}

func (p *BaseProcessor) Sleep() {
	time.Sleep(50 * time.Millisecond)
}

func (p *BaseProcessor) Urls() []string {
	return p.result
}

func (p *BaseProcessor) Output() {
	bs, _ := json.Marshal(p.Urls())
	fmt.Println("result => ", string(bs))
}

//Parse the html
func (p *BaseProcessor) processMain() error {
	p.result = make([]string, 0, 100)
	buffer := bytes.NewBuffer(p.data)
	var msgs string
	str, err := buffer.ReadString('\n')
	for err == nil {
		if strings.Contains(str, "msgList = ") {
			msgs = str
			break
		}
		str, err = buffer.ReadString('\n')
	}
	if msgs == "" {
		return stacktrace.Propagate(MsgNotFound, "Failed parse main")
	}
	msgs = replacer.Replace(msgs)
	p.result = urlRegex.FindAllString(msgs, -1)
	if len(p.result) < 1 {
		return stacktrace.Propagate(MsgNotFound, "Failed find url in  main")
	}
	idMatcher := idRegex.FindAllStringSubmatch(msgs, -1)
	if len(idMatcher) < 1 {
		return stacktrace.Propagate(MsgNotFound, "Failed find id in  main")
	}
	p.lastId = idMatcher[len(idMatcher)-1][1]
	return nil
}

func (p *BaseProcessor) processPages() (err error) {
	var pageUrl = p.genUrl()
	p.logf("process pages....")
	req, err := http.NewRequest("GET", pageUrl, nil)
	if err != nil {
		return stacktrace.Propagate(err, "Failed new page request")
	}
	for k, _ := range p.req.Header {
		req.Header.Set(k, p.req.Header.Get(k))
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return stacktrace.Propagate(err, "Failed get page response")
	}
	bs, _ := ioutil.ReadAll(resp.Body)
	str := replacer.Replace(string(bs))
	result := urlRegex.FindAllString(str, -1)
	if len(result) < 1 {
		return stacktrace.Propagate(err, "Failed get page url")
	}
	idMatcher := idRegex.FindAllStringSubmatch(str, -1)
	if len(idMatcher) < 1 {
		return stacktrace.Propagate(err, "Failed get page id")
	}
	p.lastId = idMatcher[len(idMatcher)-1][1]
	p.logf("Page Get => %d,lastid: %s", len(result), p.lastId)
	p.result = append(p.result, result...)
	if p.lastId != "" {
		p.Sleep()
		return p.processPages()
	}
	return nil
}

func (P *BaseProcessor) Save() {

}

func (p *BaseProcessor) genUrl() string {
	url := "http://mp.weixin.qq.com/mp/getmasssendmsg?" + p.req.URL.RawQuery
	url += "&frommsgid=" + p.lastId + "&f=json&count=100"
	return url
}

func (P *BaseProcessor) logf(format string, msg ...interface{}) {
	if Verbose {
		Logger.Printf(format, msg...)
	}
}
