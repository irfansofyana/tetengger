package cmd

import (
	"fmt"
	"os"

	"github.com/irfansofyana/tetengger/pkg/content"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	commit string
	folder string
	branch string
	tags   string

	saveCmd = &cobra.Command{
		Use:   "save",
		Short: "Save the web content into your GitHub repository in markdown format",
		Run:   save,
	}
)

func save(ccmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "No URL is specified. Please specify a valid URL name that you want to save!")
		return
	}
	if len(args) == 1 {
		fmt.Fprintln(os.Stderr, "No name of the content is specified. Please specify a name for the content that you want to save!")
		return
	}

	content.Save(args[0], args[1])

	saveConfig()
}

func init() {
	saveCmd.PersistentFlags().StringVarP(&commit, "commit", "m", "save a new content", "the commit message that will be add when save the content to the Github repository.")
	saveCmd.PersistentFlags().StringVarP(&folder, "folder", "f", "bookmark", "the parent folder that will be used to store the content in the GitHub repository.")
	saveCmd.PersistentFlags().StringVarP(&branch, "branch", "b", "main", "the GitHub branch where the content will be saved.")
	saveCmd.PersistentFlags().StringVarP(&tags, "tags", "t", "", "tags of the content (comma separated values)")

	viper.BindPFlag("commit", saveCmd.PersistentFlags().Lookup("commit"))
	viper.BindPFlag("folder", saveCmd.PersistentFlags().Lookup("folder"))
	viper.BindPFlag("branch", saveCmd.PersistentFlags().Lookup("branch"))
	viper.BindPFlag("tags", saveCmd.PersistentFlags().Lookup("tags"))
}

func saveConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	viper.WriteConfigAs(home + "/.tetengger.yaml")
}
