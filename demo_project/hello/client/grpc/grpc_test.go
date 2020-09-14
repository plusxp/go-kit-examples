package grpc

import (
	"context"
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	c := NewClient()
	defer c.Stop()
	reply, err := c.SayHi(context.Background(), "Hanmeimei")
	if err != nil {
		t.Error(err)
	}
	log.Print("rsp:", reply)
}
