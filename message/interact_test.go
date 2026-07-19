package message

import (
	"encoding/base64"
	"testing"

	"github.com/xifan2333/blivedm-go/pb"
	"google.golang.org/protobuf/proto"
)

func TestInteractWordParse(t *testing.T) {
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
	body := []byte(`{"cmd":"INTERACT_WORD_V2","data":{"pb":"` + base64.StdEncoding.EncodeToString(raw) + `"}}`)
	got := new(InteractWord)
	got.Parse(body)
	if got.Uid != 123456 || got.Uname != "测试用户" || got.MsgType != InteractMsgTypeEnter || got.Face == "" {
		t.Fatalf("%+v", got)
	}
}

func TestUserToastParse(t *testing.T) {
	body := []byte(`{"cmd":"USER_TOAST_MSG_V2","data":{"sender_uinfo":{"uid":42,"base":{"name":"舰长甲"}},"guard_info":{"guard_level":3,"start_time":100,"end_time":100},"pay_info":{"num":1,"price":198000,"unit":"月"},"gift_info":{"gift_id":10003},"option":{"source":0},"toast_msg":"开通了舰长"}}`)
	u := new(UserToast)
	u.Parse(body)
	if u.Uid != 42 || u.Username != "舰长甲" || u.Price != 198000 {
		t.Fatalf("%+v", u)
	}
}

func TestSuperChatDeleteParse(t *testing.T) {
	s := new(SuperChatDelete)
	s.Parse([]byte(`{"cmd":"SUPER_CHAT_MESSAGE_DELETE","data":{"ids":[1,2,3]}}`))
	if len(s.Ids) != 3 {
		t.Fatalf("%+v", s)
	}
}
