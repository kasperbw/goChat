package main

import "goChat/server"

func main() {
	server.New("./config.json")
	server.Run()
}
