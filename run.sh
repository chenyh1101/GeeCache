#!/usr/bin/env bash

trap "rm server;kill 0" EXIT
go build  -o server
./server -port=8001 -api=1&
./server -port=8002 &
./server -port=8003 &



sleep 2
echo  ">>>start test"
curl "http://localhost:9999/api?key=Shen" &
curl "http://localhost:9999/api?key=Jack" &
curl "http://localhost:9999/api?key=Sam" &

wait