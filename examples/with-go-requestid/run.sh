#!/bin/bash

set -x

go build -o testsrv .
./testsrv > server-logs.txt &

SERVER_PID=$!
ps -ef | grep $SERVER_PID
echo "Server started with PID $SERVER_PID"
sleep 2
make run-client

ps -ef | grep $SERVER_PID

kill -HUP $SERVER_PID
echo "Server stopped."

rm ./testsrv
