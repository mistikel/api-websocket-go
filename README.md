# warpin code-test

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
  go run test -race ./...
  ```

- Run server
  ```
  go run main.go
  ```

- Run client
  ```
  cd client && go run client.go
  ```

- Run with docker
  ```
  docker build -t warpin .
  docker run -p 8080:8080 warpin
  ```