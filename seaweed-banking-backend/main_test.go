package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/handler"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/model"
	"github.com/gorilla/mux"
)

type fakeAccount struct {
	account      model.Account
	transactions []model.Transaction
}

var r *mux.Router
var writer *httptest.ResponseRecorder
var accountRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var NoOfTransactionsPerAccount = 1
var NoOfFakeAccounts = 4
var updaterInterval = 10

var testData []fakeAccount

func TestMain(m *testing.M) {
	setUp()
	initTestData()
	test := m.Run()
	os.Exit(test)
}

func setUp() {
	r = mux.NewRouter()
	r.HandleFunc("/accounts", handler.GetAllAccounts).Methods("GET")
	r.HandleFunc("/accounts", handler.CreateAccount).Methods("POST")
	r.HandleFunc("/accounts/{bic}/{iban}/transactions", handler.GetTransaction).Methods("GET")
	r.HandleFunc("/accounts/{bic}/{iban}/transactions", handler.CreateTransaction).Methods("POST")
	writer = httptest.NewRecorder()
}

func initTestData() {

	for i := 0; i < NoOfFakeAccounts; i++ {

		var newFakeAccount fakeAccount

		newFakeAccount.account.BIC = RandString(11)
		newFakeAccount.account.IBAN = RandString(24)
		newFakeAccount.account.Balance = 500

		testData = append(testData, newFakeAccount)
	}
}

/*
*	TEST CASES
 */

func TestAccountCreate(t *testing.T) {
	// if cap(testData) > 0 {

	for _, data := range testData {
		err := CreateAccount(data.account)

		if err != nil {
			t.Errorf("CreateAccount %v", err)
		}
	}

	WaitForUpdater()

	for _, data := range testData {
		err := VerifyAccount(data.account)

		if err != nil {
			t.Errorf("VerifyAccount %v", err)
		}
	}
}

func TestTransactionsCreate(t *testing.T) {
	fmt.Println("Len :" + strconv.Itoa(len(testData)))

	if len(testData) > 1 {

		for i := range testData {
			if len(testData) > i+1 {
				var newTransaction model.Transaction
				newTransaction.BIC = testData[i+1].account.BIC
				newTransaction.IBAN = testData[i+1].account.IBAN
				newTransaction.BookingDate = time.Now()
				newTransaction.Currency = model.EUR
				newTransaction.IntendedUse = "testingTransactions"
				newTransaction.ValueInSmallestUnit = RandTransactionUnit(testData[i].account.Balance)

				err := CreateTransaction(testData[i].account, newTransaction)

				if err != nil {
					t.Errorf("CreateTransaction %v", err)
				}

				fmt.Printf("transaction created for %v \n", testData[i].account.IBAN)
			}
		}

	} else {
		t.Errorf("We need at least 2 accounts for transactions, current number of accounts: %v", cap(testData))
	}

	WaitForUpdater()

	for i := range testData {
		if i > 0 {
			err := VerifyTransactions(testData[i].account)

			if err != nil {
				t.Errorf("VerifyAccount Error: %v", err)
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
// 	// Readout All Transactions for every account parallel
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
		err = errors.New(request.URL.String() + " Code: " + strconv.Itoa(writer.Code))
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
		err = errors.New(request.URL.String() + " Code: " + strconv.Itoa(writer.Code))
		return
	}

	err = nil
	return
}

func GetAllAccounts() (readData []byte, err error) {

	request, _ := http.NewRequest("GET", "/accounts", nil)

	r.ServeHTTP(writer, request)

	if writer.Code != 200 && writer.Code != 201 {
		err = errors.New(request.URL.String() + " Code: " + strconv.Itoa(writer.Code))
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
		err = errors.New(request.URL.String() + " Code: " + strconv.Itoa(writer.Code))
		return
	}

	err = nil
	return
}

func VerifyTransactions(account model.Account) (err error) {

	request, _ := http.NewRequest("GET", "/accounts/"+account.BIC+"/"+account.IBAN+"/transactions", nil)

	r.ServeHTTP(writer, request)

	if writer.Code != 200 && writer.Code != 201 {
		err = errors.New(request.URL.String() + " Code: " + strconv.Itoa(writer.Code))

		return
	}
	//TODO Verifiy transactions
	err = nil
	return
}

func WaitForUpdater() {
	time.Sleep(time.Second * 10)
}

// generate Random String from accountRunes slice
func RandString(n int) string {
	b := make([]rune, n)

	for i := range b {
		b[i] = accountRunes[rand.Intn(len(accountRunes))]
	}

	return string(b)
}

func RandIBAN(n int) {
	//TODO

}

func RandNumber(n int) int32 {
	num := rand.New(rand.NewSource(time.Now().UnixNano()))
	return int32(num.Intn(n))
}

func RandTransactionUnit(max int32) int32 {
	rand.Seed(time.Now().Unix())
	unit := rand.Intn(int(max) - 0)
	return int32(unit)
}

//
// for j := 0; j < NoOfTransactionsPerAccount; j++ { // Parameterize noOfTransactions per account
//
// 	var newTransaction model.Transaction
//
// 	newTransaction.BIC = RandString(11)
// 	newTransaction.IBAN = RandString(24)
// 	newTransaction.BookingDate = time.Now()
// 	newTransaction.Currency = model.EUR
// 	newTransaction.IntendedUse = "BLUB"
// 	newTransaction.ValueInSmallestUnit = 150
//
// 	newFakeAccount.transactions = append(newFakeAccount.transactions, newTransaction)
// }
