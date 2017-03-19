package main

import (
	"encoding/json"
	"testing"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/model"
)

func benchmarkCreateAccounts(b *testing.B, n int) {
	accounts := make([]model.Account, n)

	for i := range accounts {
		accounts[i] = CreateRandomAccount()
	}

	b.StartTimer()
	for _, account := range accounts {
		PostAccount(account)
	}

}

func benchmarkCreateAccountParallel(b *testing.B) {

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			PostAccount(CreateRandomAccount())
		}
	})
}

func benchmarkReadAccounts(b *testing.B) {
	accounts, err := GetAllAccounts()

	if err != nil {
		b.Error(err)
	}

	readAccounts := make([]model.Account, 0)
	err = json.Unmarshal(accounts, &readAccounts)

	if err != nil {
		b.Error(err)
	}

	b.StartTimer()
	for _, acc := range readAccounts {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := GetAccount(acc)

				if err != nil {
					b.Error(err)
				}
			}
		})
	}

}

func BenchmarkCreateAccounts(b *testing.B) {
	b.StopTimer()
	benchmarkCreateAccounts(b, 100)
}

func BenchmarkReadAccounts(b *testing.B) {
	b.StopTimer()
	WaitForUpdater()
	benchmarkReadAccounts(b)
}

// func BenchmarkCreateAccountsParallel(b *testing.B) { benchmarkCreateAccountParallel(b) }
