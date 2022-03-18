package gqtcurl

import (
	"fmt"
	"strings"
)

func CURLCMD(curlRow *CURLRow) (cmd string) {
	r := curlRow.RequestData
	hArr := make([]string, 0)
	for k, v := range r.Headers {
		if strings.ToLower(k) == "Content-Length" {
			continue
		}
		head := fmt.Sprintf("-H '%s: %v'", k, v)
		hArr = append(hArr, head)
	}
	headers := strings.Join(hArr, " ")
	method := strings.ToUpper(r.Method)
	body := ""
	if r.Body != "" {
		body = fmt.Sprintf("'%s'", r.Body)
	}
	cmd = fmt.Sprintf("curl -X%s %s %s '%s'", method, headers, body, r.URL)

	return
}
