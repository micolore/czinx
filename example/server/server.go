package main

import (
	"czinx/znet"
)

func main() {

	s := znet.NewServer("[zinx v-0.2]")

	s.Serve()
}
