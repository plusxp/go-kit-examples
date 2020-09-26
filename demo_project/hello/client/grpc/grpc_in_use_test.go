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
	defer c.Close()
	reply, err := c.SayHi(context.Background(), "Jack Ma")
	if err != pbcommon.R_OK {
		t.Error(err)
	}
	lgr.Log("rsp", reply)
}
