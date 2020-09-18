package grpc

import (
	"context"
	"hello/pb/gen-go/pbcommon"
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	c := New()
	defer c.Stop()
	reply, err := c.SayHi(context.Background(), "Jack Ma")
	if err != pbcommon.R_OK {
		t.Error(err)
	}
	log.Print("rsp:", reply)
}
