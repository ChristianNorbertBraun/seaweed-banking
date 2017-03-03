#!/bin/bash


osascript -e 'do shell script "open -a Terminal ./restartBackend.sh"'

exec docker ps -a | awk '{ print $1,$2 }' | grep mongo | awk '{print $1 }' | xargs -I {} docker stop {}
exec docker ps -a | awk '{ print $1,$2 }' | grep postgres:9.6 | awk '{print $1 }' | xargs -I {} docker stop {}
exec docker ps -a | awk '{ print $1,$2 }' | grep chrislusf/seaweedfs | awk '{print $1 }' | xargs -I {} docker stop {}

exec docker ps -a | awk '{ print $1,$2 }' | grep mongo | awk '{print $1 }' | xargs -I {} docker rm {}
exec docker ps -a | awk '{ print $1,$2 }' | grep postgres:9.6 | awk '{print $1 }' | xargs -I {} docker rm {}
exec docker ps -a | awk '{ print $1,$2 }' | grep chrislusf/seaweedfs | awk '{print $1 }' | xargs -I {} docker rm {}

echo go to project
cd $GOPATH/src/github.com/ChristianNorbertBraun/seaweed-banking
echo start docker container...
exec docker-compose up
