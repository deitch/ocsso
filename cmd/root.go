package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	_ "github.com/zellyn/kooky/allbrowsers"
)

var (
	debug bool
)

var rootCmd = &cobra.Command{
	Use:   "ocsso",
	Short: "authenticate to a VPN and then launch openconnect",
	Long:  `Use SSO to authenticate to a VPN, and then connect via openconnect.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if debug {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debugging output")
	rootCmd.AddCommand(webviewCmd)
	rootCmd.AddCommand(connectCmd)
}

// Execute primary function for cobra
func Execute() {
	rootCmd.Execute()
}
