// Copyright © 2019 xykong <xy.kong@gmail.com>

package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xykong/github-release/github"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Users with push access to the repository can create a release.",
	Long: `This endpoint triggers notifications. Creating content too quickly 
using this endpoint may result in abuse rate limiting. 
See "Abuse rate limits" and "Dealing with abuse rate limits" 
for details.
`,
	Run: func(cmd *cobra.Command, args []string) {

		_ = viper.BindPFlag("id", cmd.PersistentFlags().Lookup("id"))

		owner := viper.GetString("user")
		repo := viper.GetString("repo")

		fmt.Printf("create called: %v, %s, %s\n", args, owner, repo)

		err := github.CreateRelease(owner, repo)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("create called")
		}
	},
	Example: `github-release create --tag_name v0.0.1\
                      --name "The name of the release."\
                      --body "Text describing the contents of the tag."`,
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")
	createCmd.PersistentFlags().StringP("token", "t", "", "The tag of the release")
	_ = viper.BindPFlag("token", createCmd.PersistentFlags().Lookup("token"))

	createCmd.PersistentFlags().StringP("tag_name", "", "", "The tag of the release")
	_ = viper.BindPFlag("tag_name", createCmd.PersistentFlags().Lookup("tag_name"))

	createCmd.PersistentFlags().StringP("target_commitish", "", "", "The tag of the release")
	_ = viper.BindPFlag("target_commitish", createCmd.PersistentFlags().Lookup("target_commitish"))

	createCmd.PersistentFlags().StringP("name", "", "", "The tag of the release")
	_ = viper.BindPFlag("name", createCmd.PersistentFlags().Lookup("name"))

	createCmd.PersistentFlags().StringP("body", "", "", "The tag of the release")
	_ = viper.BindPFlag("body", createCmd.PersistentFlags().Lookup("body"))

	createCmd.PersistentFlags().BoolP("draft", "", false, "The tag of the release")
	_ = viper.BindPFlag("draft", createCmd.PersistentFlags().Lookup("draft"))

	createCmd.PersistentFlags().BoolP("prerelease", "", false, "The tag of the release")
	_ = viper.BindPFlag("prerelease", createCmd.PersistentFlags().Lookup("prerelease"))

	createCmd.PersistentFlags().BoolP("edit", "e", false, "Users with push access to the repository can edit a release.")
	_ = viper.BindPFlag("edit", createCmd.PersistentFlags().Lookup("edit"))

	createCmd.PersistentFlags().StringP("id", "i", "", "The id of the release")
	_ = viper.BindPFlag("id", createCmd.PersistentFlags().Lookup("id"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
