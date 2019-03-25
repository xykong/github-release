// Copyright © 2019 xykong <xy.kong@gmail.com>

package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xykong/github-release/github"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "View the specified published full release for the repository.",
	Long: `This returns an upload_url key corresponding to the endpoint for uploading release assets. 
This key is a hypermedia resource.

If the id is not specified, this will return the latest release.
The latest release is the most recent non-prerelease, non-draft release,
sorted by the created_at attribute. The created_at attribute is the date
of the commit used for the release, and not the date when the release was
drafted or published.

Get a published release with the specified tag.
`,

	Run: func(cmd *cobra.Command, args []string) {

		owner := viper.GetString("user")
		repo := viper.GetString("repo")
		releaseId := viper.GetString("id")
		tag := viper.GetString("tag")

		fmt.Printf("show called: %v, %s, %s, %s\n", args, owner, repo, releaseId)

		var err error
		if tag != "" {
			_, err = github.GetReleaseByTag(owner, repo, tag)
		} else {
			_, err = github.GetRelease(owner, repo, releaseId)
		}

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("show called")
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")
	viper.SetDefault("id", "latest")

	showCmd.PersistentFlags().StringP("id", "i", "", "The id of the release")
	_ = viper.BindPFlag("id", showCmd.PersistentFlags().Lookup("id"))

	showCmd.PersistentFlags().StringP("tag", "t", "", "The tag of the release")
	_ = viper.BindPFlag("tag", showCmd.PersistentFlags().Lookup("tag"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
