package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	const port int = 8080
	li, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Panic(err)
	}
	defer li.Close()
	fmt.Println("Listening on port", port)

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Println(err)
		}

		go handle(conn)
	}
}

// handle takes the connection as argument and handles it
func handle(conn net.Conn) {
	err := conn.SetDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Println("CONN TIMEOUT")
	}
	// use bufio to read the connection
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println(ln)
		fmt.Fprintf(conn, "i heard you say: %s\n", ln)
	}
	defer conn.Close()

	fmt.Println("Connection closed.")
}
