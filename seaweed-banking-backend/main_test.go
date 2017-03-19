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
	"testing"
	"time"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/config"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/handler"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/model"
	"github.com/pressly/chi"
)

type fakeAccount struct {
	account      model.Account
	transactions []model.Transaction
}

var testAccount = model.NoBalanceAccount{Name: "TestUser", BIC: "TESTBIC", IBAN: "TESTIBAN"}

var r *chi.Mux
var testConfigPath = flag.String("testConfig", "./data/conf/testconfig.json", "Path to json formated testconfig")

func TestMain(m *testing.M) {

	setUp()
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

	r = chi.NewRouter()
	r.Get("/accounts", handler.GetAllAccounts)
	r.Post("/accounts", handler.CreateAccount)
	r.Get("/accounts/:bic/:iban", handler.GetAccount)
	r.Get("/accounts/:bic/:iban/transactions", handler.GetAccountInfo)
	r.Post("/accounts/:bic/:iban/transactions", handler.CreateTransactionAndUpdateBalance)
}

/*
*	TEST CASES
 */

func TestAccountsCreate(t *testing.T) {

	if testing.Short() {
		t.Skip("Skip Integration Tests if Short flag is set")
	}

	var testAccounts []model.Account

	for i := 0; i < config.TestConfiguration.NoOfFakeAccounts; i++ {
		newAccount := CreateRandomAccount()
		testAccounts = append(testAccounts, newAccount)

		err := PostAccount(newAccount)

		if err != nil {
			t.Error(err)
		}
	}

	WaitForUpdater()

	for _, account := range testAccounts {
		err := VerifyAccount(account)

		if err != nil {
			t.Error(err)
		}
	}
}

func TestTransactionsCreate(t *testing.T) {

	if testing.Short() {
		t.Skip("Skip Integration Tests if Short flag is set")
	}

	var testData []fakeAccount

	for i := 0; i < config.TestConfiguration.NoOfFakeAccounts; i++ {
		var newFakeAccount fakeAccount
		newFakeAccount.account = CreateRandomAccount()
		testData = append(testData, newFakeAccount)

		err := PostAccount(newFakeAccount.account)

		if err != nil {
			t.Error(err)
		}
	}

	if len(testData) > 1 {

		for i := range testData {
			if len(testData) > i+1 {
				newTransaction := CreateRandomTransaction(testData[i+1].account, "TestCreateTransaction", RandNumberWithRange(10, int(testData[i].account.Balance)))

				err := PostTransaction(testData[i].account, newTransaction)

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
			err := VerifyTransactions(testData[i])

			if err != nil {
				t.Error(err)
			}
		}
	}
}

/*
*	HELPERS
 */
func PostAccount(account model.Account) error {

	writer := httptest.NewRecorder()
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(&account)

	request, _ := http.NewRequest("POST", "/accounts", b)
	r.ServeHTTP(writer, request)

	if writer.Code != http.StatusCreated {
		return fmt.Errorf("CreateAccount: %v \nResponse Code: %v",
			request.URL.String(),
			writer.Code)
	}
	return nil
}

func PostTransaction(account model.Account, transaction model.Transaction) error {

	writer := httptest.NewRecorder()
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(&transaction)

	request, _ := http.NewRequest("POST", "/accounts/"+account.BIC+"/"+account.IBAN+"/transactions", b)
	r.ServeHTTP(writer, request)

	if writer.Code != http.StatusCreated {
		return fmt.Errorf("CreateTransaction: %v \nResponse Code: %v",
			request.URL.String(),
			writer.Code)
	}
	return nil
}

func GetAllAccounts() ([]byte, error) {

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/accounts", nil)

	r.ServeHTTP(writer, request)

	if writer.Code != http.StatusOK {
		return nil, fmt.Errorf("GetAllAccounts: %v \nResponse Code: %v",
			request.URL.String(),
			writer.Code)
	}
	return writer.Body.Bytes(), nil
}

func GetAccount(account model.Account) ([]byte, error) {

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/accounts/"+account.BIC+"/"+account.IBAN, nil)

	r.ServeHTTP(writer, request)

	if writer.Code != http.StatusOK {
		return nil, fmt.Errorf("GetAllAccounts: %v \nResponse Code: %v",
			request.URL.String(),
			writer.Code)
	}
	return writer.Body.Bytes(), nil
}

func VerifyAccount(account model.Account) error {

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/accounts/"+account.BIC+"/"+account.IBAN, nil)
	r.ServeHTTP(writer, request)

	if writer.Code != http.StatusOK {
		return fmt.Errorf("VerifyAccount: %v \nResponse Code: %v",
			request.URL.String(),
			writer.Code)
	}

	readAccount := model.Account{}
	if err := json.Unmarshal(writer.Body.Bytes(), &readAccount); err != nil {
		return fmt.Errorf("VerifyAccount: Unable to parse AccountInfo: bic %s iban %s",
			account.BIC,
			account.IBAN)
	}

	if account.BIC != readAccount.BIC ||
		account.IBAN != readAccount.IBAN ||
		account.Balance != readAccount.Balance {
		return fmt.Errorf("VerifyAccount: Account: bic: %v iban: %v not found",
			account.BIC,
			account.IBAN)
	}

	return nil
}

func VerifyTransactions(fakeAcc fakeAccount) error {

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/accounts/"+fakeAcc.account.BIC+"/"+fakeAcc.account.IBAN+"/transactions", nil)

	r.ServeHTTP(writer, request)

	if writer.Code != http.StatusOK {
		return fmt.Errorf("VerifyTransactions: %v \n Response Code: %v",
			request.URL.String(),
			writer.Code)
	}

	readAccountInfo := model.AccountInfo{}
	if err := json.Unmarshal(writer.Body.Bytes(), &readAccountInfo); err != nil {
		return fmt.Errorf("VerifyTransactions: Unable to parse AccountInfo: bic %s iban %s",
			fakeAcc.account.BIC,
			fakeAcc.account.IBAN)
	}

	for _, createdTransaction := range fakeAcc.transactions {
		for _, readTransaction := range readAccountInfo.Transactions {
			if createdTransaction.Recipient.BIC != readTransaction.Recipient.BIC ||
				createdTransaction.Recipient.IBAN != readTransaction.Recipient.IBAN ||
				createdTransaction.ValueInSmallestUnit != readTransaction.ValueInSmallestUnit ||
				createdTransaction.Currency != readTransaction.Currency ||
				createdTransaction.IntendedUse != readTransaction.IntendedUse {

				return fmt.Errorf("VerifyTransactions: Transaction: bic: %v iban: %v not found",
					createdTransaction.Recipient.BIC,
					createdTransaction.Recipient.IBAN)
			}
		}

	}
	return nil
}

func WaitForUpdater() {

	time.Sleep(time.Second * time.Duration(config.TestConfiguration.UpdaterInterval))
}

func CreateRandomAccount() model.Account {

	var newAccount model.Account

	newAccount.Name = fmt.Sprintf("RandomAccount%d", RandNumberWithRange(0, 1000))
	newAccount.BIC = RandBIC()
	newAccount.IBAN = RandIBAN()
	newAccount.Balance = RandNumberWithRange(200, 10000)

	return newAccount
}

func CreateRandomTransaction(targetAcc model.Account, intendedUse string, value int32) model.Transaction {

	var newTransaction model.Transaction
	newTransaction.Recipient = targetAcc.NoBalanceAccount
	newTransaction.Sender = testAccount
	newTransaction.Recipient.BIC = targetAcc.BIC
	newTransaction.Recipient.IBAN = targetAcc.IBAN
	newTransaction.BookingDate = time.Now()
	newTransaction.Currency = model.EUR
	newTransaction.IntendedUse = intendedUse
	newTransaction.ValueInSmallestUnit = value

	return newTransaction
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
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := low + r.Intn(hi-low)
	return int32(num)
}
