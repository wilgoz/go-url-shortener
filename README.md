# URL Shortener
My first Go (1.13.5) program.

Due to the lack of a frontend, `curl` or `Postman` is needed to create a POST request to shorten the desired URL.

An example of a POST request body:
```json
{
    "original": "https://github.com/wilgoz"
}
```

An example of a result sent by the server after sending a POST request:
```json
{
    "created_at": 1578324327,
    "original": "https://github.com/wilgoz",
    "shortened": "c4HfogPZR"
}
```

Running `localhost:{port}/{shortened}` will then redirect the client to the original URL.

## Main Features
*   URL shortening
*   Configurable backend settings
*   Multiple storage support (Redis/MongoDB)
*   Optional redis cache layer on top of MongoDB

## Installing Dependencies
*   `go mod download`

## Running
*   `go run main.go`
