package auth

import (
	"os"
	"os/exec"
)

func RunOC(token, serverCertHash, server string) error {
	ocCmd := exec.Command("openconnect", "--cookie", token, "--servercert", serverCertHash, server)
	ocCmd.Stdout = os.Stdout
	ocCmd.Stderr = os.Stderr
	return ocCmd.Run()
}
