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
	//"errors"
	"os"
	//"encoding/binary"
	"bufio"
	"strings"
	
)

var (
	users []User

	//Global variable that indicate how many language the client support.
	languages = "english\n日本語\ndeutsch"
	intro = "Please pick a language.\n言語を選択してください。\nBitte wählen Sie eine Sprache aus.\n"
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
		connection, err := listener.Accept() //Accept

		go forceShutDown(listener,connection)

		if (err != nil) {
			fmt.Println("Failed to establish connection.")
			break
		} else {
			go validateLang(connection)
		}
	}
}

func validateLang(connection net.Conn){
	write(connection,[]byte(intro + "\n" + languages))
	var lines []string
	l := strings.Split(languages, "\n")
	
	for {
		picked := strings.TrimSpace(read(connection))
		fmt.Println(picked)
		for _,lang := range(l) {
			fmt.Println(lang)
			if (picked == lang) {
				fmt.Println("Found language!", lang, picked)
				write(connection, []byte("approved"))
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
			write(connection, []byte("Please pick a real language"))
		}
	}
	user := loginSetup(connection,lines)
	handleClient(connection, lines, user)
}





func loginSetup(connection net.Conn, lines []string) User{	
	write(connection, []byte(lines[0]))
	write(connection, []byte(lines[1])) //Fråga efter kortnummer
	for {
		tmpcard := read(connection)
		cardnumber, err := strconv.Atoi(tmpcard)
		if err != nil {
			write(connection, []byte(lines[6] + "\n" + lines[1]))
		} else {
			for _, u := range(users) {
				if(u.card_number == cardnumber){
					write(connection, []byte(lines[2]))
					for {
					    
					    tmpkod := read(connection)
					    sifferkod, sifferError := strconv.Atoi(tmpkod)
					    if sifferError != nil {
						    write(connection, []byte(lines[6]+"\n" + lines[2]))
					    }else if u.sifferkod == sifferkod {
						    write(connection, []byte("approved"))
						    write(connection, []byte(lines[3]))
						    return u
					    }else {
						    write(connection, []byte(lines[6] + "\n" + lines[2]))
					    }
					}
				}
			}
			write(connection, []byte(lines[6] + "\n" + lines[1]))
		
		}
	}
	
	return *new(User)
	
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

func wait(connection net.Conn) {
	holder := make([]byte, 1)
	for {
		num, _ := connection.Read(holder)
	
		fmt.Println(num)
		if num != 0 {
			break
		}
	}
}

func handleClient(client net.Conn, lines []string, user User) {
	write(client, []byte(lines[4]))
	write(client, []byte(lines[5]))

	stillconnected := true
	for stillconnected {
		input := strings.TrimSpace(read(client))
		switch input {
		case "1" : //saldo
			tmpsaldo := strconv.Itoa(user.saldo)
			write(client, []byte(tmpsaldo))
		case "2" : //whitdraw
			write(client, []byte("Amount:"))
			amount,_ := strconv.Atoi(read(client))
			user.saldo = user.saldo - amount
		case "3" : //deposit
			write(client, []byte("Amount:"))
			amount,_ := strconv.Atoi(read(client))
			user.saldo = user.saldo + amount
		case "4" : // Exit
			write(client, []byte("GOOOOOOOOOOOOOOOOODBYE"))
			client.Close()
			stillconnected = false
		default:
			write(client, []byte("hmmm what u say"))
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
	    users = findUser()
	    server(port)
    }
}

func check(e error){
	if e != nil {
		fmt.Println(e)
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
