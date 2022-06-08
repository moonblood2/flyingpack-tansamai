# service-shipping
The REST API service for shipping.

## Stacks
* Golang for the main language.
* Postgres for database management.
* http standard lib for handle REST API request/response.
* ozzo-validator for validate input.

## Design pattern and project layout.
* Use uncle Bob's Clean Architecture for design pattern. 

## Dev
### Run database
```docker-compose -f docker-compose.yml up -d postgres```
### Run Server
```go run ./cmd/restapi/main.go```
