# tarantool-kv-storage

## API
- POST /kv body: {key: "test", "value": {SOME ARBITRARY JSON}} 
- PUT /kv/{id} body: {"value": {SOME ARBITRARY JSON}} 
- GET /kv/{id} 
- DELETE /kv/{id} 
- POST возвращает 409 если ключ уже существует
- POST, PUT возвращают 400 если body некорректное 
- PUT, GET, DELETE возвращает 404 если такого ключа нет

### Database
in directory database:
````
docker build -t tarantooldb .
docker run --rm -d -p 3301:3301 --name=tarantooldb tarantooldb
```` 

### Start server
````
go build -o bin/main cmd/main.go
bin/main
````
or (docker)
````
docker build -t service .
docker run --rm -d -p 80:80 --name=service service
````

### Tests
in directory pkg/handler
````
go test -v
````
you can generate html file that shows test coverage
````
go test -v -coverprofile="../../test/handler_cover.out"
go tool cover -html="../../test/handler_cover.out" -o "../../test/handler_cover.html"
````
