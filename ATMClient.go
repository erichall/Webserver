package main

import (
	"bufio"
	"fmt"
	"os"
	"net"
	//"strings"
	"encoding/binary"
	"strconv"
	"math"
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

func write(connection net.Conn, msg []byte) {
	size := binary.Size(msg)
	//fmt.Println(size)
	var tmpMsg = make([]byte, 10)
	var looper int = 1
	var sendSize float64 = math.Ceil(float64(size/10))
	connection.Write([]byte(strconv.Itoa(int(sendSize))))
	
	if size <= 10 {
		connection.Write(msg)
	}else {
		for i := 1; i < size; i++ {
			tmpMsg[looper-1] += msg[i-1]
			fmt.Println(msg[i])
			if i % 10 == 0 {
				fmt.Println(tmpMsg)
				tmpMsg = make([]byte,10)
				looper = 0
				connection.Write(tmpMsg)
			}else if size-1 == i {
				fmt.Println(tmpMsg)
				connection.Write(tmpMsg)
			}
			looper += 1
		}
	}
}

func client(){
	_, port := getAddr()
	fmt.Println(port)
	//connection, err := net.Dial("tcp", (ip + ":" + port))
	connection, err := net.Dial("tcp", ("130.237.223.107" + ":" + "2345")) //connection is conn object

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
	//	_,writeError := connection.Write([]byte(msg))

		
	/**	if(writeError != nil){
			connection.Close()
			fmt.Println(writeError)
			return
		}
*/
		write(connection,msg)

		incmsg := make([]byte, 10)
		bytes,_ := connection.Read(incmsg)
		stringMsg := string(incmsg[:bytes])
		fmt.Println(stringMsg)
	
		
	}
}
