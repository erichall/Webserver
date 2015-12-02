package main

import (
	"bufio"
	"fmt"
	"os"
	"net"
	//"strings"
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



func client(){
	_, port := getAddr()
	//connection, err := net.Dial("tcp", (ip + ":" + port))
	connection, err := net.Dial("tcp", ("130.237.227.22" + ":" + port))

	if err != nil {
		fmt.Println(err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	
	for  {
	
	
		msg, error := reader.ReadString('\n')
		if( error != nil){
			fmt.Println(err)
			return
		}
		
		_,writeError := connection.Write([]byte(msg))

		
		if(writeError != nil){
			connection.Close()
			fmt.Println(writeError)
			return
		}

		incmsg := make([]byte, 10)
		bytes,_ := connection.Read(incmsg)
		stringMsg := string(incmsg[:bytes])
		
	//	if len(stringMsg) > 0 {
			fmt.Println(stringMsg)
	//	}
	}
}
