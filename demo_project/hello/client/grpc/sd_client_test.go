package grpc

import (
	"context"
	"gokit_foundation"
	"hello/pb/gen-go/pb"
	"hello/pb/gen-go/pbcommon"
	"hello/pb/pbutil"
	"testing"
)

func TestSDNew(t *testing.T) {
	lgr := gokit_foundation.NewLogger(nil)
	svc := MustNewClientWithSd(lgr)
	//return
	// 0x01 SayHi
	reply, errcode := svc.SayHi(context.Background(), "Jack Ma")
	if errcode != pbcommon.R_OK {
		t.Error("SayHi err", errcode)
	}
	lgr.Log("SayHi-rsp", reply)

	// 0x02 MakeADate
	rsp2, err := svc.MakeADate(context.Background(), &pb.MakeADateRequest{
		BaseReq: pbutil.DefBaseReq(),
		DateStr: "2020-10-01",
		WantSay: "Would you like to date me?",
	})
	if err != nil {
		t.Errorf("MakeADate err %s", err)
	}
	lgr.Log("MakeADate-rsp", rsp2)

	// 0x03 UpdateUserInfo
	rsp3, err := svc.UpdateUserInfo(context.Background(), &pb.UpdateUserInfoRequest{
		BaseReq: pbutil.DefBaseReq(),
		UserId:  1,
		NewName: "Jet Li",
	})
	if err != nil {
		t.Errorf("UpdateUserInfo err %s", err)
	}
	lgr.Log("UpdateUserInfo-rsp", rsp3)
}
