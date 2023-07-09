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

func CreateReviewersMessage(reviewers map[string]string, login string) (text string, reviewerRequest bool) {
	reviewersMap := viper.GetStringMapString(constant.REVIEWERS)

	if len(reviewers) == 0 {
		name, ok := reviewersMap[login]

		if !ok {
			return
		}

		text += "<@" + name + "> " + constant.REQUEST_REVIEWER
		reviewerRequest = true
		return
	}

	text = constant.REQUEST_REVIEW + "\n"

	for reviewerName, state := range reviewers {
		name, ok := reviewersMap[reviewerName]

		if !ok {
			continue
		}

		if state == constant.COMMENTED {
			text += "APPROVE 요청 - <@" + name + ">\n"
		} else if state == constant.NOTHING {
			text += "PR 리뷰 요청 - <@" + name + ">\n"
		}

	}
	return
}

func DeleteTheadMessage(token, channel string, repoInfo *models.RepoInfo) (messageResponse models.ResponseOkMessage, err error) {

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

	var bodyAny any

	err = json.Unmarshal(body, &messageResponse)
	err = json.Unmarshal(body, &bodyAny)

	if err != nil {
		return
	}

	for _, message := range messageResponse.Messages {
		if message.ToDelete() {
			DeleteMessage(token, channel, message.Ts)
		}
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
