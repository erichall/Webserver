package main

// Import packages.
import (
	"bufio"
	"fmt"
	"os"
	"net"
	"strings"
	"strconv"
)

//Global reader that will read what the user writes.
var (
	reader *bufio.Reader = bufio.NewReader(os.Stdin)
	languages = "english\n日本語\ndeutsch\nsvenska"
	intro = "Please pick a language.\n言語を選択してください。\nBitte wählen Sie eine Sprache aus.\nVar venlig velj ett sprak.\n"
	end = "Please pick a real language.\n実際の言語を選択してください。\nBitte wählen Sie eine echte Sprache.\nSnelle velj ett riktigt sprak.\n"

	lines []string
	custlang string
)

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

func validateLang() {
	fmt.Println(intro + "\n" + languages)
	lines = make([]string,0)
	l := strings.Split(languages, "\n")
	for {
		//fmt.Println("Entering loop")
		picked := strings.TrimSpace(userInput())
		
		for _,lang := range(l) {
			fmt.Println(lang)
			if (picked == lang) {
				fmt.Println("Found language!", lang, picked)
				custlang = lang
				lang = lang + ".txt"
			
				file, err := os.Open(lang)
				check(err)

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					lines = append(lines, scanner.Text())
				 }
				file.Close()
				check(scanner.Err())
				break
			} 
		}
		if (len(lines) != 0) {
			break
		} else {
			fmt.Println(end)
		}
	}
}

func userInput() string {
	msg, err := reader.ReadBytes('\n')
	check(err)
	if len(msg) == 0 {
		msg = append(msg, 32)
	}
	//fmt.Println(msg)
	return string(msg)
}

func stringToInt(strint string) int {
	tmp,_ := strconv.Atoi(strint)
	return tmp
}

func makeMsg(opt int, msg string) []byte {
	msg = strings.TrimSpace(msg)
	var res = make([]byte, 10)
	res[0] = byte(opt)
	
	switch opt {
	case 2 : //Withdrawl
		if len(msg) > 9 {
			break
		}
		for index := range msg {
			res[index+1] = byte(msg[index])
		}
	case 3 :
		if len(msg) > 9 {
			break
		}
		for index := range msg {
			res[index+1] = byte(msg[index])
		}
	case 100 : //cardnumber
		if len(msg) != 16 {
			break
		}
		res[1] =  byte(stringToInt(msg[0:2]))
		res[2] =  byte(stringToInt(msg[2:4]))
		res[3] =  byte(stringToInt(msg[4:6]))
		res[4] =  byte(stringToInt(msg[6:8]))
		res[5] =  byte(stringToInt(msg[8:10]))
		res[6] =  byte(stringToInt(msg[10:12]))
		res[7] =  byte(stringToInt(msg[12:14]))
		res[8] =  byte(stringToInt(msg[14:16]))
	case 101 : //password
		if len(msg) != 4 {
			break	
		}
		res[1] = byte(stringToInt(msg[0:1]))
		res[2] = byte(stringToInt(msg[1:2]))
		res[3] = byte(stringToInt(msg[2:3]))
		res[4] = byte(stringToInt(msg[3:4]))
	}
	//fmt.Println(res)
	return res
	
}
func loginSetUp(connection net.Conn) {
	fmt.Println(lines[0])
	
	for {
		fmt.Print(lines[1])
		card := strings.Replace(string(userInput()), " ", "", -1) //Läser in kortnr
		connection.Write(makeMsg(100, card))
		ans := make([]byte, 10)
		connection.Read(ans)
		if ans[0] == 253 {
			break
		} else {
			fmt.Println(lines[6])
		}
	}
	
	for {
		fmt.Print(lines[2])
		password := userInput()
		connection.Write(makeMsg(101, password))
	    
		ans := make([]byte, 10)
		connection.Read(ans)
		if ans[0] == 253 {
			break
		} else {
			fmt.Println(lines[6])
		}
	}
}

func decode(array []byte) string {
	var tmp string

	for _, elem := range array {
		tmp += strconv.Itoa(int(elem))
	}
	return tmp
}


func handlingRequests(connection net.Conn) {    
	for {
		fmt.Println(lines[3])
		fmt.Println(lines[4])
		fmt.Println(lines[5])
		
		input := strings.TrimSpace(userInput())
		switch input {
		case "1": //saldo
			cmd := makeMsg(1,"")
			connection.Write(cmd)
			saldo := make([]byte, 10)
			read,_ := connection.Read(saldo)
			if read != 0 {
				fmt.Println(decode(saldo[0:read]))
			}else {
				fmt.Println("0")
			}
		case "2": //withdraw
			var amount string
			for {
				fmt.Print(lines[7])
				amount = strings.TrimSpace(userInput())
				_, err := strconv.Atoi(amount)
				//fmt.Println(err, "error i 2")
				if len(amount) <= 27 && (len(amount) > 0) && (err == nil) {
					break
				} else {
					//fmt.Println(len(amount))
					fmt.Println(lines[6])
				}
			}
			connection.Write(makeMsg(2, amount))
			//fmt.Println("here i got")
			ans := make([]byte, 10)
			connection.Read(ans)
			if ans[0] == 253 {
				break
			} else {
				fmt.Print(lines[11])
			}
		case "3": //deposit
			var amount string
			for {
				fmt.Print(lines[7])
				amount = strings.TrimSpace(userInput())
				_, err := strconv.Atoi(amount)
				//fmt.Println(err, "error i 2")
				if len(amount) <= 27 && len(amount) > 0 && err == nil {
					break
				} else {
					//fmt.Println(len(amount))
					fmt.Println(lines[6])
				}
			}
			connection.Write(makeMsg(3, amount))
			//fmt.Println("here i got")
			ans := make([]byte, 10)
			connection.Read(ans)
			if ans[0] == 253 {
				break
			} else {
				fmt.Print(lines[11])
			}
		case "4": //exit
			connection.Write(makeMsg(4, ""))
			fmt.Println(lines[8])
			os.Exit(0)
		case "5": //change language
			validateLang()	
		default:
			fmt.Println(lines[10])
		}
	}
}

//client starts the client.
func client(){
    // Connect to the server through tcp/IP.
	connection, err := net.Dial("tcp", ("127.0.0.1" + ":" + "6666"))
    // If connection failed crash.
	check(err)
    //Configure the language.
    validateLang()
    //Time to log in to the account.
    loginSetUp(connection)
    handlingRequests(connection)
}
