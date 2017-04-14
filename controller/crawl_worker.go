package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
	"compress/gzip"
	"io"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/reechou/holmes"
	"github.com/robertkrimen/otto"
)

type VideoInfo struct {
	Id  string
	Md5 string
}

var (
	RMDMY_VIDEO = map[int]string{
		1: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3FplJRmaQO0O0OO0O0O&token=%s",
		2: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3FplZVpbAO0O0OO0O0O&token=%s",
		3: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3FplplsaQO0O0OO0O0O&token=%s",
		4: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3FplplsagO0O0OO0O0O&token=%s",
		5: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3Fpl5dtaAO0O0OO0O0O&token=%s",
		6: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3FpmJNlbAO0O0OO0O0O&token=%s",
		7: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3FpmZVtagO0O0OO0O0O&token=%s",
		8: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3FpmZxlZgO0O0OO0O0O&token=%s",
		9: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3FpmZxlaAO0O0OO0O0O&token=%s",
		10: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3FpmZxlaQO0O0OO0O0O&token=%s",
		11: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3FpmphkagO0O0OO0O0O&token=%s",
		12: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3FpmphkawO0O0OO0O0O&token=%s",
		13: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3Fpm5NpYwO0O0OO0O0O&token=%s",
		14: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3Fpm5NpZAO0O0OO0O0O&token=%s",
		15: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3Fpm5xraQO0O0OO0O0O&token=%s",
		16: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3Fpm5xsZAO0O0OO0O0O&token=%s",
		17: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3Fqk5VqagO0O0OO0O0O&token=%s",
		18: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3Fqk5VqawO0O0OO0O0O&token=%s",
		19: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3FqlJZlawO0O0OO0O0O&token=%s",
		20: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3FqlZVtagO0O0OO0O0O&token=%s",
		21: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3FqlZprYwO0O0OO0O0O&token=%s",
		22: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3FqlZprZAO0O0OO0O0O&token=%s",
		23: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3Fql5VpZAO0O0OO0O0O&token=%s",
		24: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3Fql5VpZQO0O0OO0O0O&token=%s",
		25: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3FqmJhlZwO0O0OO0O0O&token=%s",
		26: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3FqmJlmZQO0O0OO0O0O&token=%s",
		27: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3FqmZxlawO0O0OO0O0O&token=%s",
		28: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3FqmZxlbAO0O0OO0O0O&token=%s",
		29: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3Fqm5NmZAO0O0OO0O0O&token=%s",
		30: "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3Fqm5NmZgO0O0OO0O0O&token=%s",
		
		31: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSimnGSyJSVbJo000oTk2qWmJdpa3TIxpWXmcWYboo000oXY8NxYpNpnmiPxWiVcFo000oZlmxnk55nlpaYY2YO0O0O&token=%s",
		32: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSibHBolZRllm5ixWpllmFxbnTHnGaZmsSUmYo000oXY8NxYpNpnmiPxWWWcFo000oZlmxnk55nlpaYY2YO0O0O&token=%s",
		33: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSia2yTl8lqa5timmlnlmZpmXTGl5hsa8SUbIo000oXY8NxYpNpnmiPm5lpaVo000oZlmxnk55nlpaYYpkO0O0O&token=%s",
		34: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSia22Sl5RolZ1klmtjZ5KcmHSVyZVra5qVnYo000oXY8NxYpNpnmiPnGaWalo000oZlmxnk55nlpaYYpkO0O0O&token=%s",
		35: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSiaHJlm5aXaJtpx2dlbGJscHSUlpNpacSXcIo000oXY8NxYpNpnmiPm2trn1o000oZlmxnk55nlpaYYpkO0O0O&token=%s",
		36: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSibnFmlcWYl3CUmGtqZpKZbHTIx2lonMVjmYo000oXY8NwYpNpnmiPm2SWmloo00oGlmhok22dlcWbaGsO0O0O&token=%s",
		37: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSianFjm5tlaWxrmmOYamFvnXSdm5ZtbMdmaoo000oXY8NwYpNpnmiPxJlpcFo000oZlmxnk55nlpaYY2YO0O0O&token=%s",
		38: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSim3Jlw5Nram2VmpiUZWdrZ3TJlJWZncdjbIo000oXY8NvYpNpnmiPm5dqbFo000oZlmxnk55nlpaYYpkO0O0O&token=%s",
		39: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSimXFpmJRtZWtjmmSYl2FxcHTIx2tom5SVbIo000oXY8NvYpNpnmiPxZqVbFo000oZlmxnk55nlpaYYpkO0O0O&token=%s",
		40: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSibWpjm8aalpqVm5WUaWKZbXTInJZnacVmbIo000oXY8NvYpNpnmiPnGpnbFo000oZlmxnk55nlpaYY2YO0O0O&token=%s",
		41: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSibJtjlMWVZ3Bol5iWZZVxmHTImGhpbJFoboo000oXY8NvYpNpnmiPm2tsnFo000oZlmxnk55nlpaYYpkO0O0O&token=%s",
		42: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSib2xlk5WYlWyXxGZpapaamnTIk5SZaplpaIo000oXY8NvYpNpnmiPnJdqn1o000oZlmxnk55nlpaYY2YO0O0O&token=%s",
		43: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSicXGSx8aVZm9immeUZmlubHSblmlpmJOTm4o000oXY8NvYpNpnmiPxJWWb1o000oZlmxnk55nlpaYY2YO0O0O&token=%s",
		44: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSim55lw5yYa2uYnJmYm2FtaXSZmJdsmsOYm4o000oXY8NuYpNpnmiPxWxsalo000oZlmxnk55nlpaYY2YO0O0O&token=%s",
		45: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSibnBjk5xnmJyXlmeVaZKdmXSZlGZtbMdicIo000oXY8NuYpNpnmiPm2WWnlo000oZlmxnk55nlpaYY2YO0O0O&token=%s",
		46: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSiaHKVlJpkaGyYxWhrbGJobnSYx2dnbZNnnYo000oXY8NuYpNpnmiPm5psnloo00oGlmhok22dlcWbZGoO0O0O&token=%s",
		47: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSibZ5kx5RkamqUyWlnZZSZcHSYnGeanJWUmYo000oXY8NuYpNpnmiPm5lknFo000oZlmxnk55nlpaYYpkO0O0O&token=%s",
		48: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSinm9kk5aaZZ1px2hjmJSZnHSYlGtomJRiboo000oXY8NuYpNpnmiPxZljnlo000oZlmxnk55nlpaYYpkO0O0O&token=%s",
		49: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSinG6XlJdmbG9nx5aWbGFqbXSXm5RtnZOXmYo000oXY8NuYpNpnmiPxGWValoo00oGlmhok22dlcWbaGsO0O0O&token=%s",
		50: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSimZtqxshmaJ9myZZqmpZub3SWlJiXcJpknYo000oXY8NsYpNpnmiPxWiXa1o000oZlmxnk55nlpaYYpkO0O0O&token=%s",
		
		51: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSicG5hlsRoa3JoxmNsa2ZubXSVyWJlaJGXm4o000oXY8NsYpNpnmiPnJVka1o000oZlmxnk55nlpaYYpkO0O0O&token=%s",
		52: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSibp1qxZlsaWpqk2WXa2qcZ3SVxmNumMOUcIo000oXY8NsYpNpnmiPxGlsnVoo00oGlmhok22dlcWbZGoO0O0O&token=%s",
		53: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSimmo000oVw8VmlWlnxZZjamNqm3SVm5aWnJeTboo000oXY8NsYpNpnmiPm2ZncVoo00oGlmhok22dlcWbZGoO0O0O&token=%s",
		54: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSicW5hw5qWmJ6TmWSYbWSdZ3SVlJabmJaUaoo000oXY8NsYpNpnmiPnJqWblo000oZlmxnk55nlpaYYpkO0O0O&token=%s",
		55: "https://qq.h6080.com/sapi.php?id=o5uiyZ3KrqSibW2WlpSVmW1ik2lpZWmbm3TKmJZmbsdrbIo000oXY8NrYpNpnmiPnGyUalo000oZlmxnk55nlpaYYpkO0O0O&token=%s",
	}
	RMDMY_VIDEO_MD5 = map[int]*VideoInfo{
		27: &VideoInfo{Id: "64145854", Md5: "5dac8630ef7b10a24f82ad64b6574a71"},
		28: &VideoInfo{Id: "64145787", Md5: "1b573add0fb41bc6a5e15d0933b961af"},
		29: &VideoInfo{Id: "64145860", Md5: "a5ba7b68e72d80a1df6ed7fd7a494789"},
		30: &VideoInfo{Id: "64145849", Md5: "1411b2b1fac35e9e48639dc30c3ba440"},
		31: &VideoInfo{Id: "64145789", Md5: "0d37d33a78c922d5fdec0937b8e13e8a"},
		32: &VideoInfo{Id: "64145864", Md5: "5ac40f591b0e076b85e510271db8801a"},
		33: &VideoInfo{Id: "64145859", Md5: "67c69e16ab6855fb7c87c73fd90eb3f2"},
		34: &VideoInfo{Id: "64145856", Md5: "769a312c4b5ee2e9c722a7ed52c61e54"},
		35: &VideoInfo{Id: "64145939", Md5: "9904537a6a1bee77f1cb32a35936fa86"},
		36: &VideoInfo{Id: "64145865", Md5: "03f360aaad0a26ec8ba77c5886dea830"},
		37: &VideoInfo{Id: "64145866", Md5: "73828fbeb1139f529c5c81247e2883ec"},
		38: &VideoInfo{Id: "64145857", Md5: "a75af0791ffe48b26acfe82a4c3db243"},
		39: &VideoInfo{Id: "64145847", Md5: "04ee9bb30189c5cd58ca510edfbe444d"},
		40: &VideoInfo{Id: "64145869", Md5: "b5ef8ab09e63a3b3648c622098e7db1f"},
		41: &VideoInfo{Id: "64145936", Md5: "b76f40bb3bee21913ac4e5ec26e6af66"},
		42: &VideoInfo{Id: "64145875", Md5: "bae796c85f5e59ef42598daf90f2c6de"},
		43: &VideoInfo{Id: "64145871", Md5: "cc33cae9824699381fe97c93aa4a90ba"},
		44: &VideoInfo{Id: "64145934", Md5: "1c37cc755ec9d1341de7b2ec8edd64d6"},
		45: &VideoInfo{Id: "64145873", Md5: "5873808cc0e504c82d0effd91b5b128e"},
		46: &VideoInfo{Id: "64145941", Md5: "626f2ecc570de0fb458bf3dec05e8bf2"},
		47: &VideoInfo{Id: "64145935", Md5: "9eee251c4b94a30fd83a66c4f84c785e"},
		48: &VideoInfo{Id: "64145944", Md5: "3683497ae268122a56a61cd5062c1685"},
		49: &VideoInfo{Id: "64145872", Md5: "22f1477fbd6eb2a195e57e65991db733"},
		50: &VideoInfo{Id: "64145940", Md5: "1549c61bfd65011368c085883bb64b93"},
		51: &VideoInfo{Id: "64145870", Md5: "c6d12b3652cce0d354e0c2d533ab6912"},
		52: &VideoInfo{Id: "64145942", Md5: "592d36b545a9c95b50ddcb6e5311b4b7"},
		53: &VideoInfo{Id: "64145938", Md5: "8f4296aa283e573267c21dd8d6393a3e"},
		54: &VideoInfo{Id: "64145858", Md5: "216977d46b074b4ebc99ccb9616a9b24"},
		55: &VideoInfo{Id: "64145861", Md5: "a57f0300f9aada8e3070d1042550616f"},
	}
)

type CrawlWorker struct {
	sync.Mutex
	videoUrlMap map[int]string

	stop chan struct{}
	done chan struct{}
}

func NewCrawlWorker() *CrawlWorker {
	cw := &CrawlWorker{
		videoUrlMap: make(map[int]string),
		stop:        make(chan struct{}),
		done:        make(chan struct{}),
	}
	go cw.run()

	return cw
}

func (self *CrawlWorker) Stop() {
	close(self.stop)
	<-self.done
}

func (self *CrawlWorker) GetVideoUrl(num int) string {
	self.Lock()
	defer self.Unlock()
	return self.videoUrlMap[num]
}

func (self *CrawlWorker) run() {
	self.crawl()
	for {
		select {
		case <-time.After(time.Minute):
			self.crawl()
		case <-self.stop:
			close(self.done)
			return
		}
	}
}

func (self *CrawlWorker) crawl() {
	var wg sync.WaitGroup
	for k, v := range RMDMY_VIDEO {
		wg.Add(1)
		go func(key int, url string) {
			defer wg.Done()
			srcUrl, err := self.getVideoInfo(url)
			if err != nil {
				holmes.Error("get video error: %v", err)
				return
			}
			self.Lock()
			self.videoUrlMap[key] = srcUrl
			self.Unlock()
		}(k, v)
	}
	wg.Wait()
	//fmt.Println(self.videoUrlMap)
	holmes.Debug("map: %v", self.videoUrlMap)

	//for k, v := range RMDMY_VIDEO {
	//	doc, err := self.getHtmlDoc(v)
	//	if err != nil {
	//		holmes.Error("[%d] [%s] new document error: %v", k, v, err)
	//		continue
	//	}
	//	fmt.Println(doc.Html())
	//	src, ifExit := doc.Find("video").Attr("src")
	//	fmt.Println(src, ifExit)
	//}
}

func (self *CrawlWorker) getHtmlDoc(url string) (*goquery.Document, error) {
	client := &http.Client{}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Referer", "http://api.svip.baiyug.cn/svip/index.php?url=64145864")
	request.Header.Add("Host", "api.svip.baiyug.cn")
	request.Header.Add("Upgrade-Insecure-Requests", "1")
	request.Header.Add("X-Requested-With", "XMLHttpRequest")
	request.Header.Add("Cookie", "Hm_lvt_dedd7873131db605c54f2ab012bb969c=1492146059,1492146080,1492146118,1492146122; Hm_lpvt_dedd7873131db605c54f2ab012bb969c=1492146144; Hm_lvt_77759e7c3da990eb50c6eda7c384e874=1492146060; Hm_lpvt_77759e7c3da990eb50c6eda7c384e874=1492146144; UM_distinctid=15b6ad7b29baa-0175e2697c0f9b-37657900-fa000-15b6ad7b29c10e; Hm_lvt_cfe875cb9404083b70c9a0e437ed57a0=1492146173; Hm_lpvt_cfe875cb9404083b70c9a0e437ed57a0=1492146187; Hm_lvt_901bde7452b7f689c3fd3d7b183cef9e=1492146173; Hm_lpvt_901bde7452b7f689c3fd3d7b183cef9e=1492146187")
	request.Header.Add("User-Agent", "Mozilla/5.0 (iPad; CPU OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return goquery.NewDocumentFromResponse(response)
}

type RspVideoInfo struct {
	Msg string `json:"msg"`
	Ext string `json:"ext"`
	Url string `json:"url"`
}

//func (self *CrawlWorker) getVideoInfo(info *VideoInfo) (string, error) {
//	resp, err := http.PostForm("http://api.svip.baiyug.cn/vip_op_yun_h6/url.php",
//		url.Values{"id": {info.Id}, "type": {"mmsid2"}, "siteuser": {"123"}, "md5": {info.Md5}, "hd": {"yh"}})
//	if err != nil {
//		holmes.Error("post form error: %v", err)
//		return "", err
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return "", err
//	}
//	var rsp RspVideoInfo
//	err = json.Unmarshal(body, &rsp)
//	if err != nil {
//		holmes.Error("json unmarshal[%s] error: %v", string(body), err)
//		return "", err
//	}
//	vUrl, err := url.QueryUnescape(rsp.Url)
//	if err != nil {
//		return "", err
//	}
//	return self.getShortUrl(vUrl)
//}

func (self *CrawlWorker) getVideoInfo(urlStr string) (string, error){
	client := &http.Client{}
	queryUrl := fmt.Sprintf(urlStr, self.getToken())
	request, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Accept-Encoding", "gzip, deflate, sdch, br")
	request.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,de;q=0.6,en;q=0.4,ko;q=0.2,pt;q=0.2,zh-TW;q=0.2")
	request.Header.Set("Upgrade-Insecure-Requests", "1")
	request.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1")
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	reader, err := gzip.NewReader(response.Body)
	if err != nil {
		holmes.Error("gzip new reader error: %v", err)
		return "", err
	}
	
	buff := make([]byte, 1024)
	for {
		n, err := reader.Read(buff)
		if err != nil && err != io.EOF {
			holmes.Error("reader read error: %v", err)
			return "", nil
		}
		
		if n == 0 {
			break
		}
	}
	
	s := fmt.Sprintf("%s", buff)
	reg := regexp.MustCompile(`src(.*?)=(.*?)\"(.*?)\"`)
	src := reg.FindString(s)
	src = strings.Replace(src, "\"", "", -1)
	src = strings.Replace(src, "src=", "", -1)
	if src == "" {
		holmes.Error("get url[%s] s[%s] error cannot found src", queryUrl, s)
		return "", fmt.Errorf("src == nil")
	}
	
	return self.getShortUrl(src)
}

func (self *CrawlWorker) getToken() string {
	vm := otto.New()
	vm.Run(`
	    function Base64() {
			_keyStr = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=";
			this.encode = function (input) {
				var output = "";
				var chr1, chr2, chr3, enc1, enc2, enc3, enc4;
				var i = 0;
				input = _utf8_encode(input);
				while (i < input.length) {
					chr1 = input.charCodeAt(i++);
					chr2 = input.charCodeAt(i++);
					chr3 = input.charCodeAt(i++);
					enc1 = chr1 >> 2;
					enc2 = ((chr1 & 3) << 4) | (chr2 >> 4);
					enc3 = ((chr2 & 15) << 2) | (chr3 >> 6);
					enc4 = chr3 & 63;
					if (isNaN(chr2)) {
						enc3 = enc4 = 64;
					} else if (isNaN(chr3)) {
						enc4 = 64;
					}
					output = output +
					_keyStr.charAt(enc1) + _keyStr.charAt(enc2) +
					_keyStr.charAt(enc3) + _keyStr.charAt(enc4);
				}
				return output;
			}
			this.decode = function (input) {
				var output = "";
				var chr1, chr2, chr3;
				var enc1, enc2, enc3, enc4;
				var i = 0;
				input = input.replace(/[^A-Za-z0-9\+\/\=]/g, "");
				while (i < input.length) {
					enc1 = _keyStr.indexOf(input.charAt(i++));
					enc2 = _keyStr.indexOf(input.charAt(i++));
					enc3 = _keyStr.indexOf(input.charAt(i++));
					enc4 = _keyStr.indexOf(input.charAt(i++));
					chr1 = (enc1 << 2) | (enc2 >> 4);
					chr2 = ((enc2 & 15) << 4) | (enc3 >> 2);
					chr3 = ((enc3 & 3) << 6) | enc4;
					output = output + String.fromCharCode(chr1);
					if (enc3 != 64) {
						output = output + String.fromCharCode(chr2);
					}
					if (enc4 != 64) {
						output = output + String.fromCharCode(chr3);
					}
				}
				output = _utf8_decode(output);
				return output;
			}
			_utf8_encode = function (string) {
				string = string.replace(/\r\n/g,"\n");
				var utftext = "";
				for (var n = 0; n < string.length; n++) {
					var c = string.charCodeAt(n);
					if (c < 128) {
						utftext += String.fromCharCode(c);
					} else if((c > 127) && (c < 2048)) {
						utftext += String.fromCharCode((c >> 6) | 192);
						utftext += String.fromCharCode((c & 63) | 128);
					} else {
						utftext += String.fromCharCode((c >> 12) | 224);
						utftext += String.fromCharCode(((c >> 6) & 63) | 128);
						utftext += String.fromCharCode((c & 63) | 128);
					}
		 
				}
				return utftext;
			}
			_utf8_decode = function (utftext) {
				var string = "";
				var i = 0;
				var c = c1 = c2 = 0;
				while ( i < utftext.length ) {
					c = utftext.charCodeAt(i);
					if (c < 128) {
						string += String.fromCharCode(c);
						i++;
					} else if((c > 191) && (c < 224)) {
						c2 = utftext.charCodeAt(i+1);
						string += String.fromCharCode(((c & 31) << 6) | (c2 & 63));
						i += 2;
					} else {
						c2 = utftext.charCodeAt(i+1);
						c3 = utftext.charCodeAt(i+2);
						string += String.fromCharCode(((c & 15) << 12) | ((c2 & 63) << 6) | (c3 & 63));
						i += 3;
					}
				}
				return string;
			}
		}
		eval(function(p,a,c,k,e,d){e=function(c){return(c<a?"":e(parseInt(c/a)))+((c=c%a)>35?String.fromCharCode(c+29):c.toString(36))};if(!''.replace(/^/,String)){while(c--)d[e(c)]=k[c]||e(c);k=[function(e){return d[e]}];e=function(){return'\\w+'};c=1;};while(c--)if(k[c])p=p.replace(new RegExp('\\b'+e(c)+'\\b','g'),k[c]);return p;}('1 5="g";1 h=(+6 7).c().d(0,8);1 4="9+"+h+"@"+5;1 2=6 e();1 a=2.3(h);1 b=2.3(4);1 f=2.3(a+b);',18,18,'|var|base86|encode|rgfgb|www|new|Date||content|||toString|slice|Base64|encrypted|caonimabi|'.split('|'),0,{}));
	`)
	var token string
	if value, err := vm.Get("encrypted"); err == nil {
		if valueStr, err := value.ToString(); err == nil {
			token = valueStr
		}
	}
	return token
}

type RspShortUrl struct {
	UrlShort string `json:"url_short"`
	UrlLong  string `json:"url_long"`
	Type     int    `json:"type"`
}

func (self *CrawlWorker) getShortUrl(urlStr string) (string, error) {
	queryUrl := fmt.Sprintf("http://api.t.sina.com.cn/short_url/shorten.json?source=209678993&url_long=%s", url.QueryEscape(urlStr))
	resp, err := http.Get(queryUrl)
	if err != nil {
		holmes.Error("get short url error: %v", err)
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		holmes.Error("io read error: %v", err)
		return "", err
	}

	if strings.Contains(string(body), "error_code") {
		holmes.Error("urlstr:%s rsp:%s error", urlStr, string(body))
		return "", fmt.Errorf("response: %s error", string(body))
	}

	var rsp []RspShortUrl
	err = json.Unmarshal(body, &rsp)
	if err != nil {
		holmes.Error("[%s] json unmarshal error: %v", string(body), err)
		return "", err
	}
	if len(rsp) == 0 {
		return "", fmt.Errorf("len(rsp) == 0")
	}
	holmes.Debug("get short url: %v", rsp[0])
	return rsp[0].UrlShort, nil
}

func (self *CrawlWorker) getTest() {
	client := &http.Client{}
	
	request, err := http.NewRequest("GET", "http://play.h6080.com/jx/sapi.php?id=n5o000oiyZ3KrqSia3FpmZxlZgO0O0OO0O0O&token=TVRRNU1qRTJPVGs9WTI5dWRHVnVkQ3N4TkRreU1UWTVPVUJqWVc5dWFXMWhZbWs9", nil)
	if err != nil {
		return
	}
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Accept-Encoding", "gzip, deflate, sdch, br")
	request.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,de;q=0.6,en;q=0.4,ko;q=0.2,pt;q=0.2,zh-TW;q=0.2")
	request.Header.Set("Upgrade-Insecure-Requests", "1")
	request.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1")
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	reader, err := gzip.NewReader(response.Body)
	if err != nil { panic(err) }
	
	buff := make([]byte, 1024)
	for {
		n, err := reader.Read(buff)
		
		if err != nil && err != io.EOF {
			panic(err)
		}
		
		if n == 0 {
			break
		}
	}
	
	s := fmt.Sprintf("%s", buff)
	fmt.Println(s)
}
