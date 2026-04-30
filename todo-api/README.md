#### create project
```bash
mkdir todo-api  
go mod init todo-api
```
create `main.go` file

## build project
```bash
go build
```

## run api
```bash
./todo-api
```

## Requests
get todos:
`curl http://localhost:8080/todos`

post todos:
`curl -X POST http://localhost:8080/todos -H "Content-Type: application/json" -d '{"title":"Go lernen","done":false}'`