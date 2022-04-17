package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/deitch/ocsso/auth"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	_ "github.com/zellyn/kooky/allbrowsers"
)

const (
	clientVersion       = "4.7.00136"
	defaultUserAgent    = "AnyConnect Linux_64 " + clientVersion
	cookieCheckInterval = 1
)

var (
	server, userAgent string
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "authenticate to a VPN and then launch openconnect",
	Long:  `Use SSO to authenticate to a VPN, and then connect via openconnect.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if server == "" {
			return errors.New("must pass argument --server <server>")
		}

		log.Debugf("starting init-auth request to %s", server)
		tunnelGroup, authMethod, configHash, ssoLoginURL, ssoCookieName, err := auth.InitAuth(server, clientVersion, userAgent, nil)
		if err != nil {
			log.Fatalf("init-auth failed: %v", err)
		}
		log.Debugf("completed init-auth request to %s", server)

		// open the webview and await the cookies
		log.Debugf("opening webview to %s", ssoLoginURL)
		token, err := auth.LoginAndWaitForCookie(ssoLoginURL, ssoCookieName)
		if err != nil {
			log.Fatalf("error authenticating: %v", err)
		}
		log.Debugf("openview successful, got token '%s", token)

		// send the auth finish request
		log.Debugf("starting finish-auth request to %s", server)
		sessionToken, serverCertHash, banner, err := auth.FinishAuth(server, clientVersion, tunnelGroup, authMethod, configHash, token, userAgent, nil)
		if err != nil {
			log.Fatalf("finish-auth failed: %v", err)
		}
		log.Debugf("finish-auth successful: session-token '%s", sessionToken)
		log.Debugf("finish-auth successful: server-cert-hash '%s", serverCertHash)
		if banner != "" {
			confirm := "Confirm"
			reader := bufio.NewReader(os.Stdin)
			fmt.Println(banner)
			fmt.Println("-------")
			fmt.Printf("type '%s' to accept, anything else to exit\n", confirm)

			fmt.Print("-> ")
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			if text != confirm {
				fmt.Println("Terms not accepted, exiting.")
				os.Exit(1)
			}

		}
		log.Debugf("running openconnect")
		if err := auth.RunOC(sessionToken, serverCertHash, server); err != nil {
			os.Exit(1)
		}

		return nil
	},
}

func addAuthHeaders(req *http.Request) {
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "identity")
	req.Header.Add("X-Transcend-Version", "1")
	req.Header.Add("X-Aggregate-Auth", "1")
	req.Header.Add("X-Support-HTTP-Auth", "true")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
}

func init() {
	connectCmd.PersistentFlags().StringVar(&server, "server", "", "server to authenticate against")
	connectCmd.PersistentFlags().StringVar(&userAgent, "user-agent", defaultUserAgent, "set the user-agent for requests")
}
