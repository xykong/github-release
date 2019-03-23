// Copyright © 2019 xykong <xy.kong@gmail.com>

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xykong/github-release/github"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		_ = viper.BindPFlag("id", cmd.PersistentFlags().Lookup("id"))

		owner := viper.GetString("user")
		repo := viper.GetString("repo")

		fmt.Printf("list called: %v, %s, %s\n", args, owner, repo)

		if viper.GetBool("assets") {
			github.ListAssets(owner, repo)
			return
		}

		_, _ = github.GetReleases(owner, repo)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")
	listCmd.PersistentFlags().BoolP("assets", "a", false, "Users with push access to the repository can edit a release.")
	_ = viper.BindPFlag("assets", listCmd.PersistentFlags().Lookup("assets"))

	listCmd.PersistentFlags().StringP("id", "i", "", "The id of the release")
	_ = viper.BindPFlag("id", listCmd.PersistentFlags().Lookup("id"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
