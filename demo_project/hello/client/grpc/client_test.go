package grpc

import (
	"context"
	"gokit_foundation"
	"hello/pb/gen-go/pbcommon"
	"testing"
)

func TestNew(t *testing.T) {
	lgr := gokit_foundation.NewLogger(nil)
	c := MustNew(lgr)
	// 在日常开发中如果是频繁使用同一个RPC client，则不需要使用完立即close，而是完全不需要时再close，频繁创建连接也是开销
	// 也可以使用sync.Pool封装client，调用者不再关心close问题
	defer c.Close()
	reply, err := c.SayHi(context.Background(), "Jack Ma")
	if err != pbcommon.R_OK {
		t.Error(err)
	}
	lgr.Log("rsp", reply)
}
