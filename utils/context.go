package utils

import (
	"context"
	"time"
	constant "tnals5152/git-pr-slack/const"

	"github.com/spf13/viper"
)

// key = config에 저장된 키
func GetContext(key string) (context.Context, context.CancelFunc) {
	// var timeout time.Duration
	timeout := GetTimeout(key)

	return context.WithTimeout(context.Background(), timeout*time.Second)
}

func GetTimeout(key string) time.Duration {
	if viper.InConfig(key) {
		return viper.GetDuration(key) * time.Second
	}
	value, ok := constant.ContextTimeoutMap[key]

	if ok {
		return value
	}
	return 10

}
