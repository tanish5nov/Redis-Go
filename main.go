package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("======STARTING REDIS=======")

	// starting to listen on Port 6379 using TCP connection
	listener, serverStartError := net.Listen("tcp", ":6379")
	if serverStartError != nil {
		fmt.Println("Error while starting Redis Server", serverStartError)
	}
	defer func(listener net.Listener) {
		if listener != nil {
			err := listener.Close()
			if err != nil {
				fmt.Println("Error while closing Redis Server", err)
			}
		}
	}(listener)

	// accepting connection on the listener which runs on Port 6379
	for {
		if listener == nil {
			break
		}
		connection, connectionAcceptError := listener.Accept()
		if connectionAcceptError != nil {
			fmt.Println("Error while accepting Connection", connectionAcceptError)
		}

		closingConnectionError := connection.Close()
		if closingConnectionError != nil {
			fmt.Println("Error while closing Connection", closingConnectionError)
		}
	}
}
