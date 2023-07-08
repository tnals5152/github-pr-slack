package slack

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
	constant "tnals5152/git-pr-slack/const"
	"tnals5152/git-pr-slack/models"
	"tnals5152/git-pr-slack/utils"

	"github.com/spf13/viper"
)

func UpsertTodayMainMessage(token, channel, text string, repoInfo *models.RepoInfo) (returnMessageID string, err error) {

	lastDate, err := time.Parse("2006-01-02", repoInfo.Message.LastMessageDate)

	if err != nil {
		return
	}
	todayString := time.Now().Format("2006-01-02")
	today, _ := time.Parse("2006-01-02", todayString)

	if !today.After(lastDate) {
		returnMessageID = repoInfo.Message.LastMessageID
		return
	}
	mainMessage := models.Message{
		Channel: channel,
		Text:    text,
	}

	requestBody, err := json.Marshal(mainMessage)

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

	body, err := io.ReadAll(response.Body)

	var messageResponse models.MessageResponse

	err = json.Unmarshal(body, &messageResponse)

	if err != nil {
		return
	}
	repoInfo.Message.LastMessageDate = todayString
	repoInfo.Message.LastMessageID = messageResponse.Ts

	viper.Set("repo-info."+repoInfo.Key, repoInfo)

	viper.WriteConfig()

	returnMessageID = messageResponse.Ts

	return
}
