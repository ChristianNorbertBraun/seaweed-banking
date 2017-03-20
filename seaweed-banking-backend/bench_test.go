package main

import "testing"

func benchmarkCreateAccounts(b *testing.B, n int) {

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// newAcc := CreateRandomAccount()
		// benchAccounts = append(benchAccounts, newAcc)
		err := PostAccount(CreateRandomAccount())

		if err != nil {
			b.Error(err)
		}
	}
}

func benchmarkCreateAccountParallel(b *testing.B) {

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		newAcc := CreateRandomAccount()
		benchAccounts = append(benchAccounts, newAcc)
		err := PostAccount(newAcc)

		if err != nil {
			b.Error(err)
		}
	})
}

func benchmarkReadAccounts(b *testing.B) {

	b.ResetTimer()

	for _, acc := range benchAccounts {
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
	benchmarkCreateAccounts(b, 10)
}

func BenchmarkReadAccounts(b *testing.B) {
	WaitForUpdater()
	benchmarkReadAccounts(b)
}

// func BenchmarkCreateAccountsParallel(b *testing.B) { benchmarkCreateAccountParallel(b) }
