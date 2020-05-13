# api-websocket-go

## Prerequisite tools

- Go version go1.14.1 
- Docker (*optional*)


## Command

- Install Dependency
  ```
  go mod download
  ```

- Test
  ```
  go test -race ./...
  ```

- Run server
  ```
  go run main.go
  ```

- Run client
  ```
  cd client && go run client.go
  ```

- Run script request to publish message
  ```
  ./request.sh
  ```

- Run with docker
  ```
  docker build -t asg .
  docker run -p 8080:8080 asg
  ```
