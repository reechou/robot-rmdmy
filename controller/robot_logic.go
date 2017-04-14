package controller

import (
	"fmt"
	"strings"
	"strconv"
	"time"
	"math/rand"
	
	"github.com/reechou/holmes"
	"github.com/reechou/robot-rmdmy/robot_proto"
)

const (
	RMDMY_NAME = "人民的名义"
	RMDMY_PREFIX = "r"
)

func (self *Logic) HandleReceiveMsg(msg *robot_proto.ReceiveMsgInfo) {
	holmes.Debug("receive robot msg: %v", msg)
	switch msg.BaseInfo.ReceiveEvent {
	case robot_proto.RECEIVE_EVENT_MSG:
		self.handleMsg(msg)
	}
}

func (self *Logic) handleMsg(msg *robot_proto.ReceiveMsgInfo) {
	switch msg.MsgType {
	case robot_proto.RECEIVE_MSG_TYPE_TEXT:
		msgStr := strings.Replace(msg.Msg, " ", "", -1)
		if strings.HasPrefix(msgStr, RMDMY_PREFIX) {
			url := self.getVideo(msgStr)
			if url != "" {
				url = fmt.Sprintf("%s \n\n链接有可能会失效,失效后请重新获取", url)
				self.sendMsg(msg, []MsgInfo{MsgInfo{MsgType: robot_proto.RECEIVE_MSG_TYPE_TEXT, Msg: url}})
			} else {
				holmes.Debug("get req[%s] url == nil", msgStr)
			}
		} else if strings.Contains(msg.Msg, RMDMY_NAME) {
			offset := rand.Intn(len(self.cfg.WxUrlList))
			self.sendMsg(msg, []MsgInfo{
				MsgInfo{MsgType: robot_proto.RECEIVE_MSG_TYPE_TEXT, Msg: "亲，福利就要分享，请发送以下内容加上配图到朋友圈，然后截图给我，我就发你人民的名义全集，决不食言！"},
				MsgInfo{MsgType: robot_proto.RECEIVE_MSG_TYPE_TEXT, Msg: "福利福利，人民的名义55集全部出来，亲测有效，要看的赶紧加她，然后发送：人民的名义，赶紧收藏起来，手慢就被和谐了"},
				MsgInfo{MsgType: robot_proto.RECEIVE_MSG_TYPE_IMG, Msg: self.cfg.WxUrlList[offset]},
			})
		}
	case robot_proto.RECEIVE_MSG_TYPE_IMG:
		self.sendMsg(msg, []MsgInfo{
			MsgInfo{
				MsgType: robot_proto.RECEIVE_MSG_TYPE_TEXT,
				Msg: "感谢亲的分享，你可以把要看的集数发给我，如第55集，那就发我：r55，我就会把第55集视频地址发给你，你点击打开就可以看了，不过注意这个链接是有有效期的，如果出现打不开了，可以重新再发我一遍，拿到新的地址来播放可以试下哦。\n输入：r55",
			},
		})
	}
}

func (self *Logic) getVideo(msg string) string {
	msg = strings.Replace(msg, RMDMY_PREFIX, "", -1)
	vid, err := strconv.Atoi(msg)
	if err != nil {
		holmes.Error("strconv msg[%s] error: %v", err)
		return ""
	}
	return self.cw.GetVideoUrl(vid)
}

type MsgInfo struct {
	MsgType string
	Msg string
}
func (self *Logic) sendMsg(msg *robot_proto.ReceiveMsgInfo, msgList []MsgInfo) {
	var sendReq robot_proto.SendMsgInfo
	for _, v := range msgList {
		sendReq.SendMsgs = append(sendReq.SendMsgs, robot_proto.SendBaseInfo{
			WechatNick: msg.BaseInfo.WechatNick,
			ChatType:   msg.BaseInfo.FromType,
			UserName:   msg.BaseInfo.FromUserName,
			NickName:   msg.BaseInfo.FromNickName,
			MsgType:    v.MsgType,
			Msg:        v.Msg,
		})
	}
	err := self.robotExt.SendMsgs("", &sendReq)
	if err != nil {
		holmes.Error("send msg[%v] error: %v", sendReq, err)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
