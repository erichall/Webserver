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

func read(client net.Conn) (string, error) {
    holder := make([]byte, 10)
    number, err := client.Read(holder)
    if (err != nil) {
        return "", errors.New("Error couldn't get how many bytes that will be sent.")
    }
    bytes, conerr := strconv.Atoi(string(holder[0:number]))
    if (conerr != nil) {
        return "", errors.New("Could not convert data to byte.")
    }
    message := ""   

    for bytes != 0 {
        letters, Rederr := client.Read(holder)
        if (Rederr != nil) {
            return "", errors.New("Error when reading from client.")
        }
        message += string(holder[0:letters])
        bytes--
    }

    return message, nil
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

func write(connection net.Conn, msg []byte) {
	size := binary.Size(msg)
	sendSize := size/10
	if (size % 10 != 0) {
		sendSize++
	}
	fmt.Println(sendSize)
	_, sizeError := connection.Write([]byte(strconv.Itoa(sendSize)))

	if sizeError != nil {
		fmt.Println(sizeError)
		os.Exit(1)
		return
	}

	for times := 1; times != sendSize; times++ {
		prev := 10*(times-1)
		_, timesError := connection.Write(msg[prev:(10*times)])
		if timesError != nil {
			fmt.Println(timesError)
			os.Exit(1)
			return
		}
	}
	_, restError := connection.Write(msg[10*(sendSize-1):])

	if restError != nil {
		fmt.Println(restError)
		os.Exit(1)
		return
	}
	
}

func client(){
	_, port := getAddr()
	fmt.Println(port)
	connection, err := net.Dial("tcp", ("130.237.223.33" + ":" + "2345")) //connection is conn object


	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for  {
		msg, error := reader.ReadBytes('\n')
		fmt.Println(msg)
		if( error != nil){
			fmt.Println(err)
			os.Exit(1)
			return
		}	
		write(connection,msg)

		incmsg := make([]byte, 10)
		bytes,_ := connection.Read(incmsg)
		stringMsg := string(incmsg[:bytes])
		fmt.Println(stringMsg)
	
		
	}
}
