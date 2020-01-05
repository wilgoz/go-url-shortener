# URL Shortener
My first Go (1.13.5) program. Follows the best practices in minimizing module coupling with scalibility in mind.

Due to the lack of frontend, `curl` or `Postman` is needed to create a POST request to shorten the desired URL, with an example body as follows:
```json
{
	"original": "https://github.com/wilgoz"
}
```

## Features
*   URL Shortening
*   Configurable backend settings (choose either Redis or MongoDB for the backend)
*   Optional cache layer on top of MongoDB

## Installing Dependencies
*   `go mod download`

## Running
*   `go run main.go`
