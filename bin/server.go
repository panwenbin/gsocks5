package main

import (
	"github.com/panwenbin/gsocks5/server"
	"log"
	"os"
	"strings"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("get env PORT failed, use 1080")
		port = ":1080"
	}
	if !strings.Contains(port, ":") {
		port = ":" + port
	}
	log.Printf("socks5 server listen at %s\n", port)
	server.ListenAndServe("tcp", port)
}
