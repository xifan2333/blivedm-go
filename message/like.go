package message

import (
	"github.com/xifan2333/blivedm-go/utils"
	"github.com/tidwall/gjson"
)

// LikeInfoV3Click is LIKE_INFO_V3_CLICK (web like button).
type LikeInfoV3Click struct {
	Uid       int    `json:"uid"`
	Uname     string `json:"uname"`
	LikeText  string `json:"like_text"`
	LikeIcon  string `json:"like_icon"`
	FansMedal struct {
		TargetID         int64  `json:"target_id"`
		MedalLevel       int    `json:"medal_level"`
		MedalName        string `json:"medal_name"`
		MedalColor       int    `json:"medal_color"`
		MedalColorStart  int    `json:"medal_color_start"`
		MedalColorEnd    int    `json:"medal_color_end"`
		MedalColorBorder int    `json:"medal_color_border"`
		IsLighted        int    `json:"is_lighted"`
		GuardLevel       int    `json:"guard_level"`
	} `json:"fans_medal"`
	Raw string `json:"-"`
}

// Parse decodes a full LIKE_INFO_V3_CLICK notification body.
func (l *LikeInfoV3Click) Parse(data []byte) {
	l.Raw = utils.BytesToString(data)
	d := gjson.ParseBytes(data).Get("data")
	if !d.Exists() {
		d = gjson.ParseBytes(data)
	}
	l.Uid = int(d.Get("uid").Int())
	l.Uname = d.Get("uname").String()
	l.LikeText = d.Get("like_text").String()
	l.LikeIcon = d.Get("like_icon").String()
	fm := d.Get("fans_medal")
	l.FansMedal.TargetID = fm.Get("target_id").Int()
	l.FansMedal.MedalLevel = int(fm.Get("medal_level").Int())
	l.FansMedal.MedalName = fm.Get("medal_name").String()
	l.FansMedal.MedalColor = int(fm.Get("medal_color").Int())
	l.FansMedal.MedalColorStart = int(fm.Get("medal_color_start").Int())
	l.FansMedal.MedalColorEnd = int(fm.Get("medal_color_end").Int())
	l.FansMedal.MedalColorBorder = int(fm.Get("medal_color_border").Int())
	l.FansMedal.IsLighted = int(fm.Get("is_lighted").Int())
	l.FansMedal.GuardLevel = int(fm.Get("guard_level").Int())
}
