package client

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"testing"
)

func TestHC(t *testing.T) {
	cc, err := grpc.Dial("192.168.1.10:8080", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := grpc_health_v1.NewHealthClient(cc)
	r, err := c.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		log.Fatalf("Check: %v", err)
	}
	log.Println(r)
}
