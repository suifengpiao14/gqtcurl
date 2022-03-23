// Copyright 2016 Davide Muzzarelli. All right reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gqtcurl

import (
	"fmt"
	"testing"

	"github.com/suifengpiao14/gqt/v2/gqttpl"
)

type GetOrderByOrderNumberEntity struct {
	OrderNumber string
	ServiceId   string
	SecretKey   string
	Ma          *map[string]interface{}
	*gqttpl.DataVolumeMap
}

func TestGetCURLRow(t *testing.T) {

	repo := NewRepositoryCURL()
	err := repo.AddByDir("example", TemplatefuncMap)
	if err != nil {
		panic(err)
	}
	tplName := "curl.service.curl.GetOrderByOrderNumber"
	// data := map[string]interface{}{
	// 	"OrderNumber": "1234354",
	// 	"ServiceId":   "110001",
	// 	"SecretKey":   "wwqCxg4e3OUzILDzdD957zuVH5iHRt4W",
	// }
	data := gqttpl.DataVolumeMap{
		"OrderNumber": "1234354",
		"ServiceId":   "110001",
		"SecretKey":   "wwqCxg4e3OUzILDzdD957zuVH5iHRt4W",
	}
	// data := GetOrderByOrderNumberEntity{
	// 	OrderNumber: "1234354",
	// 	ServiceId:   "110001",
	// 	SecretKey:   "wwqCxg4e3OUzILDzdD957zuVH5iHRt4W",
	// }

	curlRow, err := repo.GetCURL(tplName, &data)
	if err != nil {
		panic(err)
	}
	args := curlRow.Arguments
	dataVolumeMap, _ := args.(*gqttpl.DataVolumeMap)
	body, _ := dataVolumeMap.GetValue(BodyTemplateNamePrefix)
	bodyStr, _ := body.(string)
	bodyStr1, err := JsonCompact(bodyStr)
	if err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%s_%s", bodyStr1, data["SecretKey"])
	sc := GetMD5LOWER(s)

	fmt.Println(sc)
	str := fmt.Sprintf("%s_%s", curlRow.RequestData.Body, (data)["SecretKey"])
	fmt.Println(str)
	cmd := CURLCMD(curlRow)
	fmt.Printf("%#v", cmd)

}
