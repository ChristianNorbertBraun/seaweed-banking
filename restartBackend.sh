#!/bin/bash
cd $GOPATH/src/github.com/ChristianNorbertBraun/seaweed-banking

#osascript -e 'do shell script "open -a Terminal ./restartUpdater.sh"'
#osascript -e 'tell application "Terminal" to activate' -e 'tell application "System Events" to tell process "Terminal" to keystroke "t" using command down' -e 'tell application "Terminal" to do script "./restartUpdater.sh" in tab 2 of window 2'

sleep 20
cd seaweed-banking-backend
exec go run main.go
