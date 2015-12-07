package main

import (
	"bufio"
	"fmt"
	"os"
	"net"
	//"strings"
	"encoding/binary"
	"strconv"
	//"math"
)

func main() {
	client()
}

func getAddr() (string, string) {
	var ip string
	var port string
	fmt.Print("Enter IP to connect to: ")
	_, err := fmt.Scanln(&ip)
	if err != nil {
		fmt.Println(err)
		return "", ""
	}

	fmt.Println("Enter Port: ")
	_, erro := fmt.Scanln(&port) 
	if erro != nil {
		fmt.Println(erro)
		return "", ""
	}
	
	return ip, port
}

func transform(msg []byte) []byte {
    tmp := make([]byte, 10)
    
    for index := 0; index < len(msg); index++ {
	    if (msg[index] == 10) {
		    tmp[index] = msg[index]
		    break
	    }
	    tmp[index] = msg[index]
    }
    return tmp
}

func write(connection net.Conn, msg []byte) {
	size := binary.Size(msg)
	sendSize := size/10
	if (size % 10 != 0) {
		sendSize++
	}
	fmt.Println(sendSize)
	connection.Write([]byte(strconv.Itoa(sendSize)))

	for times := 1; times != sendSize; times++ {
		prev := 10*(times-1)
		connection.Write(msg[prev:(10*times)])
	}
	connection.Write(msg[10*(sendSize-1):])
	
}

func client(){
	_, port := getAddr()
	fmt.Println(port)
	//connection, err := net.Dial("tcp", (ip + ":" + port))

	connection, err := net.Dial("tcp", ("130.237.223.33" + ":" + "2345")) //connection is conn object


	if err != nil {
		fmt.Println(err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for  {
		msg, error := reader.ReadBytes('\n')
		fmt.Println(msg)
		if( error != nil){
			fmt.Println(err)
			return
		}	
		write(connection,msg)

	//	incmsg := make([]byte, 10)
	//	bytes,_ := connection.Read(incmsg)
	//	stringMsg := string(incmsg[:bytes])
	//	fmt.Println(stringMsg)
	
		
	}
}
