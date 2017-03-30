seaweed-banking
===============

This repo contains the components of a seaweedfs powered banking-accounting-system.

## seaweed-banking-backend
The seaweed-banking-backend is in charge of creating accounts and transactions. It also handles
all read requests. The seaweed-banking-backend is the only server which will receive requests from
the clients

## seaweed-banking-account-updater
The seaweed-banking-account-updater is in charge of preserving a consistent state for each account and
its transactions. It is invoked by the seaweed-banking-backend and collects all transactions for a given account 
and merges them into an account info file. 

## Getting started

1. Start postgres and seaweedfs
```
docker-compose up
```

2. Start seaweed-banking-backend. Therefor you have to go into the seaweed-banking-backend folder.
You can either start it with the default configuration in `./data/conf/config.json` or use your own
configuration.
```
// default config
go run *.go
// custom config
go run *.go --config config/path/config.json
```

3. Start seaweed-banking-account-updater. Therefore you have to go into the seaweed-banking-account-updater folder.
```
// default config start master
go run *.go --master
// custom config
go run *.go --config config/path/config.json --master

// to start slave
go run *.go --incomingConnections --port 8182
```

If there is no slave the master will handle all updates on its own

4. Now you are good to go.

**Info**
You may have to add your docker containers ip adress to your `/etc/hosts` file with the name `docker`

## Requests
Normaly you will only communicate with the seaweed-banking-backend.

```
// create account
curl --data '{"name":"Testuser","bic": "1234","iban": "iban1234","balance": 123}' localhost:3333/accounts

// create transactions
 curl --data '{"recipient":{"name":"Testuser","bic": "1234","iban": "iban1234"}, "sender": {"name":"Testuser2","bic": "12345","iban": "iban12345"} ,"currency": "EUR","valueInSmallestUnit": 100,"intendedUse": "MoneyMoney"}' localhost:3333/accounts/1234/iban1234/transactions

// read accountinfo
curl localhost:3333/accounts/1234/iban1234/transactions?from=2017-02-16_13:05:00
```

## Testing
At seaweed-banking you can either execute startTest.sh for automate integration testing, or use startSystem.sh to just automate system startup and test manually at seaweed-backing-backend

for testing manually several flags can configure the testing behaviour

execute `go test`

and optional:

```
// additional log information
-v 

// skip integration testing
-short

// execute all benchmark tests
-bench=.

// execute particular function by typing the FunctionName without "Benchmark"
// e.g. for testing BenchmarkReadAndWriteAccounts50_50():
 go test -v -short -bench=ReadAndWriteAccounts50_50

 // without additional configuration, Benchmark tests always process the amount of iterations
 // needed to satisfy the Benchmark runner (means like having a stable median of the time needed for each iteration or so)

 // It is also possible to set a fixed benchmark time; for time type e.g. 5s
 go test -v -short -bench={FunctionName or .} -benchtime={time}

```

## Performance

Running the `main.go` within the /performance-seaweed-banking folder will run some performance tests against a running seaweed-banking-backend.

```
// Run go run main.go -h for help
 -cookie int
    	Random cookie
  -duration duration
    	Number of seconds the test should run (default 10s)
  -goroutines int
    	Number of goroutines (default 2)
  -noWaiting
    	No wating for updater
  -server string
    	The address of the seaweed-banking-backend (default "http://localhost:3333")
  -transactions int
    	Number of transactions created per account (default 50)
  -waiting duration
    	Time to wait for updater (default 20s)

```





