package endpoint

import (
	"context"
)

// endpoint层的实现不需要用pointer，是func类型
func (e AddSvcEndpoints) Sum(ctx context.Context, a, b int) (int, error) {
	// 注意，这里虽然实现了service，但service返回的err已经映射到response.RetCode
	// 调用时，这里的err 若!=nil，则按顺序应该是grpc.conn错误，断路器等中间件返回的err
	// 此时api网关应该返回 类似503的server内部错误，而不是再读取response.RetCode，因为读取到的不是被调用方返回的，而是默认的
	resp, err := e.SumEndpoint(ctx, &SumRequest{A: a, B: b})
	// 如果这里resp是nil，要么是代码bug，要么是网络导致
	// 我们必须在这里处理nil的情况，因为client调用时会直接调用endpoint层，这里的步骤是 client---grpc-->endpoint-->service
	if resp == nil {
		return 0, err
	}
	response := resp.(*SumResponse)
	return response.V, err
}

func (e AddSvcEndpoints) Concat(ctx context.Context, a, b string) (string, error) {
	resp, err := e.ConcatEndpoint(ctx, &ConcatRequest{A: a, B: b})
	if resp == nil {
		return "", err
	}
	response := resp.(*ConcatResponse)
	return response.V, err
}
