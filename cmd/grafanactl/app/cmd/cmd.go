package cmd

import (
	"fmt"
	"io"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	gapi "github.com/overdrive3000/go-grafana-api"
	"github.com/overdrive3000/grafanactl/pkg/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type globalOpts struct {
	verbose string
	cfgFile string
	url     string
	key     string
	output  string
}

var config string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "grafanactl",
	Version: version.VERSION,
	Short:   "A grafana CLI interface",
	Long: `A CLI which allows to perform operations in a Grafana
installation via command line by using Grafana's API.`,
}

// SetUpClient set up a new grafana client
func SetUpClient() (*gapi.Client, error) {
	log.Debugf("Setting up grafana client with url %s and key %s", viper.GetString("url"), viper.GetString("apiKey"))
	return gapi.New(
		viper.GetString("apiKey"),
		viper.GetString("url"),
	)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if config != "" {
		// Use config file from the flag.
		viper.SetConfigFile(config)
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
func NewGrafanaCommand() *cobra.Command {
	opts := &globalOpts{}

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := SetUpLogs(os.Stdout, &opts.verbose); err != nil {
			return err
		}
		rootCmd.SilenceUsage = true
		log.Infof("grafanactl %+v", version.VERSION)

		return nil
	}

	rootCmd.SilenceErrors = true

	rootCmd.PersistentFlags().StringVarP(&opts.verbose, "verbosity", "v", log.WarnLevel.String(), "Log level (debug, warn, error)")
	rootCmd.PersistentFlags().StringVarP(&opts.cfgFile, "config", "c", "", "config file (default is $HOME/.grafanactl.yaml)")
	// Grafana configuration flags
	rootCmd.PersistentFlags().StringVar(&opts.url, "url", "", "Grafana URL (https://localhost:3000)")
	rootCmd.PersistentFlags().StringVar(&opts.key, "key", "", "Grafana API Key")
	rootCmd.PersistentFlags().StringVarP(&opts.output, "output", "o", "table", "Output format (table, json)")

	// Bind flags to config file
	if err := viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url")); err != nil {
		log.Error(err)
	}
	viper.BindPFlag("apiKey", rootCmd.PersistentFlags().Lookup("key"))

	config = opts.cfgFile
	cobra.OnInitialize(initConfig)

	// Add main commands
	rootCmd.AddCommand(folderCmd())
	rootCmd.AddCommand(dashboardCmd())

	return rootCmd
}

// SetUpLogs set up logrus configuration
func SetUpLogs(out io.Writer, verbose *string) error {
	log.SetOutput(out)
	lvl, err := log.ParseLevel(*verbose)
	if err != nil {
		return err
	}
	log.SetLevel(lvl)
	return nil
}
