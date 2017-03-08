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

1. Start postgres, seaweedfs and mongodb
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

3. Start seaweed-banking-account-updater. Therefor you have to go into the seaweed-banking-account-updater folder.
```
// default config
go run *.go
// custom config
go run *.go --config config/path/config.json
```

4. Now you are good to go.

**Info**
You may have to add your docker containers ip adress to your `/etc/hosts` file with the name `docker`

## Requests
Normaly you will only communicate with the seaweed-banking-backend.

```
// create account
curl --data '{"bic": "1234","iban": "iban1234","balance": 123}' localhost:3333/accounts

// create transactions
 curl --data '{"iban": "iban1234","bic": "1234","currency": "EUR","valueInSmallestUnit": 100,"intendedUse": "MoneyMoney"}' localhost:3333/accounts/1234/iban1234/transactions

// read accountinfo
curl localhost:3333/accounts/1234/iban1234/transactions?from=2017-02-16_13:05:00
```
