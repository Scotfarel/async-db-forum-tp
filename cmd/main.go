package main

import (
	"github.com/Scotfarel/db-tp-api/internal/app/server"
)

func main() {
	serverApi := server.Server{
		Url:	":5000",
	}
	serverApi.StartApiServer()
}
