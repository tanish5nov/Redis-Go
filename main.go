package main

import (
	"BuildYourOwnRedis/QueryHandler"
	"BuildYourOwnRedis/RESPParser"
	"fmt"
	"net"
	"sync"
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
	wg := sync.WaitGroup{}
	// accepting connection on the listener which runs on Port 6379
	for {
		if listener == nil {
			break
		}
		connection, connectionAcceptError := listener.Accept()
		if connectionAcceptError != nil {
			fmt.Println("Error while accepting Connection", connectionAcceptError)
			continue
		}

		wg.Add(1)
		go func(connection net.Conn, wg *sync.WaitGroup) {
			defer wg.Done()
			for {
				buffer := make([]byte, 1024)
				n, err := connection.Read(buffer) // this is a blocking call, n = quantity of bytes arrived

				if err != nil {
					fmt.Println("Error while reading from Connection", err)
					err := connection.Close()
					if err != nil {
						fmt.Println("Error while closing Redis Server", err)
					}
					return
				}

				inputData := string(buffer[:n])

				queryResult, err := RESPParser.ParseFacade(inputData)
				if err == nil {
					fmt.Println("DEBUG 7: ", queryResult)
				} else {
					fmt.Println(err)
				}
				result := QueryHandler.QueryHandlerFacade(inputData)
				sendUpdate, err := connection.Write([]byte(result))

				if err != nil {
					fmt.Println("Error while sending to Connection", err)
					err := connection.Close()
					if err != nil {
						fmt.Println("Error while closing Redis Server", err)
					}
					return
				}

				fmt.Println("WRITE OUTPUT: ", sendUpdate)
			}
		}(connection, &wg)
	}
	wg.Wait()
}
