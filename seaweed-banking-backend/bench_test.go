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

func benchmarkCreateAccountsParallel(b *testing.B) {

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

	initBenchData()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		_, err := GetAccount(benchAccounts[RandNumberWithRange(0, len(benchAccounts))])

		if err != nil {
			b.Error(err)
		}
	}
}

func benchmarkReadAccountsParallel(b *testing.B) {

	initBenchData()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := GetAccount(benchAccounts[RandNumberWithRange(0, len(benchAccounts))])

			if err != nil {
				b.Error(err)
			}
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

	initBenchData()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		if i%readRatio == 0 {
			_, err := GetAccount(benchAccounts[RandNumberWithRange(0, len(benchAccounts))])

			if err != nil {
				b.Error(err)
			}
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

func BenchmarkCreateAccounts(b *testing.B)            { benchmarkPostAccounts(b) }
func BenchmarkGetAllAccounts(b *testing.B)            { benchmarkGetAllAccounts(b) }
func BenchmarkGetAccounts(b *testing.B)               { benchmarkGetAccounts(b) }
func BenchmarkReadAndWriteAccounts90_10(b *testing.B) { benchmarkReadAndWriteAccounts(b, 1, 9) }
func BenchmarkReadAndWriteAccounts80_20(b *testing.B) { benchmarkReadAndWriteAccounts(b, 1, 4) }
func BenchmarkReadAndWriteAccounts60_40(b *testing.B) { benchmarkReadAndWriteAccounts(b, 2, 3) }
func BenchmarkReadAndWriteAccounts50_50(b *testing.B) { benchmarkReadAndWriteAccounts(b, 1, 1) }
func BenchmarkReadAndWriteAccounts40_60(b *testing.B) { benchmarkReadAndWriteAccounts(b, 3, 2) }
func BenchmarkReadAndWriteAccounts20_80(b *testing.B) { benchmarkReadAndWriteAccounts(b, 4, 1) }
func BenchmarkReadAndWriteAccounts10_90(b *testing.B) { benchmarkReadAndWriteAccounts(b, 9, 1) }

// Parallel testing throws a lot of errors,
// func BenchmarkCreateAccountsParallel(b *testing.B) { BenchmarkCreateAccountsParallel(b) }
// func BenchmarkGetAllAccountsParallel(b *testing.B) { benchmarkGetAllAccountsParallel(b) }
// func BenchmarkReadAccountsParallel(b *testing.B)   { benchmarkReadAccountsParallel(b) }
