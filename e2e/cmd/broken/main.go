package main

import "time"

func main() {
	time.Sleep(5 * time.Second)
	panic("Error starting server")
}
