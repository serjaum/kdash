package main

import (
    "fmt"
    "regexp"
    "strings"
    "os/exec"
    "runtime"
    "log"
    "github.com/atotto/clipboard"
)

func main() {

    cmdName := "kubectl"
    cmdArgs := []string{"config", "current-context"}
    out, err := exec.Command(cmdName, cmdArgs...).Output()

    if err != nil {
        fmt.Printf("%s", err)
    }

    output := string(out[:])
    output = strings.TrimSpace(output) //context
    fmt.Println("Using context", output)

    fmt.Println("Retrieving token...")
    cmdName = "kubectl"
    cmdArgs = []string{"config", "view", "--minify"}
    out, err = exec.Command(cmdName, cmdArgs...).Output()

    if err != nil {
        fmt.Printf("%s", err)
    }

    output = string(out[:])
    output = strings.TrimSpace(output)

    re := regexp.MustCompile(`password.*`)
    matches := re.FindStringSubmatch(output)

    password_array := strings.Join(matches, ", ")
    password := strings.Trim(password_array, "password: ") //password

    fmt.Println("Sending token to your clipboard...")
    clipboard.WriteAll(password) //clipboard

    fmt.Println("Opening your browser... \nPress Ctrl + C to finish!! :P")
    url := "http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/"
    openbrowser(url)

    cmdName = "kubectl"
    cmdArgs = []string{"proxy"}
    out, err = exec.Command(cmdName, cmdArgs...).Output()

    if err != nil {
        fmt.Printf("%s", err)
    }

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
