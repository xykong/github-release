// Copyright © 2019 xykong <xy.kong@gmail.com>

package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xykong/github-release/utils"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "github-release",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.github-release.yaml)")

	rootCmd.PersistentFlags().StringP("user", "u", "", "The authenticated user owned the repository")
	_ = viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))

	rootCmd.PersistentFlags().StringP("repo", "r", "", "The name of the repository")
	_ = viper.BindPFlag("repo", rootCmd.PersistentFlags().Lookup("repo"))

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Verbose message for debug")
	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose message for debug")
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.PersistentFlags().BoolP("essential", "s", false, "Verbose message for debug")
	_ = viper.BindPFlag("essential", rootCmd.PersistentFlags().Lookup("essential"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".github-release" (without extension).
		viper.AddConfigPath(".") // optionally look for config in the working directory
		viper.AddConfigPath(home)
		viper.AddConfigPath("/usr/local/etc/") // path to look for the config file in
		viper.AddConfigPath("/etc/")           // path to look for the config file in
		viper.SetConfigName(".github-release")
	}

	viper.SetDefault("github", "https://api.github.com")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		utils.Verbose("Using config file:", viper.ConfigFileUsed())
	}

	level := logrus.InfoLevel
	if viper.GetBool("verbose") {
		level = logrus.DebugLevel
	}
	logrus.SetLevel(level)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableLevelTruncation: false,
		DisableTimestamp:       true,
		FullTimestamp:          true,
	})
}
