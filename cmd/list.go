// Copyright © 2019 xykong <xy.kong@gmail.com>

package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xykong/github-release/github"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `Information about published releases are available to everyone. 
Only users with push access will receive listings for draft releases.
`,
	Run: func(cmd *cobra.Command, args []string) {

		_ = viper.BindPFlag("id", cmd.PersistentFlags().Lookup("id"))

		owner := viper.GetString("user")
		repo := viper.GetString("repo")

		fmt.Printf("list called: %v, %s, %s\n", args, owner, repo)

		if viper.GetBool("assets") {
			github.ListAssets(owner, repo)
			return
		}

		_, err := github.ListReleases(owner, repo)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("list called")
		}
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
