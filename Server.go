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
	"os"
	"bufio"
	"strings"
	"sync"
	"runtime"
	"io/ioutil"
)


var (
	users []User
	//Global reader that will read what the user writes.
	reader *bufio.Reader = bufio.NewReader(os.Stdin)
	languages = "english\n日本語\ndeutsch\nsvenska"
	intro = "Please pick a language.\n言語を選択してください。\nBitte wählen Sie eine Sprache aus.\nVar venlig velj ett sprak.\n"
	end = "Please pick a real language.\n実際の言語を選択してください。\nBitte wählen Sie eine echte Sprache.\nSnelle velj ett riktigt sprak.\n"
	// Hold the connections.
	masterList []Customer
)

/*User struct definierar all information om en user.*/
type User struct {
	card_number int
	first_name string
	last_name string
	sifferkod int
	enkod []int
	saldo int
	mutex sync.Mutex
}


type Customer struct {
	connection net.Conn
	user *User
	language string
	lines []string
}

func init() {
	cpus := runtime.NumCPU() //Get how many cpus the server comuter has.
	runtime.GOMAXPROCS(cpus) //Set the program to use all cpus for parallel computation.	
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
	fmt.Println("Creating server master.")
	go srvMaster(listener)

	for {
		connection, err := listener.Accept() //Accept
		//Create thread to end connection if established. 

		if (err != nil) {
			fmt.Println("Failed to establish connection.")
			break
		} else {
			go func(){
				lines, lang := validateLang(connection)
				fmt.Println(lang)
				user := loginSetup(connection,lines)
				masterList = append(masterList, Customer{connection, user, lang, lines})
				handleClient(connection, lines, user)
			}()
			

		}
	
	}
}

func userInput() []byte {
    msg, err := reader.ReadBytes('\n')
	check(err)
    if len(msg) == 0 {
        msg = append(msg, 32)
    }
    return msg
}

func srvMaster(listener net.Listener){
	for {
		cmd := strings.TrimSpace(string(userInput()))
		switch cmd {
		case "shutdown":
			for _,cust := range masterList {
				cust.connection.Close()
			}
			listener.Close()
		case "banner":
			fmt.Println("In what language are you going to write your banner?")
			bannerlang := strings.TrimSpace(string(userInput())) + ".txt"
			fmt.Println("What is your banner?")
			banner := strings.TrimSpace(string(userInput()))

			for index := range masterList {
				if masterList[index].language == bannerlang {
					*(&masterList[index].lines[14]) = banner
				}
			}
			overrideFile(bannerlang, banner, 14)
			fmt.Println("Banner successfully changed.")
		case "welcome message":
			fmt.Println("In what language are you going to write your welcome message?")
			welcomelang := strings.TrimSpace(string(userInput())) + ".txt"
			fmt.Println("What is your welcome message?")
			welcome := strings.TrimSpace(string(userInput()))
			overrideFile(welcomelang, welcome, 3)
			fmt.Println("Welcome message successfully changed.")

		case "randompic" :
			for _, usr := range masterList {
				writeAsciiPic(usr.connection)
			}
		default:
			fmt.Println("Did not that understand command, please try again.")
		}	
	}	
}

func overrideFile(filename, replacement string, row int) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	rows := strings.Split(string(content), "\n")
	rows[row] = replacement

	file, fileErr := os.Create(filename)
	if fileErr != nil {
		fmt.Println(fileErr)
	}
	newContent := strings.Join(rows, "\n")
	file.WriteString(newContent)
	file.Close()
}

func validateLang(connection net.Conn) ([]string, string) {
	write(connection,[]byte(intro + "\n" + languages))
	var lines []string
	l := strings.Split(languages, "\n")
	var custlang string
	for {
		fmt.Println("Entering loop")
		picked := strings.TrimSpace(read(connection))
		fmt.Println(picked)
		for _,lang := range(l) {
			fmt.Println(lang)
			if (picked == lang) {
				fmt.Println("Found language!", lang, picked)
				write(connection, []byte("approved"))
				lang = lang + ".txt"
				custlang = lang
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
			write(connection, []byte(end))
		}
	}
	return lines, custlang
}

func loginSetup(connection net.Conn, lines []string) *User{	
	write(connection, []byte(lines[0]))
	write(connection, []byte(lines[1])) //Fråga efter kortnummer
	for {
		tmpcard := read(connection)//Läser in kortnr
		cardnumber, err := strconv.Atoi(tmpcard)
		if err != nil {
			write(connection, []byte(lines[6] + "\n" + lines[1]))
		} else {
			for index := range(users) {
				if(users[index].card_number == cardnumber){
					write(connection, []byte("approved"))
					write(connection, []byte(lines[2]))
					for {
					    tmpkod := read(connection)
					    sifferkod, sifferError := strconv.Atoi(tmpkod)
					    if sifferError != nil {
						    write(connection, []byte(lines[6]+"\n" + lines[2]))
					    }else if users[index].sifferkod == sifferkod {
						    write(connection, []byte("approved"))
						    write(connection, []byte(lines[3])) //welcome msg
						    return &users[index]
					    }else {
						    write(connection, []byte(lines[6] + "\n" + lines[2]))
					    }
					}
				}
			}
			write(connection, []byte(lines[6] + "\n" + lines[1]))
		}
	}
	return new(User)
	
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

func handleClient(client net.Conn, lines []string, user *User) {
	stillconnected := true
	for stillconnected {
		write(client, []byte(lines[4]))
		write(client, []byte(lines[14]))
		write(client, []byte(lines[5]))
		input := strings.TrimSpace(read(client))
		
		switch input {
		case "1" : //saldo
			tmpsaldo := strconv.Itoa((*user).saldo)
			write(client, []byte(tmpsaldo))
		case "2" : //whitdra
			randcode := -1
			for index := range user.enkod {
				if user.enkod[index] != 1 {
					randcode = index + 1
					break
				}
			}
			if randcode == -1 {
				write(client, []byte(lines[13]))
				break
			} else {
				write(client, []byte("approved"))
			}
			stringcode := strconv.Itoa(randcode)
			
			write(client, []byte(lines[7]))
			amount,_ := strconv.Atoi(read(client))
			write(client, []byte(lines[12] +" " + stringcode))
			inputcode,_ := strconv.Atoi(strings.TrimSpace(read(client)))
			//fmt.Println(inputcode, err, user.enkod[randcode-1])
			user.mutex.Lock()

			if user.enkod[randcode-1] == inputcode {
				if user.saldo - amount >= 0 {
					user.enkod[randcode-1] = 1
					write(client, []byte("approved"))
					user.saldo = (*user).saldo - amount
				}else {
					write(client, []byte(lines[11]))
				}
			} else {
				write(client, []byte("Wronge onetime code"))
			}
			user.mutex.Unlock()
			fmt.Println("Back to main menu")
		case "3" : //deposit
			write(client, []byte(lines[7]))
			amount,_ := strconv.Atoi(read(client))
			user.mutex.Lock()
			user.saldo = (*user).saldo + amount
			user.mutex.Unlock()
		case "4" : // Exit
			write(client, []byte(lines[8] +" " +user.first_name + " " +lines[9]))
			client.Close()
			stillconnected = false
		case "5" : //byt språk
			tmp := &lines
			fmt.Println("Entering language config")
			var lang string
			*tmp, lang = validateLang(client)
			for index := range masterList {
				if masterList[index].connection == client {
					masterList[index].language = lang
					masterList[index].lines = *tmp
					break
				}
			}
		default:
			write(client, []byte(lines[10]))
		}
		
    }

}



func main() {
    var port int
    _, err := fmt.Scanf("%d", &port)
    if (err != nil) {
        fmt.Println("Couldn't read user argument.")
    } else {
	    findUser()
	    server(port)
    }
}

func check(e error){
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

func findUser() {
	//filen,err := os.Open("/home/erkan/Desktop/CDATA/PROGP/WEBSERVER/databas.txt")
    filen,err := os.Open("databas.txt")
	check(err)
	
	scanner := bufio.NewScanner(filen)

	
	fmt.Println(scanner.Text())


	var stop_read int = 0
	var user_info int = 6 //Hur många rader i databasen som består av en user.
	
	var userdata []string = make([]string, 6)

	
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
				saldo,
				sync.Mutex{}}
		
			users = append(users, user)
			stop_read = 0
			userdata = make([]string, 6)
		}
	}
}


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


func writeAsciiPic(connection net.Conn){
	bilden, err := ioutil.ReadFile("ascii.txt")
	check(err)
	write(connection, []byte(string(bilden)))
	
}
