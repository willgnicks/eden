package main

import (
	"flag"
	"github.com/willgnicks/eden/client"
)

func main() {
	host := flag.String("h", "127.0.0.1", "the host name of chat server")
	port := flag.Int("p", 8888, "the port number of chat server")
	username := flag.String("u", "user_001", "the username to be logged in as")
	flag.Parse()
	c := client.New(*host, *port, *username)
	c.Spin()
}
