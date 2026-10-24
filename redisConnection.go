package main

import (
	"fmt"
	"net"
)

func main() {

	/*
		we use net package of go to start the tcp connection
		this returns a listener object
	*/
	listener, listeningError := net.Listen("tcp", ":6379")

	if listeningError != nil {
		fmt.Println("Listening error", listeningError)
		return
	}

	/*
		this should be done here only,
		coz consider a scenario where you placed it right after
		listener, listeningError := net.Listen("tcp", ":6379")
		in that case if we get error then we will first check
		if listeningError != nil {
		and then return, upon return defer listener.Close()
		will be called, coz of which it will be nil.close()
		which will give runtime errors
	*/
	defer listener.Close()

	/*
		this is intentionally kept as infinite loop,
		the reason being, if we don't do infinite,
		then one connection will be accepted and program will terminate
		listener.Accept() is a blocking call
	*/

	for {
		connection, connectionError := listener.Accept()
		if connectionError != nil {
			fmt.Println("Connection acceptance error", connectionError)
			fmt.Println(connection)
		}
	}
}
