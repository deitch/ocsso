package auth

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/webview/webview"
	"golang.org/x/net/context"
)

const (
	js                = `findCookie(document.cookie)`
	copyPasteShortcut = `
window.addEventListener("keypress", (event) => {
  if (event.metaKey && event.key === 'c') {
    document.execCommand("copy")
    event.preventDefault();
  }
  if (event.metaKey && event.key === 'v') {
    document.execCommand("paste")
    event.preventDefault();
  }
})
`
)

func LoginAndWaitForCookie(server, cookieName string) (string, error) {
	// get my command-line
	var (
		out     bytes.Buffer
		cmdName = os.Args[0]
	)
	wvCmd := exec.Command(cmdName, "webview", "--webview-url", server, "--cookie", cookieName)
	wvCmd.Stdout = &out
	wvCmd.Stderr = os.Stderr
	if err := wvCmd.Run(); err != nil {
		return "", err
	}
	result := out.String()
	parts := strings.SplitN(result, "=", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("result was '%s', which does not split into 2 on '='", result)
	}
	if parts[0] != cookieName {
		return "", fmt.Errorf("cookie name returned was '%s' and not expected '%s", parts[0], cookieName)
	}

	// we have a result, so now we can just send
	return parts[1], nil
}

func LoginAndWaitForCookieProcess(server, cookieName string) {
	// create the channel we need
	var token string
	result := make(chan string, 1)
	ctx, cancel := context.WithCancel(context.Background())
	debug := true
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle(fmt.Sprintf("SSO %s", server))
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate(server)
	w.Eval(copyPasteShortcut)
	w.Bind("findCookie", cookieFinder(cookieName, result, cancel))
	w.Dispatch(func() {
		go func() {
			for {
				select {
				case <-time.After(time.Second * 1):
					w.Eval(js)
				case <-ctx.Done():
					return
				}
			}
		}()
	})
	go func() {
		token = <-result
		fmt.Printf("%s=%s", cookieName, token)
		w.Terminate()
	}()
	w.Run()
}

func cookieFinder(cookieName string, c chan<- string, cancel context.CancelFunc) func(string) {
	return func(cookiesString string) {
		cookies := strings.Split(cookiesString, ";")
		for _, cookie := range cookies {
			parts := strings.SplitN(cookie, "=", 2)
			if len(parts) != 2 {
				continue
			}
			name, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
			if name == cookieName {
				c <- value
				cancel()
				return
			}
		}
	}
}
