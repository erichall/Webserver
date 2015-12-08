package main

// Import packages.
import (
	"bufio"
	"fmt"
	"os"
	"net"
	"strings"
	"encoding/binary"
	"strconv"
    "errors"
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

//read Reads from the server.
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

//write Writes to the server.
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

func languageConfig() {
    var lang string
    
    for {
        for _,elem := range intro {
            fmt.Println(elem)
        }
        fmt.Println()
        for _,elem := range languages {
            fmt.Println(elem)
        }
        fmt.Scanf("%s", &lang)
        lang = strings.TrimSpace(strings.ToLower(lang))
        valid := false
        
        for _,elem := range languages {
            if (strings.ToLower(elem) == lang) {
                valid = true
            }
        }
        if (valid) {
            lang = lang + ".txt"
            file, err := os.Open(lang)
            check(err)

            var lines []string
            scanner := bufio.NewScanner(file)
            for scanner.Scan() {
		        lines = append(lines, scanner.Text())
	        }
            file.Close()
	        check(scanner.Err())
            text = lines

            break
        } else {
            fmt.Println("Sorry. I did not understand.")
        }
    }
}

func userInput() string {
    msg, err := reader.ReadBytes('\n')
	check(err)
    return string(msg)
}

//client starts the client.
func client(){
    // Connect to the server through tcp/IP.
	//connection, err := net.Dial("tcp", ("130.237.223.33" + ":" + "2345"))
    // If connection failed crash.
	//check(err)
    //Configure the language.
    languageConfig()
    //Wait for server to send how the user can login.
    //login := read(connection)

    fmt.Println(text[0])
    fmt.Print(text[1])
    cardnumber := userInput()
    fmt.Println()
    fmt.Print(text[2])
    password := userInput()
    
    fmt.Println()
    fmt.Println(text)
    fmt.Println("Cardnumber = ", cardnumber)
    fmt.Println("Password = ", password)
}
