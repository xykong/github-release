// Copyright © 2019 xykong <xy.kong@gmail.com>

package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/xykong/github-release/github"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
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

		fmt.Printf("delete called: %v, %s, %s\n", args, owner, repo)

		github.DeleteRelease(owner, repo)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")
	deleteCmd.PersistentFlags().StringP("id", "i", "", "The id of the release")
	_ = viper.BindPFlag("id", deleteCmd.PersistentFlags().Lookup("id"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
