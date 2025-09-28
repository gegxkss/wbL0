# wbL0 — Микросервис для обработки заказов Wildberries

## Описание
wbL0 — это микросервис на Go, реализующий обработку заказов, интеграцию с Kafka и хранение данных в PostgreSQL. Проект поддерживает кэширование, REST API, миграции и генерацию тестовых заказов.

## Стек технологий
- Go
- PostgreSQL
- Kafka
- GORM
- Docker
- HTML/JS/CSS (фронтенд)

## Структура проекта
```
├── cmd/                
├── front/              
├── internal/
│   ├── cache/         
│   ├── config/        
│   ├── handlers/       
│   └── models/         
├── kafka/              
├── migrations/         
├── producer/           
├── docker-compose.yml  
├── main.go             
```

## Быстрый старт
1. Запустите инфраструктуру:
   ```sh
   docker-compose up -d
   ```
2. Выполните миграции:
   ```sh
   go run main.go -m
   ```
3. Запустите микросервис:
   ```sh
   go run main.go
   ```
4. Откройте фронтенд:
   - Перейдите на [http://localhost:8081](http://localhost:8081)

## Тесты
Запуск всех unit-тестов:
```sh
go test ./...
```

## Авторы
- gegxkss

---
Для вопросов и предложений — открывайте issue или пишите напрямую!
