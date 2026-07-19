package message

import (
	"encoding/base64"
	"testing"

	"github.com/xifan2333/blivedm-go/pb"
	"google.golang.org/protobuf/proto"
)

func TestInteractWordV2ParseProtobuf(t *testing.T) {
	src := &pb.InteractWordV2{
		Uid:       123456,
		Uname:     "测试用户",
		MsgType:   InteractMsgTypeEnter,
		Timestamp: 1710000000,
		Uinfo: &pb.InteractWordV2UserInfo{
			Base: &pb.InteractWordV2UserBaseInfo{
				Face: "https://i0.hdslb.com/bfs/face/test.jpg",
			},
		},
	}
	raw, err := proto.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}
	b64 := base64.StdEncoding.EncodeToString(raw)
	body := []byte(`{"cmd":"INTERACT_WORD_V2","data":{"pb":"` + b64 + `"}}`)

	got := new(InteractWordV2)
	got.Parse(body)

	if got.Uid != 123456 {
		t.Fatalf("uid=%d", got.Uid)
	}
	if got.Uname != "测试用户" {
		t.Fatalf("uname=%q", got.Uname)
	}
	if got.MsgType != InteractMsgTypeEnter {
		t.Fatalf("msg_type=%d", got.MsgType)
	}
	if got.Timestamp != 1710000000 {
		t.Fatalf("timestamp=%d", got.Timestamp)
	}
	if got.Face != "https://i0.hdslb.com/bfs/face/test.jpg" {
		t.Fatalf("face=%q", got.Face)
	}
	if !got.IsEnter() {
		t.Fatal("expected IsEnter")
	}
	if got.IsLike() {
		t.Fatal("unexpected IsLike")
	}
}

func TestInteractWordV2ParseLike(t *testing.T) {
	src := &pb.InteractWordV2{
		Uid:     1,
		Uname:   "liker",
		MsgType: InteractMsgTypeLike,
	}
	raw, _ := proto.Marshal(src)
	body := []byte(`{"cmd":"INTERACT_WORD_V2","data":{"pb":"` + base64.StdEncoding.EncodeToString(raw) + `"}}`)
	got := new(InteractWordV2)
	got.Parse(body)
	if !got.IsLike() || got.IsEnter() {
		t.Fatalf("msg_type=%d", got.MsgType)
	}
}

func TestParseInteractWordJSON(t *testing.T) {
	data := `{"cmd":"INTERACT_WORD","data":{"uid":99,"uname":"旧协议","msg_type":1,"timestamp":100,"roomid":1}}`
	w := ParseInteractWordJSON(data)
	if w.Uid != 99 || w.Uname != "旧协议" || w.MsgType != 1 {
		t.Fatalf("%+v", w)
	}
	// data-only
	w2 := ParseInteractWordJSON(`{"uid":7,"uname":"x","msg_type":2}`)
	if w2.Uid != 7 || w2.MsgType != 2 {
		t.Fatalf("%+v", w2)
	}
}
