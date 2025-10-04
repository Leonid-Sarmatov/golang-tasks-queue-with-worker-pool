for i in {1..100}; do
  curl -X POST http://localhost:8080/enqueue \
    -H "Content-Type: application/json" \
    -d "{
      \"id\": \"$i\",
      \"payload\": \"обработать данные $i\",
      \"max_retries\": 3
    }" &
done
wait
echo "Все 100 задач отправлены"