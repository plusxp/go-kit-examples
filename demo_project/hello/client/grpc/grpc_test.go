package grpc

import (
	"context"
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	c := New()
	defer c.Stop()
	reply, err := c.SayHi(context.Background(), "Jack Ma")
	if err != nil {
		t.Error(err)
	}
	log.Print("rsp:", reply)
}
