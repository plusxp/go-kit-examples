package client

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"go-util/_util"
	"testing"
)

/*
模拟client对server进行调用测试，对于client来说，操作的是endpoint，比以往的req&rsp模式更为便捷
*/

func TestSum(t *testing.T) {
	consulAddr := "192.168.1.168:8500"

	svc, err := New(consulAddr, log.NewNopLogger())
	_util.PanicIfErr(err, nil)

	r, err := svc.Sum(context.Background(), 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Sum rsp:%d\n", r)
}
