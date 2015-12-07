package main

// Import packages.
import (
	"fmt"
	"net"
	"strconv"
	"errors"
//	"io/ioutil"
	"os"
	"encoding/binary"
	"bufio"
	"strings"
	
)

type User struct {
	card_number int
	first_name string
	last_name string
	sifferkod int
	enkod []int
	saldo int
}

func server(port int) {
    fmt.Println("Server startup!")
    stringPort := strconv.Itoa(port)
    // start listening to a part
    listener, err := net.Listen("tcp", ":" + stringPort)
    if (err != nil) {
        fmt.Println("Couldn't listen on port " + stringPort)
        return
    }
    for {
        connection, err := listener.Accept()
        fmt.Println("Connection established!")

         go forceShutDown(listener, connection) 

        if (err != nil) {
            fmt.Println("Failed to establish connection.")
            break
        } else {
            go handleClient(connection)
        }
    }
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


func handleClient(client net.Conn) {
	for {
		message, err := read(client)
	    
		if (err != nil) {
			fmt.Println(err)
			client.Close()
			break
		} else {
			fmt.Print(message)
			_, errWrite := client.Write([]byte("Message recieved!"))
			if (errWrite != nil) { 
                fmt.Println(errWrite) 
            }
		}
	}
}

func forceShutDown(listener net.Listener, client net.Conn) {
    for {
        var shutdown string
        fmt.Scanf("%s", &shutdown)

        if (shutdown == "shutdown") {
            listener.Close()
            client.Close()
            break
        }
    }
}

func main() {
    var port int
    _, err := fmt.Scanf("%d", &port)
    if (err != nil) {
        fmt.Println("Couldn't read user argument.")
    } else {
	    user, errUser := findUser("123456")
	    check(errUser)
	    fmt.Println(user)
	   // server(port)
    }
}

func check(e error){
	if e != nil {
		panic(e)
	}
}

func findUser(kortnummer string) (User, error){
	filen,err := os.Open("/home/erkan/Desktop/CDATA/PROGP/WEBSERVER/databas.txt")
	check(err)

	scanner := bufio.NewScanner(filen)

	
	fmt.Println(scanner.Text())

	var found_flag int = 0
	var stop_read int = 0

	
	var userdata [6]string 
	for scanner.Scan(){
		if stop_read < 6 {
			if scanner.Text() == kortnummer {
				found_flag = 1
			}

			if found_flag == 1 {
				userdata[stop_read] = scanner.Text()
				stop_read++
				
			}
		}
	}

	if len(userdata) == 0 {
		lol := new(User)
		return *lol, errors.New("Kan ikke hitta")
	}

	tmpenkod := strings.Split(userdata[4], " ")
	var intkodlist []int = make([]int, len(tmpenkod))
	

	for index, elem := range tmpenkod {
		intkodlist[index],_ = strconv.Atoi(elem)
	} 

	card_number,_ := strconv.Atoi(userdata[0])
	sifferkod,_ :=  strconv.Atoi(userdata[3])
	saldo,_ := strconv.Atoi(userdata[5])
	user := User{card_number,
		userdata[1],
		userdata[2],
		sifferkod,
		intkodlist,
		saldo}

	return user,nil
}
