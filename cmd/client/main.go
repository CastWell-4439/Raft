package client

import (
	"flag"
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	serveraddr := flag.String("serveraddr", "localhost:8080", "server address")
	command := flag.String("command", "", "command")
	flag.Parse()

	if *command == "" {
		fmt.Println("command is required")
		return
	}
	client, err := rpc.Dial("tcp", *serveraddr)
	if err != nil {
		log.Fatalln("fail to connect  %v", err)
	}
	defer client.Close()

}
