package worker

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/config"
)

var HealthTicker *time.Ticker

func SetUpSlavePing(duration time.Duration) {
	ipAddress, err := getIPAddress()
	if err != nil {
		log.Println("Failed to get ip address", err)
		return
	}
	runHealthPing(ipAddress)

	HealthTicker = time.NewTicker(duration)
	go func() {
		for t := range HealthTicker.C {
			log.Println("Pinging master on: ", t)

			runHealthPing(ipAddress)
		}
	}()
}

func runHealthPing(ipAddress string) {
	url := fmt.Sprintf("%s:%s/register",
		config.Configuration.Master.Host,
		config.Configuration.Master.Port)

	body := bytes.Buffer{}
	body.WriteString(fmt.Sprintf("http://%s:%s",
		ipAddress,
		config.Configuration.Server.Port))

	resp, err := http.Post(url, "application/json", &body)

	if err != nil {
		log.Println("Unable to send heartbeat to master:", url, err)

		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("Bad Statuscode on health ping")

		return
	}
}

func getIPAddress() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", errors.New("Can't get my ip address!")
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", errors.New("Can't find my ip address")
}
