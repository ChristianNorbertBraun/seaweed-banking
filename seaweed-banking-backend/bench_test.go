package main

import "testing"

func benchmarkCreateAccounts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := PostAccount(CreateRandomAccount())

		if err != nil {
			b.Error(err)
		}
	}
}

func benchmarkCreateAccountParallel(b *testing.B) {

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := PostAccount(CreateRandomAccount())

			if err != nil {
				b.Error(err)
			}
		}
	})
}

func benchmarkReadAccounts(b *testing.B) {

	for i := 0; i < b.N; i++ {

		_, err := GetAccount(benchAccounts[RandNumberWithRange(0, len(benchAccounts)-1)])

		if err != nil {
			b.Error(err)
		}
	}
}

func benchmarkReadAccountsParallel(b *testing.B) {

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := GetAccount(benchAccounts[RandNumberWithRange(0, len(benchAccounts)-1)])

			if err != nil {
				b.Error(err)
			}
		}
	})
}

func benchmarkGetAllAccounts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GetAllAccounts()

		if err != nil {
			b.Error(err)
		}
	}
}

func benchmarkGetAllAccountsParallel(b *testing.B) {

	b.RunParallel(func(pb *testing.PB) {

		for pb.Next() {
			_, err := GetAllAccounts()

			if err != nil {
				b.Error(err)
			}
		}
	})
}

func benchmarkReadAndWriteAccounts(b *testing.B, writeRatio, readRatio int) {
	for i := 0; i < b.N; i++ {

		if i%readRatio == 0 {
			_, err := GetAccount(benchAccounts[RandNumberWithRange(0, len(benchAccounts)-1)])

			if err != nil {
				b.Error(err)
			}
		}

		if i%writeRatio == 0 {
			err := PostAccount(CreateRandomAccount())

			if err != nil {
				b.Error(err)
			}

		}

	}
}

func BenchmarkCreateAccounts(b *testing.B)            { benchmarkCreateAccounts(b) }
func BenchmarkGetAllAccounts(b *testing.B)            { benchmarkGetAllAccounts(b) }
func BenchmarkReadAccounts(b *testing.B)              { benchmarkReadAccounts(b) }
func BenchmarkReadAndWriteAccounts50_50(b *testing.B) { benchmarkReadAndWriteAccounts(b, 1, 1) }
func BenchmarkReadAndWriteAccounts10_90(b *testing.B) { benchmarkReadAndWriteAccounts(b, 1, 9) }
func BenchmarkReadAndWriteAccounts90_10(b *testing.B) { benchmarkReadAndWriteAccounts(b, 9, 1) }

// Parallel testing throws a lot of errors,
// func BenchmarkCreateAccountsParallel(b *testing.B) { BenchmarkCreateAccountsParallel(b) }
// func BenchmarkGetAllAccountsParallel(b *testing.B) { benchmarkGetAllAccountsParallel(b) }
// func BenchmarkReadAccountsParallel(b *testing.B)   { benchmarkReadAccountsParallel(b) }
