package message

import (
	"github.com/xifan2333/blivedm-go/utils"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// SuperChatDelete is SUPER_CHAT_MESSAGE_DELETE.
type SuperChatDelete struct {
	Ids []int  `json:"ids"`
	Raw string `json:"-"`
}

// Parse decodes a full SUPER_CHAT_MESSAGE_DELETE notification body.
func (s *SuperChatDelete) Parse(data []byte) {
	s.Raw = utils.BytesToString(data)
	root := gjson.ParseBytes(data)
	arr := root.Get("data.ids")
	if !arr.Exists() {
		arr = root.Get("ids")
	}
	if !arr.IsArray() {
		log.Error("parse SuperChatDelete failed: ids missing")
		return
	}
	s.Ids = make([]int, 0, len(arr.Array()))
	for _, v := range arr.Array() {
		s.Ids = append(s.Ids, int(v.Int()))
	}
}
