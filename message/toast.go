package message

import (
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// UserToast is USER_TOAST_MSG_V2.
type UserToast struct {
	Uid        int    `json:"uid"`
	Username   string `json:"username"`
	GuardLevel int    `json:"guard_level"` // 0 none, 1 总督, 2 提督, 3 舰长
	Num        int    `json:"num"`
	Price      int    `json:"price"` // 金瓜子
	Unit       string `json:"unit"`
	GiftId     int    `json:"gift_id"`
	StartTime  int    `json:"start_time"`
	EndTime    int    `json:"end_time"`
	// Source: 0 paid, 2 gift/transfer.
	Source   int    `json:"source"`
	ToastMsg string `json:"toast_msg"`
}

func (u *UserToast) Parse(data []byte) {
	d := gjson.ParseBytes(data).Get("data")
	if !d.Exists() {
		log.Error("parse UserToast failed")
		return
	}
	u.Uid = int(d.Get("sender_uinfo.uid").Int())
	u.Username = d.Get("sender_uinfo.base.name").String()
	u.GuardLevel = int(d.Get("guard_info.guard_level").Int())
	u.Num = int(d.Get("pay_info.num").Int())
	u.Price = int(d.Get("pay_info.price").Int())
	u.Unit = d.Get("pay_info.unit").String()
	u.GiftId = int(d.Get("gift_info.gift_id").Int())
	u.StartTime = int(d.Get("guard_info.start_time").Int())
	u.EndTime = int(d.Get("guard_info.end_time").Int())
	u.Source = int(d.Get("option.source").Int())
	u.ToastMsg = d.Get("toast_msg").String()
}
