#!/bin/bash
cd $GOPATH/src/github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater
sleep 20
exec go run main.go --master
