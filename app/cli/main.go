package main

import (
	_ "github.com/lib/pq"
	"github.com/typical-go/typical-rest-server/app/server"
	"github.com/typical-go/typical-rest-server/typical"
)

func main() {
	typical.Container().Invoke(func(s *server.Server) error {
		return s.Serve()
	})
}
