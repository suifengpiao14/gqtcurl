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
	gqttpl.DataVolumeMap
}

func TestGetCURLRow(t *testing.T) {

	repo := NewRepositoryCRUL()
	err := repo.AddByDir("example", TemplatefuncMap)
	if err != nil {
		panic(err)
	}
	tplName := "curl.service.curl.GetOrderByOrderNumber"
	// data := map[string]interface{}{
	// 	"OrderNumber": "1234354",
	// }
	data := &GetOrderByOrderNumberEntity{
		OrderNumber: "1234354",
	}

	curlRow, err := repo.GetCURL(tplName, data)
	if err != nil {
		panic(err)
	}
	cmd := CURLCMD(curlRow)
	fmt.Printf("%#v", cmd)

}
