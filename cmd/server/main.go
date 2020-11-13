package main

import "githib.com/aellwein/coap/v1/server"

func main() {
	srv := server.Builder().Build()
	if err := srv.Listen(); err != nil {
		panic(err)
	}
}
