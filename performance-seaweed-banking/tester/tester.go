package tester

import (
	"log"
	"time"
)

type Tester interface {
	RunFor(seconds time.Duration)
}

type BasicTester struct {
	baseURL    string
	cookie     int
	goroutines int
}

func NewBasicTester(baseurl string, cookie int, goroutines int) *BasicTester {
	return &BasicTester{baseURL: baseurl, cookie: cookie, goroutines: goroutines}
}

func listen(success chan bool, failure chan bool, quit <-chan time.Time) (int, int) {
	numberOfSuccess := 0
	numberOfFailure := 0

L:
	for {
		select {
		case <-success:
			numberOfSuccess++
		case <-failure:
			numberOfFailure++
		case <-quit:
			break L
		default:
		}
	}

	return numberOfSuccess, numberOfFailure
}

func printResults(name string, success int, failure int, duration time.Duration) {
	log.Println()
	log.Printf("Running tests for %v seconds for %s", duration.Seconds(), name)
	log.Println("-------------------------------------------------------------")
	log.Printf("Total requests: \t\t\t\t\t%d", success+failure)
	log.Printf("Total successful requests: \t\t\t\t%d", success)
	log.Printf("Total failed requests: \t\t\t\t%d", failure)
	log.Println()
	log.Printf("Successful requests per second \t\t\t%f", float64(success)/duration.Seconds())
	log.Printf("Failed requests per second \t\t\t\t%f", float64(failure)/duration.Seconds())
	log.Println()
	log.Printf("Time for a single request \t\t\t\t%fms", (duration.Seconds()*1000)/float64(success+failure))
}
