#!/usr/bin/env bash

HOST=localhost
PORT=30080
REQUEST_NUM=100
SLEEP_TIME=1.0
MAX_ATTEMPTS=5

if ! HOST=$(minikube ip) || [[ -z "$HOST" ]]; then
    echo "failed to get Minikube IP" >&2
    exit 1
fi

for arg in "$@"; do
    case $arg in
        --port=*)
            PORT="${arg#*=}"
            ;;
        --requests-numbers=*)
            REQUEST_NUM="${arg#*=}"
            ;;
        --sleep-time=*)
            SLEEP_TIME="${arg#*=}"
            ;;
    esac
done

echo 
echo "Connect to:" 
echo "    http://${HOST}:${PORT}/enqueue"
echo 

for ((i=1; i<=REQUEST_NUM; i++)); do
  curl --noproxy '*' -v --connect-timeout 10 http://${HOST}:${PORT}/enqueue \
    -H "Content-Type: application/json" \
    -d "{
      \"id\": \"$i\",
      \"payload\": \"обработать данные $i\",
      \"max_attempts\": ${MAX_ATTEMPTS}
    }" &
  sleep ${SLEEP_TIME}
done
wait
echo "Все ${REQUEST_NUM} задач отправлены"
