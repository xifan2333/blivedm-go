package message

import (
	"encoding/base64"

	"github.com/xifan2333/blivedm-go/pb"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"google.golang.org/protobuf/proto"
)

// Parse INTERACT_WORD_V2 body: {"cmd":"INTERACT_WORD_V2","data":{"pb":"..."}}
func (i *InteractWord) Parse(data []byte) {
	pbB64 := gjson.GetBytes(data, "data.pb").String()
	if pbB64 == "" {
		log.Error("parse InteractWord failed: empty data.pb")
		return
	}
	raw, err := base64.StdEncoding.DecodeString(pbB64)
	if err != nil {
		log.Error("parse InteractWord base64 failed")
		return
	}
	msg := new(pb.InteractWordV2)
	if err := proto.Unmarshal(raw, msg); err != nil {
		log.Error("parse InteractWord protobuf failed")
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
