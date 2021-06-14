package main

import (
	"fmt"
	"github.com/szaher/logit/internal/server"
	"log"
)

func main()  {
	fmt.Println("Welcome to my application")
	srv := server.NewHttpServer("0.0.0.0:8888")
	log.Fatal(srv.ListenAndServe())
}
