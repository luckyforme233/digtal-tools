package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type conf struct {
	Token      string `toml:"token"`
	PrvKeyPath string `toml:"prvKeyPath"`
	PubKeyPath string `toml:"pubKeyPath"`
	CLEmail    string `toml:"CLEmail"`
	CLApiKey   string `toml:"CLApiKey"`
	CLDomain   string `toml:"CLDomain"`
}

var C conf

func InitConfig() {
	fmt.Println("init config")
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("toml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		log.Fatalln("文件加载失败")
		fmt.Println("err", err)
		return
	}
	err = viper.Unmarshal(&C)
	if err != nil {
		log.Fatalln("配置文件反序列化失败")
		fmt.Println("err", err)
		return
	}
}
