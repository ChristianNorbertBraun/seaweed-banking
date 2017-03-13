#!/bin/bash

#osascript -e 'do shell script "open -a Terminal ./restartBackend.sh"'
echo start Backend inside new Tab...
osascript -e 'tell application "Terminal" to activate' -e 'tell application "System Events" to tell process "Terminal" to keystroke "t" using command down'
sleep 1
osascript -e 'tell application "Terminal" to do script "./restartBackend.sh" in selected tab of window 1'

echo start Updater inside new Tab...
osascript -e 'tell application "Terminal" to activate' -e 'tell application "System Events" to tell process "Terminal" to keystroke "t" using command down'
sleep 1
osascript -e 'tell application "Terminal" to do script "./restartUpdater.sh" in selected tab of window 1'

echo go to project
cd $GOPATH/src/github.com/ChristianNorbertBraun/seaweed-banking

echo stop and restart docker container...
docker-compose stop && docker-compose rm -f && docker-compose up -d

echo Wait Updater and Backend...
sleep 18

echo start Tests...
cd seaweed-banking-backend
exec go test -v
