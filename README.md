# HTTP multiplexer


## Описание
Приложение принимает POST-запрос на  `/multiplexer` c переданным массивом `url` 

Пример body запроса: 

```json
{
    "url": [
        "https://jsonplaceholder.typicode.com/todos/1",
        "https://jsonplaceholder.typicode.com/todos/2",
        "https://jsonplaceholder.typicode.com/todos/3",
        "https://jsonplaceholder.typicode.com/todos/4",
        "https://jsonplaceholder.typicode.com/todos/5"
    ]
}
```

Ответ возвращается в виде массива `key->value`, где `key` - это `url`, а `value` это `body` ответа в виде необработанной строки.

Пример ответа: 

```json
{
    "https://jsonplaceholder.typicode.com/todos/1": "{\n  \"userId\": 1,\n  \"id\": 1,\n  \"title\": \"delectus aut autem\",\n  \"completed\": false\n}",
    "https://jsonplaceholder.typicode.com/todos/2": "{\n  \"userId\": 1,\n  \"id\": 2,\n  \"title\": \"quis ut nam facilis et officia qui\",\n  \"completed\": false\n}",
    "https://jsonplaceholder.typicode.com/todos/3": "{\n  \"userId\": 1,\n  \"id\": 3,\n  \"title\": \"fugiat veniam minus\",\n  \"completed\": false\n}",
    "https://jsonplaceholder.typicode.com/todos/4": "{\n  \"userId\": 1,\n  \"id\": 4,\n  \"title\": \"et porro tempora\",\n  \"completed\": true\n}",
    "https://jsonplaceholder.typicode.com/todos/5": "{\n  \"userId\": 1,\n  \"id\": 5,\n  \"title\": \"laboriosam mollitia et enim quasi adipisci quia provident illum\",\n  \"completed\": false\n}"
}
``` 

Если во время запроса произошла ошибка, то ответ будет содержать поле `message`

Пример ответа с ошибкой: 

```json
{
    "message": "Get \"https://domainnotexists.com\": dial tcp: lookup domainnotexists.com on 127.0.0.53:53: no such host"
}
```

## Билд 

`make build`, затем `./multiplexer`

## Линтинг

Необходимо локально поставить пакет [golangci-lint](https://github.com/golangci/golangci-lint)

Затем вызвать `make lint`
