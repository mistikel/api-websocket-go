package main

import (
	"github.com/mistikel/api-websocket-go/server"
)

func main() {
	s := server.New()
	s.Serve()
}
