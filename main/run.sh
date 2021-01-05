#!/bin/bash
trap "rm server;kill 0" EXIT

go build -o server
./server -port=8001 &
./server -port=8002 &
./server -port=8003 -api=1 &

sleep 2
echo ">>> start test"
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Sam" &
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Jack" &
curl "http://localhost:9999/api?key=Sam" &
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Jack" &
curl "http://localhost:9999/api?key=Sam" &
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Jack" &
curl "http://localhost:9999/api?key=Sam" &
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Jack" &
curl "http://localhost:9999/api?key=Sam" &
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Jack" &
curl "http://localhost:9999/api?key=Sam" &
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Jack" &
curl "http://localhost:9999/api?key=Sam" &
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Jack" &
curl "http://localhost:9999/api?key=Sam" &
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Jack" &
curl "http://localhost:9999/api?key=Sam" &
curl "http://localhost:9999/api?key=Tom1" &
curl "http://localhost:9999/api?key=Sam" &

# 测试后发现，在同一短时刻内，多个请求只会响应一次，所以应该一次请求完成的时间，响应只会执行一次，下次请求时间，同一响应时间同样只执行一次
wait