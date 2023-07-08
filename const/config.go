package constant

import "time"

const (
	REPO_INFO        = "repo-info"
	REVIEWERS        = "reviewers"
	GIT_TOKEN        = "git-token"
	SLACK_BOT_TOKEN  = "slack.bot-token"
	SLACK_CHANNEL    = "slack.channel"
	SLACK_USER_TOKEN = "slack.user-token"

	TIMEOUT_HTTP = "timeout.http"
)

// 기본 타임아웃 세팅
var ContextTimeoutMap map[string]time.Duration = map[string]time.Duration{
	TIMEOUT_HTTP: 10,
}
