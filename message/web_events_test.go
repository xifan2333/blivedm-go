package message

import "testing"

func TestUserToastV2Parse(t *testing.T) {
	body := []byte(`{
		"cmd":"USER_TOAST_MSG_V2",
		"data":{
			"sender_uinfo":{"uid":42,"base":{"name":"舰长甲"}},
			"guard_info":{"guard_level":3,"start_time":100,"end_time":100},
			"pay_info":{"num":1,"price":198000,"unit":"月"},
			"gift_info":{"gift_id":10003},
			"option":{"source":0},
			"toast_msg":"<%舰长甲%> 开通了舰长"
		}
	}`)
	u := new(UserToastV2)
	u.Parse(body)
	if u.Uid != 42 || u.Username != "舰长甲" || u.GuardLevel != 3 || u.Price != 198000 || u.Source != 0 {
		t.Fatalf("%+v", u)
	}
	if u.ToastMsg == "" || u.Unit != "月" {
		t.Fatalf("%+v", u)
	}
}

func TestSuperChatDeleteParse(t *testing.T) {
	body := []byte(`{"cmd":"SUPER_CHAT_MESSAGE_DELETE","data":{"ids":[1,2,3]}}`)
	s := new(SuperChatDelete)
	s.Parse(body)
	if len(s.Ids) != 3 || s.Ids[0] != 1 || s.Ids[2] != 3 {
		t.Fatalf("%+v", s)
	}
}

func TestLikeInfoV3ClickParse(t *testing.T) {
	body := []byte(`{"cmd":"LIKE_INFO_V3_CLICK","data":{"uid":9,"uname":"点赞侠","like_text":"为主播点赞了"}}`)
	l := new(LikeInfoV3Click)
	l.Parse(body)
	if l.Uid != 9 || l.Uname != "点赞侠" || l.LikeText == "" {
		t.Fatalf("%+v", l)
	}
}
