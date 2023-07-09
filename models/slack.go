package models

import (
	"strconv"
	"strings"
	"time"
	constant "tnals5152/git-pr-slack/const"
	"tnals5152/git-pr-slack/utils"

	"github.com/spf13/viper"
)

type MessageResponse struct {
	Ts string `json:"ts"`
}

type Message struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
	Ts      string `json:"ts"`
	TheadTs string `json:"thread_ts"` // parent message ID
	User    string `json:"user"`      // message를 쓴 user의 ID
}

type ResponseOkMessage struct {
	Messages []*Message `json:"messages"`
}

func (m *Message) ToDelete() bool {

	botUserID := viper.GetString(constant.SLACK_BOT_USER_ID)
	// 현재 메시지가 부모 메시지거나, 아래 쓰레드 메시지가 없거나, 해당 메시지의 작성자가 해당 봇이 아니면 메시지 삭제하지 않는다.
	if m.TheadTs == m.Ts ||
		m.TheadTs == "" ||
		m.User != botUserID {
		return false
	}
	// 리뷰어 지정 메시지면 삭제하지 않는다.
	return !strings.ContainsAny(m.Text, constant.REQUEST_REVIEWER)
}

// 리뷰어 지정 메시지 이후 처음에만 바로 보내지고 나머지는 3일이 지난 후 다시 보낼 수 있도록 한다.
func (responseMessage *ResponseOkMessage) HasRequestReviewAfterAssign(PRUrl string) bool {
	botUserID := viper.GetString(constant.SLACK_BOT_USER_ID)

	var check bool

	for _, message := range responseMessage.Messages {
		if message.TheadTs == message.Ts || message.TheadTs == "" || message.User != botUserID {
			continue
		}
		isContain := strings.Contains(message.Text, constant.REQUEST_REVIEW)
		if !isContain { //&&
			//!strings.ContainsAny(message.Text, PRUrl) {
			continue
		}
		check = true

		unixTimeSlice := strings.Split(message.Ts, ".")

		var sec, nsec int64

		if len(unixTimeSlice) > 1 {
			sec, _ = strconv.ParseInt(unixTimeSlice[0], 10, 64)
		}

		if len(unixTimeSlice) > 2 {
			nsec, _ = strconv.ParseInt(unixTimeSlice[1], 10, 64)
		}
		messageTime := time.Unix(sec, nsec)
		reviewDay := viper.GetInt(constant.REVIEW_DAY)

		if reviewDay == 0 {
			reviewDay = 3
		}
		isAfterDay := utils.IsAfterDay(messageTime, reviewDay)

		if isAfterDay {
			return isAfterDay
		}
	}
	if !check {
		return true
	}
	return false
}
