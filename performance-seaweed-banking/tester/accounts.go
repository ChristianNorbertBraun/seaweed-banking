package tester

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"time"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/model"
)

type AccountTester struct {
	BasicTester
}

func NewAccountTester(tester BasicTester) *AccountTester {
	return &AccountTester{tester}
}

func (at *AccountTester) RunFor(duration time.Duration) (int, int) {
	success := make(chan bool, 100)
	failure := make(chan bool, 100)
	timer := time.NewTimer(duration)

	go func() {
		for i := 0; i < at.goroutines; i++ {
			go func(i int) {
				for k := 0; ; k++ {
					id := fmt.Sprintf("%did%did%d", i, at.cookie, k)
					at.createAccount(id, success, failure)
				}
			}(i)
		}
	}()

	good, bad := listen(success, failure, timer.C)

	printResults("Creating accounts", good, bad, duration)

	return good, bad
}

func (at *AccountTester) createAccount(id string, success chan bool, failure chan bool) {

	noBalanceAccount := model.NoBalanceAccount{BIC: "BIC" + id, IBAN: "IBAN" + id, Name: "Name" + id}
	account := model.Account{NoBalanceAccount: noBalanceAccount, Balance: 100}

	buffer := bytes.Buffer{}

	json.NewEncoder(&buffer).Encode(account)

	resp, err := http.Post(fmt.Sprintf("%s/accounts", at.baseURL), "application/json", &buffer)

	if err != nil || resp.StatusCode > 201 {
		failure <- true
		return
	}

	resp.Body.Close()
	success <- true
}
