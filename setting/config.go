package setting

import "github.com/spf13/viper"

func SetConfig(paths ...string) {
	if len(paths) == 0 {
		paths = append(paths, "./config")
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	for _, path := range paths {
		viper.AddConfigPath(path) // 안에 path가 배열로 들어가 있어, 여러 경로를 등록해 놓을 수 있다!
	}

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

}
