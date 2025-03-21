package proxy

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Proxy struct {
	Id         string             `db:"id" json:"id"`
	Name       string             `db:"name" json:"name"`
	Listen     string             `db:"listen" json:"listen"`
	Target     string             `db:"target" json:"target"`
	ReqHeaders types.JSONMap[any] `db:"req_headers" json:"req_headers"`
	ResHeaders types.JSONMap[any] `db:"res_headers" json:"res_headers"`
	Https      bool               `db:"https" json:"https"`
	Rebase     bool               `db:"rebase" json:"rebase"`
	Created    types.DateTime     `db:"created" json:"created"`
	Updated    types.DateTime     `db:"updated" json:"updated"`
}

type ProxyLog struct {
	Id         string         `db:"id" json:"id"`
	Proxy      string         `db:"proxy" json:"proxy"`
	ClientIP   string         `db:"client_ip" json:"client_ip"`
	Method     string         `db:"method" json:"method"`
	Target     string         `db:"target" json:"target"`
	StatusCode int            `db:"status_code" json:"status_code"`
	Error      string         `db:"error" json:"error"`
	Duration   int            `db:"duration" json:"duration"`
	ReqHeaders string         `db:"req_headers" json:"req_headers"`
	ResHeaders string         `db:"res_headers" json:"res_headers"`
	Created    types.DateTime `db:"created" json:"created"`
	Updated    types.DateTime `db:"updated" json:"updated"`
}

func GetAllProxy(app *pocketbase.PocketBase) (records []*core.Record, err error) {
	records, err = app.FindAllRecords("proxy")

	if err != nil {
		return nil, err
	}

	return records, nil
}

func FindProxyByPath(app *pocketbase.PocketBase, path string) (record *core.Record) {
	record, err := app.FindFirstRecordByFilter("proxy", "listen={:listen}", dbx.Params{"listen": path})
	if err != nil {
		return nil
	}
	return record
}

// 获取请求参数 并合并到目标url
func GetTargetMergeQuery(target string, query url.Values) string {
	targetUrl, err := url.Parse(target)
	if err != nil {
		return target
	}
	targetQuery := targetUrl.Query()
	for k, v := range query {
		targetQuery[k] = v
	}
	targetUrl.RawQuery = targetQuery.Encode()
	target = targetUrl.String()

	// 去掉 http:// 或 https:// 并返回
	scheme := "http://"
	if strings.HasPrefix(target, "https://") {
		scheme = "https://"
	}

	target = strings.TrimPrefix(target, scheme)

	return target
}

func JsonMapToString(jsonMap types.JSONMap[any]) string {
	jsonString, err := json.Marshal(jsonMap)
	if err != nil {
		return ""
	}
	return string(jsonString)
}

func RunProxy(app *pocketbase.PocketBase, c *core.RequestEvent, record *core.Record) (err error) {
	name := record.GetString("name")
	log.Println("RunProxy", name, "开始")
	start := time.Now()

	query := c.Request.URL.Query()
	target := record.GetString("target")

	rebase := record.GetBool("rebase")
	https := record.GetBool("https")
	scheme := "http"
	if https {
		scheme = "https"
	}

	targetUrl := GetTargetMergeQuery(target, query)

	reqHeaders := types.JSONMap[any]{}
	resHeaders := types.JSONMap[any]{}

	record.UnmarshalJSONField("req_headers", &reqHeaders)
	record.UnmarshalJSONField("res_headers", &resHeaders)

	// 设置请求头
	headers := make(http.Header)
	for k, v := range reqHeaders {
		if str, ok := v.(string); ok {
			headers.Set(k, str)
		}
	}

	if rebase {
		parsedUrl, err := url.Parse(scheme + "://" + target)
		if err == nil {
			headers.Set("Host", parsedUrl.Host)
			headers.Set("Referer", parsedUrl.String())
			headers.Set("Origin", fmt.Sprintf("%s://%s", scheme, parsedUrl.Host))
		}
	}

	// 将headers重写回reqHeaders
	for k, v := range headers {
		reqHeaders[k] = v
	}

	schemeTarget := scheme + "://" + targetUrl

	log.Println("RunProxy", name, "创建请求", schemeTarget, c.Request.Method)
	// 创建请求
	req, err := http.NewRequest(c.Request.Method, schemeTarget, c.Request.Body)
	if err != nil {
		log.Println("RunProxy", name, "创建请求失败", err)
		return err
	}
	req.Header = headers

	// 发送请求
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse("http://127.0.0.1:7890")
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("RunProxy", name, "请求失败", err)
		return err
	}
	defer resp.Body.Close()

	// 设置响应头
	for k, v := range resHeaders {
		if str, ok := v.(string); ok {
			c.Response.Header().Set(k, str)
		}
	}

	// 将headers重写回resHeaders
	for k, v := range c.Response.Header() {
		resHeaders[k] = v
	}

	c.Response.WriteHeader(resp.StatusCode)

	_, err = io.Copy(c.Response, resp.Body)

	log.Println("RunProxy", name, "结束")

	proxyLog := &ProxyLog{
		Proxy:      record.GetString("id"),
		ClientIP:   c.Request.RemoteAddr,
		Method:     c.Request.Method,
		Target:     targetUrl,
		StatusCode: resp.StatusCode,
		Error:      "",
		Duration:   int(time.Since(start).Milliseconds()),
		ReqHeaders: JsonMapToString(reqHeaders),
		ResHeaders: JsonMapToString(resHeaders),
	}

	if err != nil {
		proxyLog.Error = err.Error()
	}

	RecordLog(app, proxyLog)

	return err
}

func RecordLog(app *pocketbase.PocketBase, proxyLog *ProxyLog) {
	log.Println("RecordLog", "开始记录日志", proxyLog.Proxy)

	collection, err := app.FindCollectionByNameOrId("proxy_log")
	if err != nil {
		return
	}
	record := core.NewRecord(collection)

	record.Set("proxy", proxyLog.Proxy)
	record.Set("client_ip", proxyLog.ClientIP)
	record.Set("method", proxyLog.Method)
	record.Set("target", proxyLog.Target)
	record.Set("status_code", proxyLog.StatusCode)
	record.Set("error", proxyLog.Error)
	record.Set("duration", proxyLog.Duration)
	record.Set("req_headers", proxyLog.ReqHeaders)
	record.Set("res_headers", proxyLog.ResHeaders)

	err = app.Save(record)
	if err != nil {
		log.Println("RecordLog", "记录日志失败", err)
		return
	}

	log.Println("RecordLog", "记录日志成功", record.Id)
}
