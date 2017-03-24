package tester

import (
	"bytes"
	"fmt"
	"time"

	"encoding/json"

	"net/http"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/model"
)

type TransactionCreationTester struct {
	BasicTester
	NumberOfTransactions int
}

type TransactionReadTester struct {
	BasicTester
}

func NewTransactionCreationTester(tester BasicTester, numberOfTransactions int) *TransactionCreationTester {
	return &TransactionCreationTester{tester, numberOfTransactions}
}

func NewTransactionReadTester(tester BasicTester) *TransactionReadTester {
	return &TransactionReadTester{tester}
}

func (tct *TransactionCreationTester) RunFor(duration time.Duration) (int, int) {
	success := make(chan bool, 100)
	failure := make(chan bool, 100)

	timer := time.NewTimer(duration)
	go func() {
		for i := 0; i < tct.goroutines; i++ {
			go func(i int) {
				for k := 0; ; k++ {
					for j := 0; j < tct.NumberOfTransactions; j++ {
						id := fmt.Sprintf("%did%did%d", i, tct.cookie, k)
						tct.createTransaction(id, success, failure)
					}
				}
			}(i)
		}
	}()

	good, bad := listen(success, failure, timer.C)

	printResults("Creating Transactions", good, bad, duration)

	return good, bad
}

func (trt *TransactionReadTester) RunFor(duration time.Duration) (int, int) {
	success := make(chan bool, 100)
	failure := make(chan bool, 100)

	timer := time.NewTimer(duration)
	go func() {
		for i := 0; i < trt.goroutines; i++ {
			go func(i int) {
				for k := 0; ; k++ {
					for j := 0; j < 700; j++ {
						id := fmt.Sprintf("%did%did%d", i, trt.cookie, k)
						trt.readTransaction(id, success, failure)
					}
				}
			}(i)
		}
	}()

	good, bad := listen(success, failure, timer.C)

	printResults("Reading Transactions", good, bad, duration)

	return good, bad
}

func (tct *TransactionCreationTester) createTransaction(id string, success chan bool, failure chan bool) {
	bic := "BIC" + id
	iban := "IBAN" + id
	buffer := bytes.Buffer{}

	url := fmt.Sprintf("%s/accounts/%s/%s/transactions", tct.baseURL, bic, iban)

	noBalanceAccount := model.NoBalanceAccount{BIC: bic, IBAN: iban, Name: "Name" + id}
	transaction := model.Transaction{
		Recipient:           noBalanceAccount,
		Sender:              noBalanceAccount,
		Currency:            model.EUR,
		IntendedUse:         "Test",
		ValueInSmallestUnit: 100}

	json.NewEncoder(&buffer).Encode(transaction)

	resp, err := http.Post(url, "application/json", &buffer)

	if err != nil || resp.StatusCode > 201 {
		failure <- true

		return
	}

	resp.Body.Close()
	success <- true
}

func (trt *TransactionReadTester) readTransaction(id string, success chan bool, failure chan bool) {
	bic := "BIC" + id
	iban := "IBAN" + id

	url := fmt.Sprintf("%s/accounts/%s/%s/transactions", trt.baseURL, bic, iban)

	accountInfo := model.AccountInfo{}

	resp, err := http.Get(url)

	if err != nil || resp.StatusCode > 200 {
		failure <- true
		return
	}

	json.NewDecoder(resp.Body).Decode(&accountInfo)
	resp.Body.Close()

	if len(accountInfo.Transactions) == 0 {
		// log.Printf("No Transactions found for bic: %s iban: %s", bic, iban)
		failure <- true

		return
	}

	success <- true
}
