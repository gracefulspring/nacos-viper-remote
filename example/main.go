package main

import (
	"fmt"
	remote "github.com/gracefulspring/nacos-viper-remote"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config_viper := viper.New()
	runtime_viper := config_viper
	runtime_viper.SetConfigFile("./example_config.yaml")
	_ = runtime_viper.ReadInConfig()
	var option *remote.Option
	_ = runtime_viper.Sub("cloud.discovery.metadata").Unmarshal(&option)

	remote.SetOptions(option)

	//remote.SetOptions(&remote.Option{
	//	Url:         "localhost",
	//	Port:        80,
	//	NamespaceId: "public",
	//	GroupName:   "DEFAULT_GROUP",
	//	Config: 	 remote.Config{ DataId: "config_dev" },
	//	Auth:        nil,
	//})
	//localSetting := runtime_viper.AllSettings()
	remote_viper := viper.New()
	err := remote_viper.AddRemoteProvider("nacos", "10.1.120.30", "")
	remote_viper.SetConfigType("yaml")
	err = remote_viper.ReadRemoteConfig()

	if err == nil {
		config_viper = remote_viper
		fmt.Println("used remote viper")
		provider := remote.NewRemoteProvider("yaml")
		respChan := provider.WatchRemoteConfigOnChannel(config_viper)

		go func(rc <-chan bool) {
			for {
				<-rc
				fmt.Printf("remote async: %s", config_viper.GetString("application.name"))
			}
		}(respChan)

	}

	appName := config_viper.GetString("service.name")

	fmt.Println(appName)

	go func() {
		for {
			time.Sleep(time.Second * 30) // delay after each request
			appName = config_viper.GetString("application.name")
			fmt.Println("sync:" + appName)
		}
	}()

	onExit()
}

func onExit() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT, syscall.SIGKILL)

	for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			fmt.Println("Program Exit...", s)

		default:
			fmt.Println("other signal", s)
		}
	}
}
