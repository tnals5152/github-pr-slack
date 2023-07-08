package constant

const (
	SLASH = "/"

	PR_LIST_URL        = "https://api.github.com/repos/{owner}/{repo}/pulls"
	PR_INFO_URL        = "https://api.github.com/repos/{owner}/{repo}/pulls/{pull_number}?state=all"
	SLACK_DELETE_URL   = "https://slack.com/api/chat.delete"
	SLACK_REPLY_URL    = "https://slack.com/api/conversations.replies"
	SLACK_POST_MESSAGE = "https://slack.com/api/chat.postMessage"

	OWNER_PARAM       = "{owner}"
	REPO_PARAM        = "{repo}"
	PULL_NUMBER_PARAM = "{pull_number}"

	APPROVED = "APPROVED"
	// COMMENTED = "COMMENTED"
)
