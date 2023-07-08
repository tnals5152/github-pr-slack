package models

type MessageResponse struct {
	Ts string `json:"ts"`
}

type Message struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
	Ts      string `json:"ts"`
	TheadTs string `json:"thread_ts"` // parent message ID
}

type ResponseOkMessage struct {
	Messages []*Message `json:"messages"`
}
