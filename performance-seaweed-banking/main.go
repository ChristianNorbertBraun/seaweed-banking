package main

import (
	"flag"
	"math/rand"

	"github.com/ChristianNorbertBraun/performance-seaweed-banking/tester"

	"log"
	"time"
)

var server = flag.String("server", "http://localhost:3333", "The address of the seaweed-banking-backend")
var duration = flag.Duration("duration", 10*time.Second, "Number of seconds the test should run")
var goroutines = flag.Int("goroutines", 2, "Number of goroutines")
var waiting = flag.Duration("waiting", 20*time.Second, "Time to wait for updater")
var noWaiting = flag.Bool("noWaiting", false, "No wating for updater")
var transactions = flag.Int("transactions", 50, "Number of transactions created per account")
var cookie = flag.Int("cookie", 0, "Random cookie")

func init() {
	flag.Parse()
}

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if *cookie == 0 {
		*cookie = r.Intn(100000)
	}

	log.Println("Cookie: ", *cookie)

	baseTester := tester.NewBasicTester(*server, *cookie, *goroutines)
	accountTester := tester.NewAccountTester(*baseTester)
	transactionCreationTester := tester.NewTransactionCreationTester(*baseTester, *transactions)
	transactionReadTester := tester.NewTransactionReadTester(*baseTester)

	accountTester.RunFor(*duration)
	transactionCreationTester.RunFor(*duration)

	if !*noWaiting {
		log.Println()
		log.Println("Giving updater time to work")
		time.Sleep(*waiting)
		log.Println("Start reading transactions")
	}

	transactionReadTester.RunFor(*duration)
}
