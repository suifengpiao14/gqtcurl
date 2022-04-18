package gqtcurl

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/suifengpiao14/gqt/v2/gqttpl"
)

var curlClient *http.Client
var curlClientOnce sync.Once

func GetCurl() *http.Client {
	if curlClient == nil {
		curlClient = InitHTTPClient()
	}
	return curlClient
}

func InitHTTPClient() *http.Client {
	// 使用单例创建client
	var client *http.Client
	curlClientOnce.Do(func() {
		client = &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   10 * time.Second, // 连接超时时间
					KeepAlive: 30 * time.Second, // 连接保持超时时间
				}).DialContext,
				MaxIdleConns:        2000,             // 最大连接数,默认0无穷大
				MaxIdleConnsPerHost: 2000,             // 对每个host的最大连接数量(MaxIdleConnsPerHost<=MaxIdleConns)
				IdleConnTimeout:     90 * time.Second, // 多长时间未使用自动关闭连接
			},
		}
	})
	return client
}

//CURLRawExecAsync
func CURLRawExecAsync(curlRepository func() *RepositoryCURL, httpClient func() *http.Client, entity gqttpl.TplEntityInterface) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("CURLRawExecAsync panic err:%#v", err)
		}
	}()
	curlRow := &CURLRow{}
	err := curlRepository().GetCURLRowByTplEntityRef(entity, curlRow)
	if err != nil {
		logger.Errorf("curlRepository().GetCURLRowByTplEntityRef Err: %s", err.Error())
		return
	}
	requestData := curlRow.RequestData

	// 3 分钟超时
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Minute)
	req, err := http.NewRequestWithContext(ctx, requestData.Method, requestData.URL, bytes.NewReader([]byte(requestData.Body)))
	if err != nil {
		logger.Errorf("http.NewRequestWithContext() Err: %s", err.Error())
		return
	}

	for k, vArr := range requestData.Header {
		for _, v := range vArr {
			req.Header.Add(k, v)
		}
	}

	rsp, err := httpClient().Do(req)
	if err != nil {
		logger.Errorf("client.Do() Err: %s", err.Error())
		return
	}
	defer rsp.Body.Close()
	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		logger.Errorf("response  Err: %s", err.Error())
	}
	if err != nil {
		response := string(b)
		logger.Infof("response : %s", response)
	}
	return
}
