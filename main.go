package main

import (
	"fmt"
	"github.com/atotto/clipboard"
	"log"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

func main() {

	getcontext := executecmd("kubectl", "config", "current-context")
	fmt.Println("Using context", getcontext)

	fmt.Println("Retrieving token")
	gettoken := executecmd("kubectl", "config", "view", "--minify")

	re := regexp.MustCompile(`password.*`)
	matches := re.FindStringSubmatch(gettoken)

	password_array := strings.Join(matches, ", ")
	password := strings.Trim(password_array, "password: ")

	fmt.Println("Sending token to your clipboard...")
	clipboard.WriteAll(password)

	fmt.Println("Opening your browser... \nPress Ctrl + C to finish!! :P")
	url := "http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/"
	openbrowser(url)

	executecmd("kubectl", "proxy")

}

func executecmd(cmd string, args ...string) string {

	out, err := exec.Command(cmd, args...).Output()

	if err != nil {
		fmt.Printf("%s", err)
		panic(err)
	}

	output := string(out[:])
	output = strings.TrimSpace(output) //context
	return output

}

func openbrowser(url string) {
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
