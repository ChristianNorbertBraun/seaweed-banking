package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/config"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/handler"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/model"
	"github.com/gorilla/mux"
)

type fakeAccount struct {
	account      model.Account
	transactions []model.Transaction
}

var testData []fakeAccount
var r *mux.Router
var writer *httptest.ResponseRecorder
var testConfigPath = flag.String("testConfig", "./data/conf/testconfig.json", "Path to json formated testconfig")

func TestMain(m *testing.M) {
	setUp()
	initTestData()
	test := m.Run()
	os.Exit(test)
}

func setUp() {
	flag.Parse()

	err := config.ParseTestConfig(*testConfigPath)

	if err != nil {
		log.Fatalf("Unable to parse testconfig from: %s because: %s",
			*configPath,
			err)
	}

	r = mux.NewRouter()
	r.HandleFunc("/accounts", handler.GetAllAccounts).Methods("GET")
	r.HandleFunc("/accounts", handler.CreateAccount).Methods("POST")
	r.HandleFunc("/accounts/{bic}/{iban}/transactions", handler.GetTransaction).Methods("GET")
	r.HandleFunc("/accounts/{bic}/{iban}/transactions", handler.CreateTransaction).Methods("POST")
	writer = httptest.NewRecorder()

}

func initTestData() {

	for i := 0; i < config.TestConfiguration.NoOfFakeAccounts; i++ {

		var newFakeAccount fakeAccount

		newFakeAccount.account.BIC = RandBIC()
		newFakeAccount.account.IBAN = RandIBAN()
		newFakeAccount.account.Balance = RandNumberWithRange(200, 10000)

		testData = append(testData, newFakeAccount)
	}
}

/*
*	TEST CASES
 */

func TestAccountsCreate(t *testing.T) {

	for _, data := range testData {
		err := CreateAccount(data.account)

		if err != nil {
			t.Error(err)
		}
	}

	WaitForUpdater()

	for _, data := range testData {
		err := VerifyAccount(data.account)

		if err != nil {
			t.Error(err)
		}
	}
}

func TestTransactionsCreate(t *testing.T) {

	if len(testData) > 1 {

		for i := range testData {
			if len(testData) > i+1 {
				var newTransaction model.Transaction
				newTransaction.BIC = testData[i+1].account.BIC
				newTransaction.IBAN = testData[i+1].account.IBAN
				newTransaction.BookingDate = time.Now()
				newTransaction.Currency = model.EUR
				newTransaction.IntendedUse = "testingTransactions"
				newTransaction.ValueInSmallestUnit = RandNumberWithRange(0, int(testData[i].account.Balance))

				err := CreateTransaction(testData[i].account, newTransaction)

				if err != nil {
					t.Error(err)
				}

				testData[i+1].transactions = append(testData[i+1].transactions, newTransaction)
			}
		}
	} else {
		t.Errorf("At least 2 accounts are necessary for transactions, current number of accounts: %v",
			cap(testData))
	}

	WaitForUpdater()

	for i := range testData {
		if i > 0 {
			err := VerifyTransactions(testData[i].account)

			if err != nil {
				t.Error(err)
			}
		}
	}
}

//func TestCreateTransactionsParallel(t *testing.T) {
//TODO
//}
//
// func TestReadAllAccountsParallel(t *testing.T) {
// 	//TODO
// 	// Readout all Accounts parallel
// }
//
// func TestReadAllTransactions(t *testing.T) {
// 	// TODO
// 	// Readout All Transactions for every account parallel with benchmark
// }

/*
*	HELPERS
 */
func CreateAccount(account model.Account) (err error) {

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(&account)

	request, _ := http.NewRequest("POST", "/accounts", b)
	r.ServeHTTP(writer, request)

	if writer.Code != 200 && writer.Code != 201 {
		err = fmt.Errorf("CreateAccount: %v \nResponse Code: %v",
			request.URL.String(),
			writer.Code)
		return
	}

	err = nil
	return
}

func CreateTransaction(account model.Account, transaction model.Transaction) (err error) {

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(&transaction)

	request, _ := http.NewRequest("POST", "/accounts/"+account.BIC+"/"+account.IBAN+"/transactions", b)
	r.ServeHTTP(writer, request)

	if writer.Code != 200 && writer.Code != 201 {
		err = fmt.Errorf("CreateTransaction: %v \nResponse Code: %v",
			request.URL.String(),
			writer.Code)
		return
	}

	err = nil
	return
}

func GetAllAccounts() (readData []byte, err error) {

	request, _ := http.NewRequest("GET", "/accounts", nil)

	r.ServeHTTP(writer, request)

	if writer.Code != 200 && writer.Code != 201 {
		err = fmt.Errorf("GetAllAccounts: %v \nResponse Code: %v",
			request.URL.String(),
			writer.Code)
		readData = nil
		return
	}

	readData = writer.Body.Bytes()
	err = nil
	return
}

func VerifyAccount(acc model.Account) (err error) {

	request, _ := http.NewRequest("GET", "/accounts/"+acc.BIC+"/"+acc.IBAN, nil)
	r.ServeHTTP(writer, request)

	if writer.Code != 200 && writer.Code != 201 {
		err = fmt.Errorf("VerifyAccount: %v \nResponse Code: %v",
			request.URL.String(),
			writer.Code)
		return
	}

	err = nil
	return
}

func VerifyTransactions(account model.Account) (err error) {

	request, _ := http.NewRequest("GET", "/accounts/"+account.BIC+"/"+account.IBAN+"/transactions", nil)

	r.ServeHTTP(writer, request)

	if writer.Code != 200 && writer.Code != 201 {
		err = fmt.Errorf("VerifyTransactions: %v \n Response Code: %v",
			request.URL.String(),
			writer.Code)
		return
	}

	response := strings.Split(writer.Body.String(), "}")

	var readTransactions []model.Transaction

	for _, v := range response {
		strTransaction := v + "}"
		var trans model.Transaction
		json.Unmarshal([]byte(strTransaction), trans)

		readTransactions = append(readTransactions, trans)
	}

	for _, v := range testData {
		for _, createdTransaction := range v.transactions {

			found := false

			for _, readTransaction := range readTransactions {

				if createdTransaction == readTransaction {
					found = true
				}
			}

			if found == false {
				err = fmt.Errorf("VerifyTransactions: Transaction: bic: %v iban: %v not found",
					createdTransaction.BIC,
					createdTransaction.IBAN)
				return
			}
		}
	}

	err = nil
	return
}

func WaitForUpdater() {
	time.Sleep(time.Second * time.Duration(config.TestConfiguration.UpdaterInterval))
}

// generate Random String from accountRunes slice
func RandBIC() string {

	bicRunes := []rune(config.TestConfiguration.BicRunes)

	b := make([]rune, 11)
	num := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range b {
		b[i] = bicRunes[num.Intn(len(bicRunes))]
	}
	return string(b)
}

func RandIBAN() string {
	iban := "DE"
	iban += fmt.Sprintf("%v%v", RandNumberWithRange(100000000, 999999999), RandNumberWithRange(100000000, 999999999))
	return iban
}

func RandNumberWithRange(low, hi int) int32 {
	num := low + rand.Intn(hi-low)
	return int32(num)
}
