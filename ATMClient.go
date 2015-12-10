package main

// Import packages.
import (
	"bufio"
	"fmt"
	"os"
	"net"
	"strings"
)

//Global reader that will read what the user writes.
var reader *bufio.Reader = bufio.NewReader(os.Stdin)

//main function that will create the client.
func main() {
	client()
}

func check(err error) {
    if (err != nil) {
        fmt.Println(err)
        os.Exit(1)
    }
}

func format(rest []byte) []byte { 
    tmp := make([]byte, 10) 
    for index := range tmp { 
        if (index < len(rest)) { 
            tmp[index] = rest[index] 
        } else { 
            tmp[index] = 32 
        } 
    } 
    return tmp 
}

//read Reads from the server.
func write(connection net.Conn, msg []byte) {
    size := len(msg)
    sendSize := size/10
    if size % 10 != 0 {
        sendSize++
    }

    number := make([]byte, 1)
    number[0] = byte(sendSize)
    connection.Write(number)

    for times := 1; times != sendSize; times++ {
	    prev := 10*(times-1)
	    _, timesError := connection.Write(msg[prev:(10*times)])
	    check(timesError)
    }
     _, restError := connection.Write(format(msg[10*(sendSize-1):]))
    check(restError)
}

func read(connection net.Conn) string {
    times := make([]byte, 1)
    for {
        read, err := connection.Read(times)
        check(err)
        if read == 1 {
            break
        }
    }
    bytes := int(times[0])
    holder := make([]byte, 10)
    message := ""
    
    for bytes != 0 {
        letters, Rederr := connection.Read(holder)
        check(Rederr)
        message += string(holder[0:letters])
        bytes--
    }

    return strings.TrimSpace(message)
}

//getAddr gets the address.
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

func languageConfig(connection net.Conn) {
	intro := read(connection)
   
	fmt.Println(intro)
	for {
	//	fmt.Println("Enter loop")
		picked := string(userInput())
		picked = strings.ToLower(picked)
	//	fmt.Println("Send the picked language", picked)
		write(connection, []byte(picked))
		answer := strings.TrimSpace(read(connection))
		if answer == "approved" {
			//fmt.Println("Let's break!")
			break
		} else {
			fmt.Println(answer)
		}
	} 
	//fmt.Println("Out of loop!")   
}

func userInput() []byte {
	msg, err := reader.ReadBytes('\n')
	check(err)
	if len(msg) == 0 {
		msg = append(msg, 32)
	}
	//fmt.Println(msg)
	return msg
}

func loginSetUp(connection net.Conn) {
	intro := read(connection)
	cardText := read(connection)
	fmt.Println(intro)
	fmt.Print(cardText)
	
	for {
		card := strings.Replace(string(userInput()), " ", "", -1) //Läser in kortnr
		write(connection, []byte(card))//Skickar iväg kortnr
		
		ans := strings.TrimSpace(read(connection))
		if ans == "approved" {
			break
		} else {
			fmt.Print(ans)
		}
	}
	
	fmt.Print(read(connection)) //Skriver ut password:
	for {
		password := userInput()
		write(connection, []byte(password))
	    
		ans := strings.TrimSpace(read(connection))
		if ans == "approved" {
			fmt.Println(read(connection))
			break
		} else {
			fmt.Print(ans)
		}
	}
}



func handlingRequests(connection net.Conn) {    
	for {
		
		intro := read(connection) //What would you like to do?
		banner := read(connection) // our banner
		options := read(connection) //(1)Balance, ...
		fmt.Println(intro)
		fmt.Println(banner)
		fmt.Println(options)
		input := string(userInput())
		write(connection, []byte(input))
		input = strings.TrimSpace(input)
		switch input {
		case "1": //saldo
			fmt.Println(read(connection))
		case "2": //withdraw
			
			ans := strings.TrimSpace(read(connection)) 
			if ans != "approved" {
				fmt.Println(ans)
				break
			}		
			fmt.Print(read(connection))
			
			amount := string(userInput())
			write(connection, []byte(amount))
			
			fmt.Println(read(connection))
			write(connection, []byte(userInput()))
			
			app := strings.TrimSpace(read(connection))
			if app != "approved" {
				//fmt.Println("I am not approved")
				fmt.Println(app)
			}
		case "3": //deposit
			fmt.Print(read(connection))
			amount := string(userInput())
			write(connection, []byte(amount))
		case "4": //exit
			fmt.Println(read(connection))
			os.Exit(0)
		case "5": //change lang
			languageConfig(connection)	
		default:
			fmt.Println(read(connection))           
		}
	}
}

//client starts the client.
func client(){
    // Connect to the server through tcp/IP.
	connection, err := net.Dial("tcp", ("127.0.0.1" + ":" + "5678"))
    // If connection failed crash.
	check(err)
    //Configure the language.
    languageConfig(connection)
    //Time to log in to the account.
    loginSetUp(connection)
    handlingRequests(connection)
}
