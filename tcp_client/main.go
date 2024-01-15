package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	const port int = 8080
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	//userName := "Dennis Phantom"

	bs, err := io.ReadAll(conn)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(bs))

}
