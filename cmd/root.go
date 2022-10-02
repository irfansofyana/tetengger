package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	username   string
	token      string
	repository string
	cfgFile    string

	ttCmd = &cobra.Command{
		Use:   "tetengger",
		Short: `tetengger - CLI to bookmark web content to your github repository`,
		Long: `Tetengger (adapted from Sundanese which means marker) 
		is a CLI tool that everyone can use to save (i.e bookmark) any content from the web 
		with free to their GitHub repository`,
		Version:       "0.2.0",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
)

func Execute() error {
	return ttCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	ttCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.tetengger.yaml)")
	ttCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "the GitHub username used to store the content.")
	ttCmd.PersistentFlags().StringVar(&token, "token", "", "the GitHub token to authenticate and authorize tetengger with the GitHub account.")
	ttCmd.PersistentFlags().StringVarP(&repository, "repository", "r", "", "The Github repository that will be used to store the content.")

	viper.BindPFlag("config", ttCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("username", ttCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("token", ttCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("repository", ttCmd.PersistentFlags().Lookup("repository"))

	ttCmd.AddCommand(saveCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.SetConfigFile(fmt.Sprintf("%s/.tetengger.yaml", home))
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Config file used for tetengger: ", viper.ConfigFileUsed())
	}
}
