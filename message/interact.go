package message

import (
	"encoding/base64"

	"github.com/xifan2333/blivedm-go/pb"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"google.golang.org/protobuf/proto"
)

// InteractWord msg_type: 1进入 2关注 3分享 4特别关注 5互粉 6点赞
const (
	InteractMsgTypeEnter         = 1
	InteractMsgTypeFollow        = 2
	InteractMsgTypeShare         = 3
	InteractMsgTypeSpecialFollow = 4
	InteractMsgTypeMutualFollow  = 5
	InteractMsgTypeLike          = 6
)

// InteractWord 互动消息（INTERACT_WORD_V2，data.pb）
type InteractWord struct {
	Uid       int    `json:"uid"`
	Uname     string `json:"uname"`
	Face      string `json:"face"`
	MsgType   int    `json:"msg_type"`
	Timestamp int    `json:"timestamp"`
}

func (i *InteractWord) Parse(data []byte) {
	pbB64 := gjson.GetBytes(data, "data.pb").String()
	if pbB64 == "" {
		log.Error("parse interact word failed")
		return
	}
	raw, err := base64.StdEncoding.DecodeString(pbB64)
	if err != nil {
		log.Error("parse interact word failed")
		return
	}
	msg := new(pb.InteractWordV2)
	if err := proto.Unmarshal(raw, msg); err != nil {
		log.Error("parse interact word failed")
		return
	}
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
