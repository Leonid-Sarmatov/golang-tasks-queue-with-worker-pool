
HOST=localhost
PORT=8080
REQUEST_NUM=100

for arg in "$@"; do
    case $arg in
        --port=*)
            PORT="${arg#*=}"
            ;;
        --host=*)
            HOST="${arg#*=}"
            ;;
        --requests-numbers=*)
            REQUEST_NUM="${arg#*=}"
            ;;
    esac
done

for ((i=1; i<=REQUEST_NUM; i++)); do
  curl -X POST http://${HOST}:${PORT}/enqueue \
    -H "Content-Type: application/json" \
    -d "{
      \"id\": \"$i\",
      \"payload\": \"обработать данные $i\",
      \"max_retries\": 3
    }" &
  sleep 0.1
done
wait
echo "Все ${REQUEST_NUM} задач отправлены"