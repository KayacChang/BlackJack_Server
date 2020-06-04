package conf

import (
	"fmt"

	"github.com/spf13/viper"
)

var ULG168APIHost, MaintainAPI, MaintainToken string

func init() {
	viper.SetConfigName("thirdparty") // name of config file (without extension)
	viper.SetConfigType("yml")
	viper.AddConfigPath("./conf/conf.d") // optionally look for config in the working directory

	if err := viper.ReadInConfig(); err != nil { // Find and read the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	ULG168APIHost = viper.GetString("host")
	MaintainAPI = viper.GetString("maintain")
	MaintainToken = viper.GetString("maintain_token")
}
