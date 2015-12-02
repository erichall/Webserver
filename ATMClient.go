package main

import (
	//"io"
	"bufio"
	"fmt"
	"os"
	"net"
	"encoding/gob"
)

func main() {
	client()
}

func connect() (string, string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter IP to connect to: ")
	ip, err := reader.ReadString('\n') //Read untill enter
	if err != nil {
		fmt.Println(err)
		return "", ""
	}

	fmt.Println("Enter Port: ")
	port, err := reader.ReadString('\n') //Read untill enter
	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	
	return ip, port
}



func client(){
	ip, port := connect()
	connection, err := net.Dial("tcp", (ip + ":" + port))

	if err != nil {
		fmt.Println(err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	
	for {
		msg, err := reader.ReadString('\n')
		if( err != nil){
			fmt.Println(err)
			return
		}

		err = gob.NewEncoder(connection).Encode(msg)
	}
}
