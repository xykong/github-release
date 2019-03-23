package utils

import (
	"fmt"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("github-release")        // name of config file (without extension)
	viper.AddConfigPath("/etc/")                 // path to look for the config file in
	viper.AddConfigPath("$HOME/.github-release") // call multiple times to add many search paths
	viper.AddConfigPath(".")                     // optionally look for config in the working directory
	err := viper.ReadInConfig()                  // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
