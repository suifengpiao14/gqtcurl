package gqtcurl

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/suifengpiao14/gqt/v2"
	"github.com/suifengpiao14/gqt/v2/gqttpl"
)

var TemplatefuncMap = template.FuncMap{
	"zeroTime":          gqt.ZeroTime,
	"currentTime":       gqt.CurrentTime,
	"permanentTime":     gqt.PermanentTime,
	"newPreComma":       gqt.NewPreComma,
	"toCamel":           gqttpl.ToCamel,
	"toLowerCamel":      gqttpl.ToLowerCamel,
	"snakeCase":         gqttpl.SnakeCase,
	"Contains":          strings.Contains,
	"fen2yuan":          Fen2yuan,
	"timestampSecond":   TimestampSecond,
	"xid":               Xid,
	"withDefault":       WithDefault,
	"withEmptyStr":      WithEmptyStr,
	"withZeroNumber":    WithZeroNumber,
	"getMD5LOWER":       GetMD5LOWER,
	"jsonCompact":       JsonCompact,
	"standardizeSpaces": gqttpl.StandardizeSpaces,
	"getBody":           GetBody,
}

func GetMD5LOWER(s ...string) string {
	allStr := strings.Join(s, "")
	return gqt.GetMD5LOWER(allStr)
}

func Fen2yuan(fen interface{}) string {
	var yuan float64
	intFen, ok := fen.(int)
	if ok {
		yuan = float64(intFen) / 100
		return strconv.FormatFloat(yuan, 'f', 2, 64)
	}
	strFen, ok := fen.(string)
	if ok {
		intFen, err := strconv.Atoi(strFen)
		if err == nil {
			yuan = float64(intFen) / 100
			return strconv.FormatFloat(yuan, 'f', 2, 64)
		}
	}
	return strFen
}

// 秒计数的时间戳
func TimestampSecond() int64 {
	return time.Now().Unix()
}

func Xid() string {
	guid := xid.New()
	return guid.String()
}

// 模板中预先写入的变量，在接口中可能没有传该字段，此时会出现<no value> ，需要使用默认值
func WithDefault(val interface{}, def interface{}) interface{} {
	if val == nil {
		val = def
	}
	return val
}

func WithEmptyStr(val interface{}) interface{} {
	def := ""
	return WithDefault(val, def)
}

func WithZeroNumber(val interface{}) interface{} {
	def := ""
	return WithDefault(val, def)
}

func JsonCompact(src string) (out string, err error) {
	var buff bytes.Buffer
	err = json.Compact(&buff, []byte(src))
	if err != nil {
		return
	}
	out = buff.String()
	return

}

func GetBody(data interface{}) (body string, err error) {

	body = ""
	dataVolume, ok := gqttpl.Interface2DataVolume(data)
	if !ok {
		err = errors.Errorf("func geBody only used for object who implement gqttpl.DataVolumeInterface")
		return
	}
	bodyInterface, ok := dataVolume.GetValue(DataVolumeMapBodyKey)
	if !ok {
		err = errors.Errorf("not defined body template _bodyXXX[XXX is template define name] or not use ")
		return
	}
	body, ok = bodyInterface.(string)
	if !ok {
		err = errors.Errorf(" body %#v can not convert to string", bodyInterface)
	}

	return
}
