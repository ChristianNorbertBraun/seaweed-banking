#!/bin/bash
cd $GOPATH/src/github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater
sleep 10
exec go run main.go
