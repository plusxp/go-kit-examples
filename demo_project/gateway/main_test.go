package main

import (
	"bytes"
	"encoding/json"
	"go-util/_util"
	"gopkg.in/jose.v1/crypto"
	"gopkg.in/jose.v1/jws"
	"hello/pb/gen-go/pb"
	"hello/pb/gen-go/pbcommon"
	"net/http"
	"testing"
	"time"
)

func TestMyGateWay_MakeADate(t *testing.T) {
	test := []struct {
		name      string
		urlPath   string
		wantReply string
	}{
		{
			name:      "[wrong url params test]",
			urlPath:   "/hello/make_a_date/2020-09-09/xxx",
			wantReply: "Sorry, I am too busy~",
		},
		{
			name:      "[right url params test]",
			urlPath:   "/hello/make_a_date/2020-10-01/xxx",
			wantReply: "OK~, I will arrive on 10.1",
		},
	}

	gwAddr := "http://127.0.0.1:8000"
	client := http.Client{
		Timeout: time.Second * 2,
	}
	for _, tt := range test {
		rsp, err := client.Get(gwAddr + tt.urlPath)
		if err != nil {
			t.Fatalf("name:%s got err:%v", tt.name, err)
		} else if rsp.StatusCode != 200 {
			t.Errorf("name:%s got err code:%v", tt.name, rsp.StatusCode)
			continue
		}
		//buf := bytes.NewBuffer(nil)
		//_, err = buf.ReadFrom(rsp.Body)
		//log.Print(111, buf.String())
		rs := new(pb.MakeADateResponse)
		err = json.NewDecoder(rsp.Body).Decode(rs)
		_util.PanicIfErr(err, nil)
		if tt.wantReply != rs.Reply {
			t.Errorf("name:%s wantReply:%s != gotReply:%s", tt.name, tt.wantReply, rs.Reply)
		}
	}
}

func genJWToken() string {
	// a simplest token
	cls := jws.Claims{
		"aud": "example_aud",
		"sub": "example_sub",
		"iss": "example_iss",
	}
	token := jws.NewJWT(cls, crypto.SigningMethodHS512)
	b, _ := token.Serialize([]byte("your_secret"))
	return string(b)
}

func TestMyGateWay_UpdateUserInfo(t *testing.T) {
	test := []struct {
		name      string
		urlPath   string
		token     string
		emptyBody bool
		want      int
	}{
		{
			name:    "[empty token test]",
			urlPath: "/hello/update_user_info",
			token:   "",
			want:    401,
		},
		{
			name:    "[right token test]",
			urlPath: "/hello/update_user_info",
			token:   genJWToken(),
			want:    200,
		},
		{
			name:      "[empty body test]",
			urlPath:   "/hello/update_user_info",
			token:     genJWToken(),
			emptyBody: true,
			want:      400,
		},
	}
	gwAddr := "http://127.0.0.1:8000"
	client := http.Client{
		Timeout: time.Second * 2,
	}

	buf := bytes.NewBuffer(nil)
	_ = json.NewEncoder(buf).Encode(&pb.UpdateUserInfoRequest{
		BaseReq: &pbcommon.BaseReq{},
		UserId:  1,
		NewName: "Nick",
	})

	for _, tt := range test {
		s := buf.String()
		if tt.emptyBody {
			s = ""
		}
		req, _ := http.NewRequest("POST", gwAddr+tt.urlPath, bytes.NewBufferString(s))
		req.Header.Set("Authorization", "BEARER "+tt.token)
		rsp, err := client.Do(req)
		if err != nil {
			t.Fatalf("name:%s got err:%v", tt.name, err)
		} else if rsp.StatusCode != tt.want {
			t.Errorf("name:%s got code:%d want code:%d", tt.name, rsp.StatusCode, tt.want)
		}
	}
}
