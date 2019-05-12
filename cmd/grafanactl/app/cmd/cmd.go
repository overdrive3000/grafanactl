package cmd

import (
	"fmt"
	"io"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/overdrive3000/grafanactl/pkg/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	v       string
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "grafanactl",
	Version: version.VERSION,
	Short:   "A grafana CLI interface",
	Long: `A CLI which allows to perform operations in a Grafana
installation via command line by using Grafana's API.`,
}

func init() {
	cobra.OnInitialize(initConfig)
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

		// Search config in home directory with name ".grafanactl.yaml".
		viper.AddConfigPath(home)
		viper.SetConfigName(".grafanactl")
	}

	viper.SetConfigType("yaml") // set config file format to YAML
	viper.AutomaticEnv()        // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

// NewGrafanaCommand add subcommands to main CLI
func NewGrafanaCommand(out, err io.Writer) *cobra.Command {
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := SetUpLogs(err, v); err != nil {
			return err
		}
		rootCmd.SilenceUsage = true
		log.Infof("grafanactl %+v", version.VERSION)

		return nil
	}

	rootCmd.SilenceErrors = true
	rootCmd.AddCommand(folderCmd)

	rootCmd.PersistentFlags().StringVarP(&v, "verbosity", "v", log.WarnLevel.String(), "Log level (debug, warn, error)")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.grafanactl.yaml)")

	return rootCmd
}

// SetUpLogs set up logrus configuration
func SetUpLogs(out io.Writer, level string) error {
	log.SetOutput(out)
	lvl, err := log.ParseLevel(v)
	if err != nil {
		return err
	}
	log.SetLevel(lvl)
	return nil
}
