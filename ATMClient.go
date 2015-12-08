package main

// Import packages.
import (
	"bufio"
	"fmt"
	"os"
	"net"
	"strings"
)

//Global variable that indicate how many language the client support.
var languages = [...]string {"English","日本語","Deutsch"}
var intro = [...]string {"Please pick a language.","言語を選択してください。","Bitte wählen Sie eine Sprache aus."}

//Global variable that hold the language that the user chose.
var text []string

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
    //check(introerr)
    fmt.Println(intro)
    for {
        picked := string(userInput())
        picked = strings.ToLower(picked)
        write(connection, []byte(picked))
        answer := read(connection)
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
        card := userInput()
        write(connection, []byte(card))
        fmt.Print(read(connection))
        password := userInput()
        write(connection, []byte(password))
    
        ans := read(connection)
        if ans == "approved" {
            fmt.Println(read(connection))
            break
        } else {
            fmt.Print(ans)
        }
    }
}

func action(userInput string) {
    
}

func handlingRequests(connection net.Conn) {
    intro := read(connection) //What would you like to do?
    options := read(connection) //(1)Balance, ...    
    for {
        fmt.Println(intro)
        fmt.Println(options)
        input := strings.TrimSpace(string(userInput()))
        write(connection, []byte(input))
        switch input {
            case "1": //saldo
                fmt.Println(read(connection))
            case "2": //withdraw
                fmt.Print(read(connection))
                amount := strings.TrimSpace(string(userInput()))
                write(connection, []byte(amount))
            case "3": //deposit
                fmt.Print(read(connection))
                amount := strings.TrimSpace(string(userInput()))
                write(connection, []byte(amount))
            case "4": //exit
                fmt.Println(read(connection))
                os.Exit(0)
            default:
                fmt.Println(read(connection))           
        }
    }
}

//client starts the client.
func client(){
    // Connect to the server through tcp/IP.
	connection, err := net.Dial("tcp", ("130.229.156.183" + ":" + "2333"))
    // If connection failed crash.
	check(err)
    //Configure the language.
    languageConfig(connection)
    //Time to log in to the account.
    loginSetUp(connection)
    handlingRequests(connection)
}
