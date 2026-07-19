package client

import (
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/xifan2333/blivedm-go/message"
	"github.com/xifan2333/blivedm-go/packet"
	"github.com/xifan2333/blivedm-go/utils"
	log "github.com/sirupsen/logrus"
)

// knownCMD lists CMDs that are expected on the web live danmaku stream.
// Unknown CMDs outside this set are logged at debug level.
// Aligned with common bilibili live CMDs; first-class typed handlers cover the
// subset that xfgryujk/blivedm web BaseHandler implements.
var (
	knownCMD = []string{
		// first-class (typed handlers)
		"DANMU_MSG", "SEND_GIFT", "SUPER_CHAT_MESSAGE", "SUPER_CHAT_MESSAGE_DELETE",
		"GUARD_BUY", "USER_TOAST_MSG", "USER_TOAST_MSG_V2",
		"INTERACT_WORD", "INTERACT_WORD_V2", "LIKE_INFO_V3_CLICK",
		"LIVE", "PREPARING",
		// known / ignored without typed model
		"HOT_RANK_SETTLEMENT", "DANMU_GIFT_LOTTERY_START", "WELCOME_GUARD", "PK_PROCESS",
		"PK_BATTLE_PRO_TYPE", "MATCH_TEAM_GIFT_RANK", "PK_BATTLE_CRIT", "LUCK_GIFT_AWARD_USER",
		"SCORE_CARD", "ONLINE_RANK_V2", "PK_BATTLE_SPECIAL_GIFT", "SEND_TOP", "SUPER_CHAT_MESSAGE_JPN",
		"ANIMATION", "GUARD_LOTTERY_START", "WEEK_STAR_CLOCK", "WELCOME", "WIN_ACTIVITY",
		"ROOM_KICKOUT", "CHANGE_ROOM_INFO", "ROOM_SKIN_MSG", "ROOM_BLOCK_MSG", "SUPER_CHAT_ENTRANCE",
		"PK_BATTLE_RANK_CHANGE", "ROOM_LOCK", "TV_END", "PK_PRE", "ROOM_SILENT_OFF",
		"ANCHOR_LOT_START", "ROOM_BOX_USER", "ONLINE_RANK_TOP3", "WIDGET_BANNER", "PK_BATTLE_START",
		"ACTIVITY_MATCH_GIFT", "PK_AGAIN", "PK_MATCH", "RAFFLE_START", "WISH_BOTTLE",
		"GUARD_ACHIEVEMENT_ROOM", "ONLINE_RANK_COUNT", "COMMON_NOTICE_DANMAKU", "LOL_ACTIVITY",
		"HOT_RANK_CHANGED", "ROOM_BLOCK_INTO", "ROOM_LIMIT", "PANEL", "RAFFLE_END", "ENTRY_EFFECT",
		"STOP_LIVE_ROOM_LIST", "TV_START", "WATCH_LPL_EXPIRED", "PK_BATTLE_PRE", "BOX_ACTIVITY_START",
		"PK_MIC_END", "LIVE_INTERACTIVE_GAME", "ROOM_BANNER", "PK_BATTLE_GIFT",
		"MESSAGEBOX_USER_GAIN_MEDAL", "LITTLE_TIPS", "HOUR_RANK_AWARDS", "NOTICE_MSG",
		"ROOM_REAL_TIME_MESSAGE_UPDATE", "ANCHOR_LOT_END", "ROOM_CHANGE", "room_admin_entrance",
		"CHASE_FRAME_SWITCH", "DANMU_GIFT_LOTTERY_AWARD", "PK_BATTLE_VOTES_ADD", "PK_BATTLE_END",
		"CUT_OFF", "PK_BATTLE_PROCESS", "PK_BATTLE_SETTLE_USER", "ANCHOR_LOT_AWARD", "WIN_ACTIVITY_USER",
		"VOICE_JOIN_STATUS", "DANMU_GIFT_LOTTERY_END", "ROOM_RANK", "ACTIVITY_BANNER_UPDATE_V2",
		"SPECIAL_GIFT", "ROOM_SILENT_ON", "WARNING", "ROOM_ADMINS", "COMBO_SEND",
		"HOT_RANK_SETTLEMENT_V2", "ANCHOR_LOT_CHECKSTATUS", "HOT_RANK_CHANGED_V2", "PK_END", "PK_SETTLE",
		"ROOM_REFRESH", "PK_START", "COMBO_END", "PK_LOTTERY_START", "GUARD_WINDOWS_OPEN",
		"REENTER_LIVE_ROOM", "MESSAGEBOX_USER_MEDAL_CHANGE", "MESSAGEBOX_USER_MEDAL_COMPENSATION",
		"LITTLE_MESSAGE_BOX", "PK_BATTLE_PRE_NEW", "PK_BATTLE_START_NEW", "PK_BATTLE_PROCESS_NEW",
		"PK_BATTLE_FINAL_PROCESS", "PK_BATTLE_SETTLE_V2", "PK_BATTLE_SETTLE_NEW", "PK_BATTLE_PUNISH_END",
		"PK_BATTLE_VIDEO_PUNISH_BEGIN", "PK_BATTLE_VIDEO_PUNISH_END", "ENTRY_EFFECT_MUST_RECEIVE",
		"SUPER_CHAT_AUDIT", "VIDEO_CONNECTION_JOIN_START", "VIDEO_CONNECTION_JOIN_END", "VIDEO_CONNECTION_MSG",
		"VTR_GIFT_LOTTERY", "RED_POCKET_START", "FULL_SCREEN_SPECIAL_EFFECT", "POPULARITY_RED_POCKET_START",
		"POPULARITY_RED_POCKET_WINNER_LIST", "USER_PANEL_RED_ALARM", "SHOPPING_CART_SHOW",
		"THERMAL_STORM_DANMU_BEGIN", "THERMAL_STORM_DANMU_UPDATE", "THERMAL_STORM_DANMU_CANCEL",
		"THERMAL_STORM_DANMU_OVER", "MILESTONE_UPDATE_EVENT", "WEB_REPORT_CONTROL", "DANMU_TAG_CHANGE",
		"RANK_REM", "LIVE_PLAYER_LOG_RECYCLE", "LIVE_INTERNAL_ROOM_LOGIN", "LIVE_OPEN_PLATFORM_GAME",
		"WATCHED_CHANGE", "DANMU_AGGREGATION", "POPULARITY_RED_POCKET_NEW", "POPULAR_RANK_CHANGED",
		"DM_INTERACTION", "LIKE_INFO_V3_UPDATE", "HOT_ROOM_NOTIFY", "PLAY_TAG", "OTHER_SLICE_LOADING_RESULT",
		// open-live cmds (not typed here; use RegisterCustomEventHandler)
		"LIVE_OPEN_PLATFORM_DM", "LIVE_OPEN_PLATFORM_DM_MIRROR", "LIVE_OPEN_PLATFORM_SEND_GIFT",
		"LIVE_OPEN_PLATFORM_GUARD", "LIVE_OPEN_PLATFORM_SUPER_CHAT", "LIVE_OPEN_PLATFORM_SUPER_CHAT_DEL",
		"LIVE_OPEN_PLATFORM_LIKE", "LIVE_OPEN_PLATFORM_LIVE_ROOM_ENTER", "LIVE_OPEN_PLATFORM_LIVE_START",
		"LIVE_OPEN_PLATFORM_LIVE_END", "LIVE_OPEN_PLATFORM_INTERACTION_END",
	}
	knownCMDMap map[string]int
	cmdReg      = regexp.MustCompile(`"cmd":"([^"]+)"`)
)

type eventHandlers struct {
	danmakuMessageHandlers []func(*message.Danmaku)
	superChatHandlers      []func(*message.SuperChat)
	superChatDeleteHandlers []func(*message.SuperChatDelete)
	giftHandlers           []func(*message.Gift)
	guardBuyHandlers       []func(*message.GuardBuy)
	liveStartHandlers      []func(start *message.LiveStart)
	liveStopHandlers       []func(start *message.LiveStop)
	userToastHandlers      []func(*message.UserToast)
	userToastV2Handlers    []func(*message.UserToastV2)
	interactWordHandlers   []func(*message.InteractWord)
	interactWordV2Handlers []func(*message.InteractWordV2)
	likeClickHandlers      []func(*message.LikeInfoV3Click)
}

type customEventHandlers map[string]func(s string)

func init() {
	knownCMDMap = make(map[string]int, len(knownCMD))
	for _, c := range knownCMD {
		knownCMDMap[c] = 0
	}
}

// RegisterCustomEventHandler 注册自定义事件处理器。
// 若 cmd 与内置 typed handler 相同，自定义 handler 会覆盖内置逻辑。
func (c *Client) RegisterCustomEventHandler(cmd string, handler func(s string)) {
	(*c.customEventHandlers)[cmd] = handler
}

// OnDanmaku 弹幕 DANMU_MSG
func (c *Client) OnDanmaku(f func(*message.Danmaku)) {
	c.eventHandlers.danmakuMessageHandlers = append(c.eventHandlers.danmakuMessageHandlers, f)
}

// OnSuperChat 醒目留言 SUPER_CHAT_MESSAGE
func (c *Client) OnSuperChat(f func(*message.SuperChat)) {
	c.eventHandlers.superChatHandlers = append(c.eventHandlers.superChatHandlers, f)
}

// OnSuperChatDelete 删除醒目留言 SUPER_CHAT_MESSAGE_DELETE
func (c *Client) OnSuperChatDelete(f func(*message.SuperChatDelete)) {
	c.eventHandlers.superChatDeleteHandlers = append(c.eventHandlers.superChatDeleteHandlers, f)
}

// OnGift 礼物 SEND_GIFT
func (c *Client) OnGift(f func(gift *message.Gift)) {
	c.eventHandlers.giftHandlers = append(c.eventHandlers.giftHandlers, f)
}

// OnGuardBuy 开通大航海 GUARD_BUY
func (c *Client) OnGuardBuy(f func(*message.GuardBuy)) {
	c.eventHandlers.guardBuyHandlers = append(c.eventHandlers.guardBuyHandlers, f)
}

// OnLiveStart 开播 LIVE
func (c *Client) OnLiveStart(f func(start *message.LiveStart)) {
	c.eventHandlers.liveStartHandlers = append(c.eventHandlers.liveStartHandlers, f)
}

// OnLiveStop 关播 PREPARING
func (c *Client) OnLiveStop(f func(start *message.LiveStop)) {
	c.eventHandlers.liveStopHandlers = append(c.eventHandlers.liveStopHandlers, f)
}

// OnUserToast 旧版 USER_TOAST_MSG
func (c *Client) OnUserToast(f func(*message.UserToast)) {
	c.eventHandlers.userToastHandlers = append(c.eventHandlers.userToastHandlers, f)
}

// OnUserToastV2 新版上舰提示 USER_TOAST_MSG_V2
func (c *Client) OnUserToastV2(f func(*message.UserToastV2)) {
	c.eventHandlers.userToastV2Handlers = append(c.eventHandlers.userToastV2Handlers, f)
}

// OnInteractWord 旧版 INTERACT_WORD（JSON）
func (c *Client) OnInteractWord(f func(*message.InteractWord)) {
	c.eventHandlers.interactWordHandlers = append(c.eventHandlers.interactWordHandlers, f)
}

// OnInteractWordV2 新版互动 INTERACT_WORD_V2（protobuf，含进房/关注/分享/点赞）
func (c *Client) OnInteractWordV2(f func(*message.InteractWordV2)) {
	c.eventHandlers.interactWordV2Handlers = append(c.eventHandlers.interactWordV2Handlers, f)
}

// OnLikeClick 点赞按钮 LIKE_INFO_V3_CLICK
func (c *Client) OnLikeClick(f func(*message.LikeInfoV3Click)) {
	c.eventHandlers.likeClickHandlers = append(c.eventHandlers.likeClickHandlers, f)
}

// Handle 处理一个包
func (c *Client) Handle(p packet.Packet) {
	switch p.Operation {
	case packet.Notification:
		cmd := parseCmd(p.Body)
		sb := utils.BytesToString(p.Body)
		// 新的弹幕 cmd 可能带参数，如 DANMU_MSG:4:0:2:2:2:0
		if ind := strings.Index(cmd, ":"); ind >= 0 {
			cmd = cmd[:ind]
		}
		// 自定义 handler 覆盖内置
		if f, ok := (*c.customEventHandlers)[cmd]; ok {
			go cover(func() { f(sb) })
			return
		}
		switch cmd {
		case "DANMU_MSG":
			d := new(message.Danmaku)
			d.Parse(p.Body)
			for _, fn := range c.eventHandlers.danmakuMessageHandlers {
				go cover(func() { fn(d) })
			}
		case "SUPER_CHAT_MESSAGE":
			s := new(message.SuperChat)
			s.Parse(p.Body)
			for _, fn := range c.eventHandlers.superChatHandlers {
				go cover(func() { fn(s) })
			}
		case "SUPER_CHAT_MESSAGE_DELETE":
			s := new(message.SuperChatDelete)
			s.Parse(p.Body)
			for _, fn := range c.eventHandlers.superChatDeleteHandlers {
				go cover(func() { fn(s) })
			}
		case "SEND_GIFT":
			g := new(message.Gift)
			g.Parse(p.Body)
			for _, fn := range c.eventHandlers.giftHandlers {
				go cover(func() { fn(g) })
			}
		case "GUARD_BUY":
			g := new(message.GuardBuy)
			g.Parse(p.Body)
			for _, fn := range c.eventHandlers.guardBuyHandlers {
				go cover(func() { fn(g) })
			}
		case "LIVE":
			l := new(message.LiveStart)
			l.Parse(p.Body)
			for _, fn := range c.eventHandlers.liveStartHandlers {
				go cover(func() { fn(l) })
			}
		case "PREPARING":
			l := new(message.LiveStop)
			l.Parse(p.Body)
			for _, fn := range c.eventHandlers.liveStopHandlers {
				go cover(func() { fn(l) })
			}
		case "USER_TOAST_MSG":
			u := new(message.UserToast)
			u.Parse(p.Body)
			for _, fn := range c.eventHandlers.userToastHandlers {
				go cover(func() { fn(u) })
			}
		case "USER_TOAST_MSG_V2":
			u := new(message.UserToastV2)
			u.Parse(p.Body)
			for _, fn := range c.eventHandlers.userToastV2Handlers {
				go cover(func() { fn(u) })
			}
		case "INTERACT_WORD":
			w := message.ParseInteractWordJSON(sb)
			for _, fn := range c.eventHandlers.interactWordHandlers {
				go cover(func() { fn(w) })
			}
		case "INTERACT_WORD_V2":
			w := new(message.InteractWordV2)
			w.Parse(p.Body)
			for _, fn := range c.eventHandlers.interactWordV2Handlers {
				go cover(func() { fn(w) })
			}
		case "LIKE_INFO_V3_CLICK":
			l := new(message.LikeInfoV3Click)
			l.Parse(p.Body)
			for _, fn := range c.eventHandlers.likeClickHandlers {
				go cover(func() { fn(l) })
			}
		default:
			if _, ok := knownCMDMap[cmd]; ok {
				return
			}
			log.Debugf("unknown cmd(%s), body: %s", cmd, p.Body)
		}
	case packet.HeartBeatResponse:
	case packet.RoomEnterResponse:
	default:
		log.WithField("protover", p.ProtocolVersion).
			WithField("data", string(p.Body)).
			Warn("unknown protover")
	}
}

func parseCmd(d []byte) string {
	str := utils.BytesToString(d)
	match := cmdReg.FindStringSubmatch(str)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func cover(f func()) {
	defer func() {
		if pan := recover(); pan != nil {
			log.Errorf("event error: %v\n%s", pan, debug.Stack())
		}
	}()
	f()
}
