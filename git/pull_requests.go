package git

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	constant "tnals5152/git-pr-slack/const"
	"tnals5152/git-pr-slack/models"
	"tnals5152/git-pr-slack/slack"
	"tnals5152/git-pr-slack/utils"
	util_error "tnals5152/git-pr-slack/utils/error"

	"github.com/spf13/viper"
)

func GetPRList(gitToken string, repoInfos []*models.RepoInfo) (err error) {
	defer util_error.DeferWrap(err)

	// wg := sync.WaitGroup{}

	for _, repoInfo := range repoInfos {
		// wg.Add(1)
		// go func(repoInfo models.RepoInfo) {
		// defer wg.Done()
		errorSlice := GetAndDoRepo(gitToken, *repoInfo)
		fmt.Println(errorSlice, repoInfo)
		// }(*repoInfo)
	}

	// wg.Wait()

	return
}

func GetAndDoRepo(gitToken string, repoInfo models.RepoInfo) (errorSlice []string) {
	var request *http.Request
	var response *http.Response
	client := http.Client{
		Timeout: utils.GetTimeout(constant.TIMEOUT_HTTP),
	}
	defer client.CloseIdleConnections()

	slackBot := viper.GetString(constant.SLACK_BOT_TOKEN)
	slackChannel := viper.GetString(constant.SLACK_CHANNEL)
	url := constant.PR_LIST_URL

	url = strings.ReplaceAll(url, constant.OWNER_PARAM, repoInfo.Owner)
	url = strings.ReplaceAll(url, constant.REPO_PARAM, repoInfo.Repo)
	request, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		errorSlice = append(errorSlice, err.Error())
		return
	}
	request.Header.Add("Accept", "application/vnd.github+json")
	request.Header.Add("Authorization", "Bearer "+gitToken)
	request.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	response, err = client.Do(request)

	if err != nil {
		errorSlice = append(errorSlice, err.Error())
		return
	}

	defer response.Body.Close()

	var body []byte

	body, err = io.ReadAll(response.Body)

	var PRInfos []*models.PRInfo
	var bodyAny []any

	err = json.Unmarshal(body, &PRInfos)
	json.Unmarshal(body, &bodyAny)

	if err != nil {

		err = errors.New(err.Error() + ", url: " + url)
		errorSlice = append(errorSlice, err.Error())
		return
	}

	if len(PRInfos) == 0 {
		return
	}

	messageID, err := slack.UpsertTodayMainMessage(
		slackBot,
		slackChannel,
		repoInfo.Repo+" 리뷰 요청",
		&repoInfo,
	)

	if err != nil {
		errorSlice = append(errorSlice, err.Error())
		return
	}

	slack.DeleteTheadMessage(
		slackBot,
		slackChannel,
		&repoInfo,
	)

	for _, PRInfo := range PRInfos {
		if IsTimeoutPR(PRInfo) {
			var reviewers []string
			reviewers, err = GetReviewers(PRInfo, gitToken)

			if err != nil {
				errorSlice = append(errorSlice, err.Error())
				continue
			}

			fmt.Println(reviewers)

			text := slack.CreateReviewersMessage(reviewers, PRInfo.User.Login)

			text = PRInfo.Url + "\n해당 PR 리뷰를 하거나 승인하세요.\n" + text

			slack.SendTheadMessage(
				slackBot,
				slackChannel,
				messageID,
				text,
			)
		}
	}
	return
}

func CallGitUrl(gitUrl, gitToken string) (body []byte, err error) {
	request, err := http.NewRequest(http.MethodGet, gitUrl, nil)

	if err != nil {
		return
	}

	request.Header.Add("Accept", "application/vnd.github+json")
	request.Header.Add("Authorization", "Bearer "+gitToken)
	request.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	client := http.Client{
		Timeout: utils.GetTimeout(constant.TIMEOUT_HTTP),
	}
	defer client.CloseIdleConnections()

	response, err := client.Do(request)

	if err != nil {
		return
	}

	defer response.Body.Close()

	body, err = io.ReadAll(response.Body)

	return
}

func IsTimeoutPR(PRInfo *models.PRInfo) (timout bool) {
	if PRInfo.State == "closed" {
		return
	}

	now := time.Now().Add(9 * time.Hour)

	created := time.Time(PRInfo.CreatedAt).Add(9 * time.Hour)

	PRDay := 3

	// 3일 지나면 slack 알림 전송
	for i := 0; i <= 3; i++ {
		nextDay := created.AddDate(0, 0, i)
		if nextDay.Weekday() == time.Saturday || nextDay.Weekday() == time.Sunday {
			PRDay++
		}
	}

	PRtime := created.AddDate(0, 0, PRDay)

	for {
		if PRtime.Weekday() == time.Saturday || PRtime.Weekday() == time.Sunday {
			PRtime = PRtime.AddDate(0, 0, 1)
			continue
		}
		break
	}

	return PRtime.Before(now)

}

func GetReviewers(PRInfo *models.PRInfo, gitToken string) (reviewerSlice []string, err error) {
	// configReviewers := viper.GetStringMapString(constant.REVIEWERS)

	reviewers := make(map[string]bool)

	for _, requestedReviewer := range PRInfo.RequestedReviewers {
		reviewers[requestedReviewer.Login] = false
	}

	body, err := CallGitUrl(PRInfo.Url+"/reviews", gitToken)

	if err != nil {
		return
	}

	var PRReviews []*models.PRReview

	err = json.Unmarshal(body, &PRReviews)

	if err != nil {
		return
	}

	for _, PRReview := range PRReviews {
		if reviewers[PRReview.User.Login] {
			continue
		}
		state := false
		if PRReview.State == constant.APPROVED {
			state = true
		}

		reviewers[PRReview.User.Login] = state
	}

	for reviewer, state := range reviewers {
		if state || reviewer == PRInfo.User.Login {
			continue
		}

		reviewerSlice = append(reviewerSlice, reviewer)
	}

	return

}

// repo 정보들을 config에서 가져온다.
func GetRepoInfos() (repoInfos []*models.RepoInfo, err error) {
	defer util_error.DeferWrap(err)

	repoInfoAny := viper.GetStringMap(constant.REPO_INFO)

	for key, repoInfoValue := range repoInfoAny {
		var repoInfoByte []byte
		var repoInfo *models.RepoInfo

		repoInfoByte, err = json.Marshal(repoInfoValue)

		if err != nil {
			return
		}

		err = json.Unmarshal(repoInfoByte, &repoInfo)

		if err != nil {
			return
		}
		repoInfo.Key = key

		repoInfos = append(repoInfos, repoInfo)

	}

	return
}
