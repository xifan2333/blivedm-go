package message

import (
	"encoding/base64"

	"github.com/xifan2333/blivedm-go/pb"
	"github.com/xifan2333/blivedm-go/utils"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"google.golang.org/protobuf/proto"
)

// InteractWord msg_type values (shared by INTERACT_WORD / INTERACT_WORD_V2).
const (
	InteractMsgTypeEnter         = 1 // 进入直播间
	InteractMsgTypeFollow        = 2 // 关注
	InteractMsgTypeShare         = 3 // 分享
	InteractMsgTypeSpecialFollow = 4 // 特别关注
	InteractMsgTypeMutualFollow  = 5 // 互粉
	InteractMsgTypeLike          = 6 // 点赞
)

// InteractWordV2 is the decoded INTERACT_WORD_V2 event.
// Prefer Face from protobuf when present; Uid/Uname/MsgType/Timestamp always filled when decode succeeds.
type InteractWordV2 struct {
	Uid       int    `json:"uid"`
	Uname     string `json:"uname"`
	Face      string `json:"face"`
	MsgType   int    `json:"msg_type"`
	Timestamp int    `json:"timestamp"`
	// Raw is the original notification JSON body.
	Raw string `json:"-"`
	// Pb holds the decoded protobuf when available.
	Pb *pb.InteractWordV2 `json:"-"`
}

// Parse decodes INTERACT_WORD_V2 notification body.
// Payload shape: {"cmd":"INTERACT_WORD_V2","data":{"pb":"<base64 protobuf>"}}
func (i *InteractWordV2) Parse(data []byte) {
	i.Raw = utils.BytesToString(data)
	pbB64 := gjson.GetBytes(data, "data.pb").String()
	if pbB64 == "" {
		// Fallback: some gateways may still embed flat JSON fields under data.
		i.Uid = int(gjson.GetBytes(data, "data.uid").Int())
		i.Uname = gjson.GetBytes(data, "data.uname").String()
		i.MsgType = int(gjson.GetBytes(data, "data.msg_type").Int())
		i.Timestamp = int(gjson.GetBytes(data, "data.timestamp").Int())
		i.Face = gjson.GetBytes(data, "data.uinfo.base.face").String()
		if i.Face == "" {
			i.Face = gjson.GetBytes(data, "data.face").String()
		}
		return
	}

	raw, err := base64.StdEncoding.DecodeString(pbB64)
	if err != nil {
		// Try raw URL encoding just in case.
		raw, err = base64.RawStdEncoding.DecodeString(pbB64)
		if err != nil {
			log.Error("InteractWordV2 base64 decode failed:", err)
			return
		}
	}

	msg := new(pb.InteractWordV2)
	if err := proto.Unmarshal(raw, msg); err != nil {
		log.Error("InteractWordV2 protobuf unmarshal failed:", err)
		return
	}
	i.Pb = msg
	i.Uid = int(msg.GetUid())
	i.Uname = msg.GetUname()
	i.MsgType = int(msg.GetMsgType())
	i.Timestamp = int(msg.GetTimestamp())
	if uinfo := msg.GetUinfo(); uinfo != nil {
		if base := uinfo.GetBase(); base != nil {
			i.Face = base.GetFace()
		}
	}
}

// IsEnter reports whether this interaction is a room-enter event.
func (i *InteractWordV2) IsEnter() bool {
	return i.MsgType == InteractMsgTypeEnter
}

// IsLike reports whether this interaction is a like event.
func (i *InteractWordV2) IsLike() bool {
	return i.MsgType == InteractMsgTypeLike
}

// ParseInteractWordJSON decodes legacy INTERACT_WORD JSON data object (or full packet).
func ParseInteractWordJSON(data string) *InteractWord {
	w := new(InteractWord)
	root := gjson.Parse(data)
	// Accept either full packet or data object.
	src := root
	if root.Get("data").Exists() {
		src = root.Get("data")
	}
	w.Uid = int(src.Get("uid").Int())
	w.Uname = src.Get("uname").String()
	w.MsgType = int(src.Get("msg_type").Int())
	w.Timestamp = int(src.Get("timestamp").Int())
	w.Roomid = int(src.Get("roomid").Int())
	w.Score = src.Get("score").Int()
	w.TriggerTime = src.Get("trigger_time").Int()
	w.UnameColor = src.Get("uname_color").String()
	w.IsSpread = int(src.Get("is_spread").Int())
	w.SpreadDesc = src.Get("spread_desc").String()
	w.SpreadInfo = src.Get("spread_info").String()
	w.TailIcon = int(src.Get("tail_icon").Int())
	w.Dmscore = int(src.Get("dmscore").Int())
	return w
}
