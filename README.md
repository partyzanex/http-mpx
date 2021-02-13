# Тестовое задание HTTP-мультиплексор

Приложение представляет собой http-сервер с одним хендлером,
хендлер на вход получает POST-запрос со списком url в json-формате. 
Сервер запрашивает данные по всем этим url и возвращает результат клиенту в json-формате.
Eсли в процессе обработки хотя бы одного из url получена ошибка,
обработка всего списка прекращается и клиенту возвращается текстовая ошибка

## Ограничения
* для реализации задачи следует использовать Go 1.13 или выше
* использовать можно только компоненты стандартной библиотеки Go
* сервер не принимает запрос если количество url в нем больше 20
* сервер не обслуживает больше чем 100 одновременных входящих http-запросов
* для каждого входящего запроса должно быть не больше 4 одновременных исходящих
* таймаут на запрос одного url - секунда
* обработка запроса может быть отменена клиентом в любой момент, 
  это должно повлечь за собой остановку всех операций связанных с этим запросом
* сервис должен поддерживать 'graceful shutdown'
* результат должен быть выложен на github

## Решение

### Структура проекта

```bash
.
├── build                   # директория для билдов
├── cmd                     # для CLI
    ├── http-server         # CLI запуска HTTP-сервера
├── internal                
    ├── assert              # реализация assert.Equal
├── pkg
    ├── fetcher
    ├── limiter
    ├── pool
    ├── types
```

### Сборка и запуск
```bash
# запуск тестов
make test
# билд
make build
# запуск
./build/http-server
# для получения информации по параметрам
./build/http-server -h

# или в dev-режиме
go run ./cmd/http-server/
go run ./cmd/http-server/ -h
```

### Сборка и запуск в Docker
```bash
# билд
make docker-build
# запуск в контейнере
make docker-run

# удаление контейнера и удаление образа
make clean
```

### Отправка запроса

```bash
curl --location --request POST 'localhost:3000/' \
--header 'Content-Type: application/json' \
--data-raw '[
    {
        "url": "http://yandex.ru/"
    },
    {
        "url": "http://ya.ru/"
    },
    {
        "url": "http://ya.ru/",
        "method": "POST",
        "headers": {
            "X-Frame-Options": [
                "DENY"
            ]
        },
        "body": "WzEsMiwzXQ=="
    },
    {
        "url": "http://google.com/"
    }
]'
```
Ответ:
```json
[
    {
        "url": "https://yandex.ru/",
        "status_code": 200,
        "headers": {
            "Accept-Ch": [
                "Viewport-Width, DPR, Device-Memory, RTT, Downlink, ECT"
            ],
            "Accept-Ch-Lifetime": [
                "31536000"
            ],
            "Cache-Control": [
                "no-cache,no-store,max-age=0,must-revalidate"
            ],
            ...
        },
        "body": "PCFET0NUWVBFIGh0bWw...HRtbD4="
    },
    ...
    {
        "url": "http://ya.ru/",
        "status_code": 403,
        "headers": {
            "Content-Type": [
                "text/html; charset=utf-8"
            ],
            "Date": [
                "Sat, 13 Feb 2021 01:06:07 GMT"
            ],
            "Etag": [
                "W/\"60254eed-3077\""
            ],
            ...
        },
        "body": "PCFET0NUWVBFIGh0bWw...HRtbD4="
    },
    ...
]
```
