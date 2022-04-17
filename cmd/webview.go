package cmd

import (
	"errors"

	"github.com/deitch/ocsso/auth"
	"github.com/spf13/cobra"
)

var cookieName, webviewURL string

var webviewCmd = &cobra.Command{
	Use:   "webview",
	Short: "launch webview to get token",
	Long: `Launch webview to get token. This is in a separate command because of limitations on Webview,
	wherein if you close the webview, the whole process closes. Will print the cookie on stdout.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if webviewURL == "" {
			return errors.New("must pass argument --webview-url <url>")
		}
		if cookieName == "" {
			return errors.New("must pass argument --cookie <name-of-cookie>")
		}

		auth.LoginAndWaitForCookieProcess(webviewURL, cookieName)
		return nil
	},
}

func init() {
	webviewCmd.Flags().StringVar(&webviewURL, "webview-url", "", "url to use for webview")
	webviewCmd.Flags().StringVar(&cookieName, "cookie", "", "cookie to use")
}
