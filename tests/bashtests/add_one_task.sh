curl -X POST http://localhost:8080/enqueue \
  -H "Content-Type: application/json" \
  -d '{
    "id": "1",
    "payload": "обработать данные",
    "max_retries": 3
  }'