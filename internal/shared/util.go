package shared

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

func GuessIsProfileSav(name string) bool {
	parts := strings.Split(name, "/")
	name = parts[len(parts)-1]
	if name == "profile.sav" {
		return true
	}
	return false
}

func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}
