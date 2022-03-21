package gqtcurl

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"text/template"

	"github.com/suifengpiao14/gqt/v2/gqttpl"

	"github.com/pkg/errors"
)

const (
	EOF = "\n"
)

var BodyTemplateNamePrefix = "_body"
var DataVolumeMapBodyKey = "_body"

// RepositoryCURL stores CURL templates.
type RepositoryCURL struct {
	templates map[string]*template.Template // namespace: template
}

// NewRepositoryCRUL create a new Repository.
func NewRepositoryCRUL() *RepositoryCURL {
	return &RepositoryCURL{
		templates: make(map[string]*template.Template),
	}
}

type RequestData struct {
	URL     string         `json:"url"`
	Method  string         `json:"method"`
	Header  http.Header    `json:"header"`
	Cookies []*http.Cookie `json:"cookies"`
	Body    string         `json:"body"`
}

var LeftDelim = "{{"
var RightDelim = "}}"

type CURLRow struct {
	Name        string
	Namespace   string
	Arguments   interface{}
	RequestData *RequestData
	Response    interface{}
}

type TplEntity interface {
	TplName() string
}

func (r *RepositoryCURL) AddByDir(root string, funcMap template.FuncMap) (err error) {
	r.templates, err = gqttpl.AddTemplateByDir(root, gqttpl.CURLNamespaceSuffix, funcMap, LeftDelim, RightDelim)
	if err != nil {
		return
	}
	return
}

func (r *RepositoryCURL) AddByFS(fsys fs.FS, root string, funcMap template.FuncMap) (err error) {
	r.templates, err = gqttpl.AddTemplateByFS(fsys, root, gqttpl.CURLNamespaceSuffix, funcMap, LeftDelim, RightDelim)
	if err != nil {
		return
	}
	return
}

func (r *RepositoryCURL) AddByNamespace(namespace string, content string, funcMap template.FuncMap) (err error) {
	t, err := gqttpl.AddTemplateByStr(namespace, content, funcMap, LeftDelim, RightDelim)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	r.templates[namespace] = t
	return
}

// 将模板名称，模板中的变量，封装到结构体中，使用结构体访问，避免拼写错误以及分散的硬编码，可以配合 gqttool 自动生成响应的结构体
func (r *RepositoryCURL) GetCURLRowByTplEntity(t TplEntity) (curlRow *CURLRow, err error) {
	return r.GetCURL(t.TplName(), t)
}

// GetCURLRowByTplEntityRef 支持只返回error 函数签名
func (r *RepositoryCURL) GetCURLRowByTplEntityRef(t TplEntity, curlRow *CURLRow) (err error) {
	curlRow1, err := r.GetCURLRowByTplEntity(t)
	if err != nil {
		return err
	}
	*curlRow = *curlRow1
	return
}

type DataInterface interface{}

type DataWrap struct {
	DataInterface
	Body string
	gqttpl.DataVolumeMap
}

func (r *RepositoryCURL) GetCURLByTplEntity(tplEntity gqttpl.TplEntityInterface) (curlRow *CURLRow, err error) {
	return r.GetCURL(tplEntity.TplName(), tplEntity)
}

func (r *RepositoryCURL) GetCURL(fullname string, data interface{}) (curlRow *CURLRow, err error) {

	if dataVolume, ok := gqttpl.Interface2DataVolume(&data); ok {
		namespace, name := gqttpl.SplitFullname(fullname)
		bodyName := fmt.Sprintf("%s%s", BodyTemplateNamePrefix, name)
		bodyFullname := gqttpl.SpellFullname(namespace, bodyName)
		bodyTplDefine, err := gqttpl.ExecuteTemplateTry(r.templates, bodyFullname, &data)
		if err != nil {
			return nil, err
		}
		if bodyTplDefine != nil {
			dataVolume.SetValue(DataVolumeMapBodyKey, bodyTplDefine.Output)
		}
	}

	var tplDefine *gqttpl.TPLDefine
	tplDefine, err = gqttpl.ExecuteTemplate(r.templates, fullname, data)
	if err != nil {
		return nil, err
	}
	curlRow, err = r.TplDefine2CURLRow(tplDefine)
	return
}

func (r *RepositoryCURL) TplDefine2CURLRow(tplDefine *gqttpl.TPLDefine) (curlRow *CURLRow, err error) {
	curlRow = &CURLRow{
		Name:      tplDefine.Name,
		Arguments: tplDefine.Input,
		Namespace: tplDefine.Namespace,
	}
	req, err := r.ReadRequest(tplDefine.Output)
	if err != nil {
		return
	}
	requestData, err := Request2RequestData(req)
	if err != nil {
		return
	}
	curlRow.RequestData = requestData
	return
}

func (r *RepositoryCURL) ReadRequest(httpRaw string) (req *http.Request, err error) {
	httpRaw = strings.Trim(httpRaw, "\r\n\t\v\f ") // （删除前后空格，对于没有body 内容的请求，后面再加上必要的换行）
	if httpRaw == "" {
		err = errors.Errorf("http raw not allow empty")
		return nil, err
	}
	httpRaw = strings.ReplaceAll(httpRaw, "\r\n", "\n") // 统一换行符
	divider := EOF + EOF                                // 连续2个换行符表示body开始
	// 插入body长度头部信息
	bodyIndex := strings.Index(httpRaw, divider)
	formatHttpRaw := httpRaw
	if bodyIndex > 0 {
		headerRaw := strings.TrimSpace(httpRaw[:bodyIndex])
		bodyRaw := httpRaw[bodyIndex+len(divider):]
		bodyLen := len(bodyRaw)
		formatHttpRaw = fmt.Sprintf("%s%sContent-Length: %d%s%s", headerRaw, EOF, bodyLen, divider, bodyRaw)
	} else {
		// 如果没有请求体，则原始字符后面必须保留一个换行符
		formatHttpRaw = fmt.Sprintf("%s%s", formatHttpRaw, divider)
	}

	buf := bufio.NewReader(strings.NewReader(formatHttpRaw))
	req, err = http.ReadRequest(buf)
	if err != nil {
		return
	}
	if req.URL.Scheme == "" {
		queryPre := ""
		if req.URL.RawQuery != "" {
			queryPre = "?"
		}
		req.RequestURI = fmt.Sprintf("http://%s%s%s%s", req.Host, req.URL.Path, queryPre, req.URL.RawQuery)
	}

	return
}

func Request2RequestData(req *http.Request) (requestData *RequestData, err error) {
	requestData = &RequestData{}
	bodyByte, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}
	req.Header.Del("Content-Length")
	requestData = &RequestData{
		URL:     req.URL.String(),
		Method:  req.Method,
		Header:  req.Header,
		Cookies: req.Cookies(),
		Body:    string(bodyByte),
	}

	return
}
