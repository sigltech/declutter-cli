package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "declutter",
		Short: "A CLI for managing you system",
		Long: `This CLI is designed to help you manage your system.
		It can be used to clean up your system of files and folders.`,
		Version: "0.0.1",
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	err := viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	if err != nil {
		fmt.Println(err)
		return
	}
	viper.SetDefault("author", "Richard Sigl")
	viper.SetDefault("license", "apache")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
