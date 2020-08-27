package endpoint

import "go-kit-examples/new_addsvc/pb/gen-go/resultcode"

/*
首先在endpoint层需要定义专门的req和rsp struct, 可称之为ep层的protocol

它们与transport层和service层的数据做相互转换，所以你在server的transport层会看到很多
	decode...request
	encode...response
说明：decode负责 rpc-request --> endpoint-request, 得到结果后，
	再encode endpoint-response --> rpc-response

-----------
在client端会看到很多的
	encode...request
	decode...response
说明：encode负责 endpoint-request --> rpc-request,得到结果后，再
	decode rpc-response --> endpoint-response

注：每个接口都需要一个decode.func和encode.func
*/

type SumRequest struct {
	A, B int
}

// SumResponse collects the response values for the Sum method.
type SumResponse struct {
	V       int                    `json:"v"`
	RetCode resultcode.RESULT_CODE `json:"ret_code"`
}

// ConcatRequest collects the request parameters for the Concat method.
type ConcatRequest struct {
	A, B string
}

// ConcatResponse collects the response values for the Concat method.
type ConcatResponse struct {
	V       string                 `json:"v"`
	RetCode resultcode.RESULT_CODE `json:"ret_code"`
}
