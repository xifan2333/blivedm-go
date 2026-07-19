# blivedm-go

bilibili 直播弹幕 golang 库（fork：对齐当前 web 协议）

## 安装

```shell
go get github.com/xifan2333/blivedm-go@v1.7.1
```

## 一等事件

| CMD | API |
| --- | --- |
| `DANMU_MSG` | `OnDanmaku` |
| `SEND_GIFT` | `OnGift` |
| `SUPER_CHAT_MESSAGE` | `OnSuperChat` |
| `SUPER_CHAT_MESSAGE_DELETE` | `OnSuperChatDelete` |
| `GUARD_BUY` | `OnGuardBuy` |
| `USER_TOAST_MSG_V2` | `OnUserToast` |
| `INTERACT_WORD_V2` | `OnInteractWord`（`data.pb` protobuf） |
| `LIVE` / `PREPARING` | `OnLiveStart` / `OnLiveStop` |

其它 CMD：`RegisterCustomEventHandler`。

旧 `INTERACT_WORD` / `USER_TOAST_MSG` 不再解析，仅静默忽略。

## 快速开始

见 `example/main.go`。

### Cookie

必要 cookie：`buvid3`, `SESSDATA`, `bili_jct`。

## 参考

协议消息模型对齐 [xfgryujk/blivedm](https://github.com/xfgryujk/blivedm) web 侧当前实现。
