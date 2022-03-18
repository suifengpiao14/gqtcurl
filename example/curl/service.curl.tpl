{{define "_bodyGetOrderByOrderNumber"}}
{{- $serviceId:="110001"}}
{
    "_head":{
        "_interface":"getOrderInfo",
        "_msgType":"request",
        "_remark":"",
        "_version":"0.01",
        "_timestamps":"{{timestampSecond}}",
        "_invokeId":"dispatch_order_{{xid}}",
        "_callerServiceId":"{{$serviceId}}",
        "_groupNo":"1"
    },
    "_param":{
        "orderNum":"{{.OrderNumber}}",
        "containInfo":[
            "basic",
            "good",
            "logistics"
        ]
    }
}
{{end}}

{{define "GetOrderByOrderNumber"}}
{{- $serviceId:="110001"}}
{{- $secretKey :="wwqCxg4e3OUzILDzdD957zuVH5iHRt4W"}}
{{- $body:=jsonCompact (getBody .)}}
POST http://ordserver.huishoubao.com/order_center/getOrderInfo HTTP/1.1
Content-Type: application/json
HSB-OPENAPI-CALLERSERVICEID: {{$serviceId}}
HSB-OPENAPI-SIGNATURE: {{getMD5LOWER  $body "_" $secretKey}}

{{$body}}
{{end}}