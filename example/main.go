package main

import (
	"fmt"

	"github.com/xifan2333/blivedm-go/client"
	"github.com/xifan2333/blivedm-go/message"
	_ "github.com/xifan2333/blivedm-go/utils"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func main() {
	log.SetLevel(log.DebugLevel)
	c := client.NewClient(23943357)
	c.SetCookie("")

	c.OnDanmaku(func(danmaku *message.Danmaku) {
		if danmaku.Type == message.EmoticonDanmaku {
			fmt.Printf("[弹幕表情] %s：%s\n", danmaku.Sender.Uname, danmaku.Emoticon.Url)
		} else {
			fmt.Printf("[弹幕] %s：%s\n", danmaku.Sender.Uname, danmaku.Content)
		}
	})
	c.OnSuperChat(func(superChat *message.SuperChat) {
		fmt.Printf("[SC|%d元] %s: %s\n", superChat.Price, superChat.UserInfo.Uname, superChat.Message)
	})
	c.OnSuperChatDelete(func(del *message.SuperChatDelete) {
		fmt.Printf("[SC删除] ids=%v\n", del.Ids)
	})
	c.OnGift(func(gift *message.Gift) {
		if gift.CoinType == "gold" {
			fmt.Printf("[礼物] %s 的 %s %d 个 共%.2f元\n", gift.Uname, gift.GiftName, gift.Num, float64(gift.Num*gift.Price)/1000)
		}
	})
	c.OnGuardBuy(func(guardBuy *message.GuardBuy) {
		fmt.Printf("[大航海] %s level=%d price=%d\n", guardBuy.Username, guardBuy.GuardLevel, guardBuy.Price/1000)
	})
	c.OnUserToast(func(u *message.UserToast) {
		fmt.Printf("[上舰toast] %s source=%d %s\n", u.Username, u.Source, u.ToastMsg)
	})
	c.OnInteractWord(func(w *message.InteractWord) {
		switch w.MsgType {
		case message.InteractMsgTypeEnter:
			fmt.Printf("[进房] %s %s\n", w.Uname, w.Face)
		case message.InteractMsgTypeFollow:
			fmt.Printf("[关注] %s\n", w.Uname)
		case message.InteractMsgTypeLike:
			fmt.Printf("[点赞] %s\n", w.Uname)
		default:
			fmt.Printf("[互动:%d] %s\n", w.MsgType, w.Uname)
		}
	})
	c.RegisterCustomEventHandler("STOP_LIVE_ROOM_LIST", func(s string) {
		fmt.Printf("STOP_LIVE_ROOM_LIST: %s\n", gjson.Get(s, "data").String())
	})

	if err := c.Start(); err != nil {
		log.Fatal(err)
	}
	log.Println("started")
	select {}
}
