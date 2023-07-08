package slack

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	constant "tnals5152/git-pr-slack/const"
	"tnals5152/git-pr-slack/models"
	"tnals5152/git-pr-slack/utils"

	"github.com/spf13/viper"
)

func SendTheadMessage(token, channel, parentID, text string) {

	message := models.Message{
		Channel: channel,
		TheadTs: parentID,
		Text:    text,
	}
	requestBody, err := json.Marshal(message)

	if err != nil {
		return
	}

	buff := bytes.NewBuffer(requestBody)

	request, err := http.NewRequest(http.MethodPost, constant.SLACK_POST_MESSAGE, buff)

	if err != nil {
		return
	}

	request.Header.Add("Content-type", "application/json")
	request.Header.Add("Authorization", "Bearer "+token)

	client := http.Client{
		Timeout: utils.GetTimeout(constant.TIMEOUT_HTTP),
	}
	defer client.CloseIdleConnections()

	response, err := client.Do(request)

	if err != nil {
		return
	}

	defer response.Body.Close()

}

func CreateReviewersMessage(reviewers []string, login string) (text string) {
	reviewersMap := viper.GetStringMapString(constant.REVIEWERS)

	if len(reviewers) == 0 {
		name, ok := reviewersMap[login]

		if !ok {
			return
		}

		text += "<@" + name + "> 리뷰어를 지정하세요."
		return
	}

	for _, reviewer := range reviewers {
		name, ok := reviewersMap[reviewer]

		if !ok {
			continue
		}

		text += "<@" + name + "> "
	}
	return
}

func DeleteTheadMessage(token, channel string, repoInfo *models.RepoInfo) (err error) {

	request, err := http.NewRequest(
		http.MethodGet,
		constant.SLACK_REPLY_URL+"?channel="+channel+"&ts="+repoInfo.Message.LastMessageID,
		nil,
	)

	if err != nil {
		return
	}

	request.Header.Add("Content-type", "application/x-www-form-urlencoded")
	request.Header.Add("Authorization", "Bearer "+token)

	client := http.Client{
		Timeout: utils.GetTimeout(constant.TIMEOUT_HTTP),
	}
	defer client.CloseIdleConnections()

	response, err := client.Do(request)

	if err != nil {
		return
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	var messageResponse models.ResponseOkMessage

	err = json.Unmarshal(body, &messageResponse)

	if err != nil {
		return
	}

	for _, message := range messageResponse.Messages {
		if message.TheadTs == message.Ts || message.TheadTs == "" {
			continue
		}
		DeleteMessage(token, channel, message.Ts)
	}
	return
}

func DeleteMessage(token, channel, ts string) {
	request, err := http.NewRequest(
		http.MethodPost,
		constant.SLACK_DELETE_URL+"?channel="+channel+"&ts="+ts,
		nil,
	)

	if err != nil {
		return
	}

	request.Header.Add("Content-type", "application/x-www-form-urlencoded")
	request.Header.Add("Authorization", "Bearer "+token)

	client := http.Client{
		Timeout: utils.GetTimeout(constant.TIMEOUT_HTTP),
	}
	defer client.CloseIdleConnections()

	response, err := client.Do(request)

	if err != nil {
		return
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	var messageResponse any

	err = json.Unmarshal(body, &messageResponse)

	if err != nil {
		return
	}
}
