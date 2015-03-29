/**
 * @author Florin Patan <florinpatan@gmail.com>
 * MIT LICENSE, see LICENSE file for full license
 */

// Command externalipchecker checks the external IP of the computer
// and reports back any changes to it
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/0xAX/notificator"
)

var (
	client       = &http.Client{}
	req          *http.Request
	notify       *notificator.Notificator
	startIP      string
	interval     uint = 60
	checkService      = "icanhazip.com"
)

func init() {
	flag.UintVar(&interval, "interval", interval, "Interval on which the external ip check is performed")
	flag.StringVar(&checkService, "service", checkService, "Service used to check the IP address")

	notify = notificator.New(notificator.Options{
		AppName: "External IP Checker",
	})
}

func fetchIP() string {
	resp, err := client.Do(req)
	if err != nil {
		displayMessage(fmt.Sprintf("error encountered: %q\n", err))
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		displayMessage(fmt.Sprintf("error encountered: %q\n", err))
		return ""
	}

	return strings.Trim(string(body), "\n")
}

func displayMessage(message string) {
	notify.Push("External IP Checker", message, "")
}

func main() {
	flag.Parse()

	var err error
	req, err = http.NewRequest("GET", fmt.Sprintf("http://%s", checkService), nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", "External IP Checker / 1.0 (github.com/dlsniper/externalipchecker)")

	startIP = fetchIP()
	if startIP == "" {
		panic("failed to get a proper IP address")
	}
	displayMessage(fmt.Sprintf("Your External IP address is: %s", startIP))

	checkInterval := time.Duration(interval) * time.Second

	for {
		select {
		case <-time.After(checkInterval):
			ip := fetchIP()
			if startIP == ip {
				continue
			}

			displayMessage(fmt.Sprintf("Your External IP address CHANGED\nOLD: %s\nNEW: %s", startIP, ip))
			startIP = ip
		}
	}
}
