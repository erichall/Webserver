/*
Webserver
Koden är skriven av Johannes Westlund & Eric Hallström
Senast modifierad : 2015-12-08 
*/
package main

// Import packages.
import (
	"fmt"
	"net"
	"strconv"
	"errors"
	"os"
	"encoding/binary"
	"bufio"
	"strings"
	
)
/*User struct definierar all information om en user.*/
type User struct {
	card_number int
	first_name string
	last_name string
	sifferkod int
	enkod []int
	saldo int
}

/*Hanterar servern, 
input : Vilken port srv ska ligga på*/
func server(port int) {
	fmt.Println("Server startup!")
	stringPort := strconv.Itoa(port) 

	// start listening to a port
	listener, err := net.Listen("tcp", ":" + stringPort)
	if (err != nil) {
		fmt.Println("Couldn't listen on port " + stringPort)
		return
	}

	for {
		connection, err := listener.Accept() //Acceptera incomming connection
		fmt.Println("Connection established!") 

		go forceShutDown(listener, connection) //Goroutine för att stänga porten helt.

		if (err != nil) {
			fmt.Println("Failed to establish connection.")
			break
		} else {
			go handleClient(connection) //Egen goroutin för varje användare
		}
	}
}
/*
read läser inc data från clienter
input : en net.Conn dvs en connection
return : läst input oavsett längd, err om misslyckad
Clienten börjar alltid med att skicka längden på sitt meddelande i bytes
*/
func read(client net.Conn) (string, error) {
	holder := make([]byte, 10) //byte array för att spara inc data max inc byte är 9999 999 999 = ca 9.9 miljarder
	number, err := client.Read(holder) //number - antalet bytes som sändes
	if (err != nil) {
		return "", errors.New("Error couldn't get how many bytes that will be sent.")
	}
	
	bytes, conerr := strconv.Atoi(string(holder[0:number])) //bytes := tar nu ut siffran för antal bytes som sändes
	
	if (conerr != nil) {
		return "", errors.New("Could not convert data to byte.")
	}
	
	message := ""   

	//samlar ihop hela meddelandet till ett message, loopar bytes gånger som sändes.
	for bytes != 0 {
		letters, Rederr := client.Read(holder) //Läser från client
		
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
	    var users []User = findUser()
	    fmt.Println(users[0])
	   // server(port)
    }
}

func check(e error){
	if e != nil {
		panic(e)
	}
}

func findUser() ([]User){
	filen,err := os.Open("/home/erkan/Desktop/CDATA/PROGP/WEBSERVER/databas.txt")
	check(err)
	
	scanner := bufio.NewScanner(filen)

	
	fmt.Println(scanner.Text())


	var stop_read int = 0
	var user_info int = 6 //Hur många rader i databasen som består av en user.
	
	var userdata []string = make([]string, 6)

	var users []User

	
	for i := 1; scanner.Scan(); i++ {
		userdata[stop_read] = scanner.Text()
		stop_read++	
		if i % user_info == 0 {
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
		
			users = append(users, user)
			stop_read = 0
			userdata = make([]string, 6)
		}

		
	

	}

	

	return users
}
