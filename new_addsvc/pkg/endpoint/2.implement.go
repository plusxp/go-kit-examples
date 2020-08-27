package endpoint

import (
	"context"
	"go-kit-examples/new_addsvc/pb/gen-go/resultcode"
)

// endpoint层的实现不需要用pointer，func类型
func (e Endpoints) Sum(ctx context.Context, a, b int) (int, resultcode.RESULT_CODE) {
	resp, _ := e.SumEndpoint(ctx, SumRequest{A: a, B: b})
	response := resp.(SumResponse)
	return response.V, response.RetCode
}
