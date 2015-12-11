package main

// Import packages.
import (
	"bufio"
	"fmt"
	"os"
	"net"
	"strings"
	"strconv"
	//"bytes"
	//"io"
)

//Global reader that will read what the user writes.
var (
	reader *bufio.Reader = bufio.NewReader(os.Stdin)
	languages = make([]string,4)
	intro = "Please pick a language.\n"
	end = "Please pick a real language.\n"

	lines []string
	custlang string
)

//main function that will create the client.
func main() {
	languages[0] = "english"
	languages[1] = "日本語"
	languages[2] = "deutsch"
	languages[3] = "svenska"
	client()
}

func check(err error) {
    if (err != nil) {
        fmt.Println(err)
        os.Exit(1)
    }
}

func validateLang() {

	fmt.Println(intro)
	for _, lang := range languages {
		fmt.Println(lang)
	}
	lines = make([]string,0)
	l := languages
	for {
		//fmt.Println("Entering loop")
		picked := strings.TrimSpace(userInput())
		
		for _,lang := range(l) {
			fmt.Println(lang)
			if (picked == lang) {
				//fmt.Println("Found language!", lang, picked)
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


func overrideFile(filename, newcontent string) {
	file, fileErr := os.Create(filename)
	if fileErr != nil {
		fmt.Println(fileErr)
	}
	file.WriteString(newcontent)
	file.Close()
	return
}

func resetConn(connection net.Conn) {
	tmp:=make([]byte, 20)
	empty := false
	for empty == false {
		connection.Read(tmp)
		fmt.Println("spam")
		if tmp[0] == 0 {
			empty = true
		}
	}
	fmt.Println("Klar i resetConn!")
}

func updateFile(connection net.Conn) {
	fmt.Println("Updating file...")
	lang := make([]byte, 10)
	var content string
	
	connection.Read(lang)

	
	
	tmp := make([]byte, 255)
	for {
		read,_ := connection.Read(tmp)
		if tmp[read-1] == 4 {
			break
		}
		content += string(tmp)
	}
	
	//fmt.Println(lang, "<- lang")
	//fmt.Println(content, "<-- content")
	//fmt.Println( "Over is  content")
	
	exist := false
	
	for _,currentlang := range languages {
		if currentlang == string(lang) {
			exist = true
		}
	}
	if exist == false {
		languages = append(languages, string(lang))
	}

	filename := strings.TrimSpace(string(lang)) + ".txt"

	//fmt.Println(filename, "<- filename" ,[]byte(filename))

	newcontent := string(content)
	overrideFile(filename, newcontent)

	//fmt.Println(custlang, "<- custlang")
	//fmt.Println(strings.TrimSpace(string(lang)),"<- lang trimmad")

	
	if custlang == strings.TrimSpace(string(lang)) {
		tmplines := make([]string,0)
		file, err := os.Open(filename)
		check(err)
		
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			tmplines = append(tmplines, scanner.Text())
		}
	//	fmt.Println("Nu har vi scannat alla text", tmplines)
		file.Close()
		*(&lines) = tmplines
		check(scanner.Err())
	}
	
}

func wait(connection net.Conn) ([]byte,int) {
	fmt.Println("Nu är vi i wait")
	tmp := make([]byte, 1)
	rest := make([]byte, 9)
	length := 0
	readFlag := false
	for readFlag != true {
		length,_ = connection.Read(tmp)
		if tmp[0] != 0 {
			readFlag = true
			connection.Read(rest)
		}
	}
	tmp = append(tmp, rest...)
	fmt.Println("kom ut ur wait", tmp, readFlag)
	return tmp, length
}

func handlingRequests(connection net.Conn) {    
	for {
		fmt.Println(lines[3])
		fmt.Println(lines[4])
		fmt.Println(lines[5])
		
		input := strings.TrimSpace(userInput())

	//	fmt.Println(input, "input in handlingReq")
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

			fmt.Println("here i got")

			enKod, length := wait(connection)
			if enKod[0] == 252 {
				fmt.Println(lines[12])
				break
			}
			fmt.Println(lines[12] + " " + string(enKod[0:length])) //Read what onetime code srv wants
			fmt.Println("Nu ska jag har skrivit ut enKod", enKod)

			inputCode := strings.TrimSpace(userInput()) //Read user input

			connection.Write([]byte(inputCode)) //Send inputCode to srv
			
			
			ans := make([]byte, 10)
			connection.Read(ans)
			if ans[0] == 252 { //if srv return inc enkod break and print error
				fmt.Println(lines[11])
				break
			}

			ans = make([]byte, 10)
			connection.Read(ans)
			if ans[0] == 252 { //check if transaction was complemeted
				fmt.Println(lines[11]) 
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
	connection, err := net.Dial("tcp", ("127.0.0.1" + ":" + "8082"))
	updateListener , listErr := net.Dial("tcp", ("127.0.0.1" + ":" + "8083"))
    // If connection failed crash.
	check(err)
	check(listErr)
	//Create separate thread for uodating client.
	go update(updateListener)
    //Configure the language.
    validateLang()
    //Time to log in to the account.
    loginSetUp(connection)
    handlingRequests(connection)
}


func update(connection net.Conn) {
	tmp := make([]byte, 1)
	for  {	 
		connection.Read(tmp)
		if tmp[0] == 255 {
			fmt.Println(tmp[0], "<--tmp[0]")
			updateFile(connection)
		}
	}
}
