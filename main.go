package main

import (
	"time"
	constant "tnals5152/git-pr-slack/const"
	"tnals5152/git-pr-slack/git"
	"tnals5152/git-pr-slack/setting"

	"github.com/spf13/viper"
)

func init() {
	setting.SetConfig()
}
func main() {

	ticker := time.NewTicker(5 * time.Minute)

	for range ticker.C {

		weekday := time.Now().Weekday()

		// 오늘이 토요일이거나 일요일이면 작동 안 하게 하기
		if weekday == time.Saturday || weekday == time.Sunday {
			continue
		}

		repoInfos, _ := git.GetRepoInfos()
		gitToken := viper.GetString(constant.GIT_TOKEN)
		git.GetPRList(gitToken, repoInfos)
	}

}
