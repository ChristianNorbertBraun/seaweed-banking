package main

import (
	"testing"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/config"
	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-backend/model"
)

var benchAccounts []model.Account

func initBenchData() {

	if len(benchAccounts) <= 0 {

		for i := 0; i < config.TestConfiguration.NoOfBenchAccounts; i++ {
			newAcc := CreateRandomAccount()
			PostAccount(newAcc)
			benchAccounts = append(benchAccounts, newAcc)
		}

		WaitForUpdater()

	}
}

func benchmarkPostAccounts(b *testing.B) {

	for i := 0; i < b.N; i++ {

		err := PostAccount(CreateRandomAccount())

		if err != nil {
			b.Error(err)
		}
	}
}

func benchmarkPostAccountsParallel(b *testing.B) {

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := PostAccount(CreateRandomAccount())

			if err != nil {
				b.Error(err)
			}
		}
	})
}

func benchmarkGetAccounts(b *testing.B) {
	var index = 0
	initBenchData()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		if index >= len(benchAccounts) {
			index = 0
		}

		_, err := GetAccount(benchAccounts[index])

		if err != nil {
			b.Error(err)
		}

		index++
	}

}

func benchmarkGetAccountsParallel(b *testing.B) {
	var index = 0
	initBenchData()

	b.ResetTimer()

	b.SetParallelism(5)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {

			if index >= len(benchAccounts) {
				index = 0
			}
			_, err := GetAccount(benchAccounts[index])

			if err != nil {
				b.Error(err)
			}

			index++
		}
	})
}

func benchmarkGetAllAccounts(b *testing.B) {

	initBenchData()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := GetAllAccounts()

		if err != nil {
			b.Error(err)
		}
	}
}

func benchmarkGetAllAccountsParallel(b *testing.B) {

	initBenchData()

	b.ResetTimer()
	b.SetParallelism(5)
	b.RunParallel(func(pb *testing.PB) {

		for pb.Next() {
			_, err := GetAllAccounts()

			if err != nil {
				b.Error(err)
			}
		}
	})
}

func benchmarkReadAndWriteAccounts(b *testing.B, readRatio, writeRatio int) {
	var index = 0
	initBenchData()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		if i%readRatio == 0 {

			if index >= len(benchAccounts) {
				index = 0
			}

			_, err := GetAccount(benchAccounts[index])

			if err != nil {
				b.Error(err)
			}

			index++

		}

		if i%writeRatio == 0 {
			newAcc := CreateRandomAccount()
			err := PostAccount(newAcc)

			if err != nil {
				b.Error(err)
			}
		}
	}
}

func BenchmarkPostAccounts(b *testing.B)              { benchmarkPostAccounts(b) }
func BenchmarkGetAllAccounts(b *testing.B)            { benchmarkGetAllAccounts(b) }
func BenchmarkGetAccounts(b *testing.B)               { benchmarkGetAccounts(b) }
func BenchmarkReadAndWriteAccounts90_10(b *testing.B) { benchmarkReadAndWriteAccounts(b, 1, 9) }
func BenchmarkReadAndWriteAccounts80_20(b *testing.B) { benchmarkReadAndWriteAccounts(b, 1, 4) }
func BenchmarkReadAndWriteAccounts60_40(b *testing.B) { benchmarkReadAndWriteAccounts(b, 2, 3) }
func BenchmarkReadAndWriteAccounts50_50(b *testing.B) { benchmarkReadAndWriteAccounts(b, 1, 1) }
func BenchmarkReadAndWriteAccounts40_60(b *testing.B) { benchmarkReadAndWriteAccounts(b, 3, 2) }
func BenchmarkReadAndWriteAccounts20_80(b *testing.B) { benchmarkReadAndWriteAccounts(b, 4, 1) }
func BenchmarkReadAndWriteAccounts10_90(b *testing.B) { benchmarkReadAndWriteAccounts(b, 9, 1) }

// With parallel testing, CreateRandomAccounts throws errors cause of math/rand mutex issue
// func BenchmarkGetAccountsParallel(b *testing.B)    { benchmarkGetAccountsParallel(b) }
// func BenchmarkGetAllAccountsParallel(b *testing.B) { benchmarkGetAllAccountsParallel(b) }
// func BenchmarkPostAccountsParallel(b *testing.B) { BenchmarkPostAccountsParallel(b) }
