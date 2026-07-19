package message

import (
	"github.com/xifan2333/blivedm-go/utils"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// UserToastV2 is USER_TOAST_MSG_V2 — richer guard/toast payload used by modern bilibili live.
// Field layout follows xfgryujk/blivedm UserToastV2Message.
type UserToastV2 struct {
	Uid        int    `json:"uid"`
	Username   string `json:"username"`
	GuardLevel int    `json:"guard_level"` // 0 none, 1 总督, 2 提督, 3 舰长
	Num        int    `json:"num"`
	Price      int    `json:"price"` // 金瓜子
	Unit       string `json:"unit"`
	GiftId     int    `json:"gift_id"`
	StartTime  int    `json:"start_time"`
	EndTime    int    `json:"end_time"`
	// Source: 0 paid, 2 gift/transfer. Official UI often shows only source=0.
	Source   int    `json:"source"`
	ToastMsg string `json:"toast_msg"`
	Raw      string `json:"-"`
}

// Parse decodes a full USER_TOAST_MSG_V2 notification body.
func (u *UserToastV2) Parse(data []byte) {
	u.Raw = utils.BytesToString(data)
	root := gjson.ParseBytes(data)
	d := root.Get("data")
	if !d.Exists() {
		d = root
	}

	sender := d.Get("sender_uinfo")
	guard := d.Get("guard_info")
	pay := d.Get("pay_info")
	gift := d.Get("gift_info")
	option := d.Get("option")

	u.Uid = int(sender.Get("uid").Int())
	u.Username = sender.Get("base.name").String()
	if u.Username == "" {
		u.Username = d.Get("username").String()
	}
	u.GuardLevel = int(guard.Get("guard_level").Int())
	u.Num = int(pay.Get("num").Int())
	u.Price = int(pay.Get("price").Int())
	u.Unit = pay.Get("unit").String()
	u.GiftId = int(gift.Get("gift_id").Int())
	u.StartTime = int(guard.Get("start_time").Int())
	u.EndTime = int(guard.Get("end_time").Int())
	u.Source = int(option.Get("source").Int())
	u.ToastMsg = d.Get("toast_msg").String()

	// Fallback flat shape (older dumps)
	if u.Uid == 0 {
		u.Uid = int(d.Get("uid").Int())
	}
	if u.GuardLevel == 0 {
		u.GuardLevel = int(d.Get("guard_level").Int())
	}
	if u.Price == 0 {
		u.Price = int(d.Get("price").Int())
	}
	if u.ToastMsg == "" {
		log.Debug("UserToastV2: empty toast_msg")
	}
}
