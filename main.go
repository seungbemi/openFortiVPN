package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/seungbemi/gofred"
)

const (
	noSubtitle     = ""
	noArg          = ""
	noAutocomplete = ""
	configFolder   = "conf"
)

// Message adds simple message
func Message(response *gofred.Response, title, subtitle string, err bool) {
	msg := gofred.NewItem(title, subtitle, noAutocomplete)
	// if err {
	// 	msg = msg.AddIcon(iconError, defaultIconType)
	// } else {
	// 	msg = msg.AddIcon(iconDone, defaultIconType)
	// }
	response.AddItems(msg)
}

func init() {
	flag.Parse()
}

func main() {
	path := os.Getenv("PATH")
	if !strings.Contains(path, "/usr/local/bin") {
		os.Setenv("PATH", path+":/usr/local/bin")
	}
	configPath := os.Getenv("alfred_workflow_data") + "/" + configFolder

	response := gofred.NewResponse()
	err := os.MkdirAll(configPath, os.ModePerm)
	if err != nil {
		Message(response, "error", err.Error(), true)
		return
	}

	configs, err := ioutil.ReadDir(configPath)
	if err != nil {
		Message(response, "error", err.Error(), true)
		return
	}
	items := []gofred.Item{}

	if flag.Arg(0) != "create" {
		for _, config := range configs {
			name := config.Name()
			status := "Off"
			command := "Start"
			n, err := exec.Command("bash", "-c", fmt.Sprintf("ps aux | grep openfortivpn | grep -v grep | grep -v osascript | grep %s | wc -l", string(name))).CombinedOutput()
			if err == nil {
				number, err := strconv.Atoi(strings.TrimSpace(string(n)))
				if err == nil && number > 0 {
					status = "On"
					command = "Stop"
				}
			}
			items = append(items, gofred.NewItem(string(name), command+" "+string(name), noAutocomplete).
				AddIcon(status+".png", "").
				Executable(strings.ToLower(command)+" "+string(name)).
				AddOptionKeyAction("Remove config", "remove "+string(name), true).
				AddCommandKeyAction("Modify config", "modify "+string(name), true))
		}
		items = append(items, gofred.NewItem("Add new config", noSubtitle, "create ").AddIcon("plus.png", ""))
	} else {
		response.AddVariable("filename", flag.Arg(1))
		items = append(items, gofred.NewItem("Add new config", fmt.Sprintf("write name ... \"%s\"", flag.Arg(1)), noAutocomplete).AddIcon("plus.png", "").
			Executable("newConfig"))
	}

	response.AddItems(items...)
	fmt.Println(response)
}
