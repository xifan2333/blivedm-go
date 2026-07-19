package main

import (
	"fmt"

	"github.com/xifan2333/blivedm-go/api"
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
	// 弹幕
	c.OnDanmaku(func(danmaku *message.Danmaku) {
		if danmaku.Type == message.EmoticonDanmaku {
			fmt.Printf("[弹幕表情] %s：表情URL： %s\n", danmaku.Sender.Uname, danmaku.Emoticon.Url)
		} else {
			fmt.Printf("[弹幕] %s：%s\n", danmaku.Sender.Uname, danmaku.Content)
		}
	})
	// 醒目留言
	c.OnSuperChat(func(superChat *message.SuperChat) {
		fmt.Printf("[SC|%d元] %s: %s\n", superChat.Price, superChat.UserInfo.Uname, superChat.Message)
	})
	c.OnSuperChatDelete(func(del *message.SuperChatDelete) {
		fmt.Printf("[SC删除] ids=%v\n", del.Ids)
	})
	// 礼物
	c.OnGift(func(gift *message.Gift) {
		if gift.CoinType == "gold" {
			fmt.Printf("[礼物] %s 的 %s %d 个 共%.2f元\n", gift.Uname, gift.GiftName, gift.Num, float64(gift.Num*gift.Price)/1000)
		}
	})
	// 上舰
	c.OnGuardBuy(func(guardBuy *message.GuardBuy) {
		fmt.Printf("[大航海] %s 开通了 %d 等级的大航海，金额 %d 元\n", guardBuy.Username, guardBuy.GuardLevel, guardBuy.Price/1000)
	})
	c.OnUserToastV2(func(u *message.UserToastV2) {
		fmt.Printf("[上舰V2|source=%d] %s level=%d price=%d %s\n", u.Source, u.Username, u.GuardLevel, u.Price, u.ToastMsg)
	})
	// 进房/关注/分享/点赞（INTERACT_WORD_V2）
	c.OnInteractWordV2(func(w *message.InteractWordV2) {
		switch w.MsgType {
		case message.InteractMsgTypeEnter:
			fmt.Printf("[进房] %s face=%s\n", w.Uname, w.Face)
		case message.InteractMsgTypeFollow:
			fmt.Printf("[关注] %s\n", w.Uname)
		case message.InteractMsgTypeShare:
			fmt.Printf("[分享] %s\n", w.Uname)
		case message.InteractMsgTypeLike:
			fmt.Printf("[点赞互动] %s\n", w.Uname)
		default:
			fmt.Printf("[互动:%d] %s\n", w.MsgType, w.Uname)
		}
	})
	// 点赞按钮
	c.OnLikeClick(func(l *message.LikeInfoV3Click) {
		fmt.Printf("[点赞] %s %s\n", l.Uname, l.LikeText)
	})
	// 自定义事件
	c.RegisterCustomEventHandler("STOP_LIVE_ROOM_LIST", func(s string) {
		data := gjson.Get(s, "data").String()
		fmt.Printf("STOP_LIVE_ROOM_LIST: %s\n", data)
	})

	err := c.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("started")
	select {}
}

func sendDanmaku() error {
	dmReq := &api.DanmakuRequest{
		Msg:      "official_13",
		RoomID:   "732",
		Bubble:   "0",
		Color:    "16777215",
		FontSize: "25",
		Mode:     "1",
		DmType:   "1",
	}
	d, err := api.SendDanmaku(dmReq, &api.BiliVerify{
		Csrf:     "",
		SessData: "",
	})
	if err != nil {
		return err
	}
	fmt.Println(d)
	return nil
}
