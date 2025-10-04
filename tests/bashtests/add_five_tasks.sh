curl -X POST http://localhost:8080/enqueue \
  -H "Content-Type: application/json" \
  -d '{
    "id": "1",
    "payload": "задача 1",
    "max_retries": 3
  }'

curl -X POST http://localhost:8080/enqueue \
  -H "Content-Type: application/json" \
  -d '{
    "id": "2",
    "payload": "задача 2",
    "max_retries": 5
  }'

curl -X POST http://localhost:8080/enqueue \
  -H "Content-Type: application/json" \
  -d '{
    "id": "3",
    "payload": "задача без возможности повторить выполнение",
    "max_retries": 0
  }'

curl -X POST http://localhost:8080/enqueue \
  -H "Content-Type: application/json" \
  -d '{
    "id": "4",
    "payload": "задача с некорректными данными",
    "max_retries": -10
  }'

curl -X POST http://localhost:8080/enqueue \
  -H "Content-Type: application/json" \
  -d '{
    "id": "5",
    "payload": "задача 5",
    "max_retries": 3
  }'