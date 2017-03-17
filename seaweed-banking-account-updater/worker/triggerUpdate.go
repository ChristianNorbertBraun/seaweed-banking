package worker

import (
	"bytes"
	"log"
	"time"

	"encoding/json"

	"net/http"

	"fmt"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/database"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/handler"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/model"
)

var UpdateTicker *time.Ticker

// SetUpUpdateWorker starts a worker which checks the update table for new entries
// and updates the matching accountInfos
func SetUpUpdateWorker(duration time.Duration) {
	UpdateTicker = time.NewTicker(duration)
	go func() {
		for t := range UpdateTicker.C {
			log.Println("-------------------------------------------------------------")
			log.Println("Start Update at: ", t)

			runUpdate()
		}
	}()
}

// StopTicker stops the given ticker
func StopTicker(ticker *time.Ticker) {
	if ticker != nil {
		ticker.Stop()
	}
}

func runUpdate() {
	updates, err := database.FindAllUpdates()
	if err != nil {
		log.Println("Unable to get updates", err)

		return
	}

	distributeUpdates(updates)
}

func distributeUpdates(updates []*model.Update) {
	numberOfSubscribers := len(subscribers)
	if numberOfSubscribers == 0 {
		log.Println("No subscribers have to do all the work allone. Pff.")
	} else {
		log.Printf("Distributing updates on %d subscribers", numberOfSubscribers)
	}
	// the master itself is also a worker therefore +1
	numberOfUpdatesPerWorker := len(updates) / (numberOfSubscribers + 1)

	if len(updates) < numberOfSubscribers+1 {
		numberOfUpdatesPerWorker = 1
	}

	start := 0
	for url := range subscribers {
		if start >= len(updates) {
			return
		}
		if err := sendUpdates(url, updates[start:start+numberOfUpdatesPerWorker]); err != nil {
			Unregister(url)
		}
		start = start + numberOfUpdatesPerWorker
	}

	handler.UpdateAccountInfo(updates[start:])
}

func sendUpdates(url string, updates []*model.Update) error {
	buffer := bytes.Buffer{}

	if err := json.NewEncoder(&buffer).Encode(updates); err != nil {
		return err
	}

	resp, err := http.Post(url+"/do/update", "application/json", &buffer)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unable to send update to %s", url)
	}

	return nil
}
