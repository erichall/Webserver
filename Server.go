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

	// Hold the connections.
	masterList []Customer
)

/*User struct definierar all information om en user.*/
type User struct {
	card_number string
	first_name string
	last_name string
	sifferkod string
	enkod []string
	saldo string
	mutex sync.Mutex
}


type Customer struct {
	connection net.Conn
	user *User
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
				user := loginSetup(connection)
				masterList = append(masterList, Customer{connection, user})
				handleClient(connection, user)
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
		/*case "banner":
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
			fmt.Println("Welcome message successfully changed.") */
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

func bytesToString(bytes []byte, start int, end int) string {
	var tmp string
	for ;start < end;start++ {
		tmp += string(bytes[start])
	}
	return tmp
}

func decode(array []byte) string {
	var tmp string

	for _, elem := range array {
		tmp += strconv.Itoa(int(elem))
	}
	return tmp
}

func validate(approved bool) []byte {
	tmp := make([]byte, 10)
	if approved {
		tmp[0] = 253
	} else {
		tmp[0] = 252
	}
	return tmp
}

func loginSetup(connection net.Conn) *User{	
	var card = make([]byte, 10)
	var userindex int
	var rightuser *User
	userFound := false
	
	for {
		connection.Read(card)//Läser in kortnr
		//fmt.Println(string(card[1:cardsize]))
		cardInt := decode(card[1:9])
		fmt.Println(cardInt)
		if card[0]  == 100 {
			//stringCard := bytesToString(card, 1, cardsize)
			for index := range(users) {
				fmt.Println(users[index].card_number)
				if(users[index].card_number == cardInt){
					connection.Write(validate(true))
					userindex = index
					userFound = true
					break
				}
			}
			if userFound == false {
				connection.Write(validate(false))
			}
		}else {
			connection.Write(validate(false))
		}

		if userFound {
			break
		}
	}
	
	var password = make([]byte, 10)
	for {
		connection.Read(password)
		sifferkod := decode(password[1:5])
		fmt.Println(sifferkod)
		if users[userindex].sifferkod == sifferkod {
			connection.Write(validate(true))
			rightuser = &users[userindex]
			break
		}else {
			connection.Write(validate(false))
		}
	}
	fmt.Println("Lets go to main meny")
	return rightuser
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

func stringToInt(strint string) int {
	tmp,_ := strconv.Atoi(strint)
	return tmp
}

func wait(connection net.Conn) ([]byte,int) {
	tmp := make([]byte, 10)
	length := 0
	readFlag := false
	for readFlag != true {
		 
		length,_ = connection.Read(tmp)
		if tmp[0] != 0 {
			readFlag = true
		}
	}
	fmt.Println("kom ut ur wait", tmp, readFlag)
	return tmp, length
}

func handleClient(connection net.Conn, user *User) {
	stillconnected := true
	for stillconnected {
		currentOperation, length := wait(connection)
		opcode := currentOperation[0]
		fmt.Println(opcode, "opcode")
		switch opcode {
		case 1 : //saldo
			tmpsaldo := (*user).saldo
			var res []byte
			fmt.Println(tmpsaldo, "det här är saldo")
			times := len(tmpsaldo)/2
			if times != 0 {
				for index := 0; index <= times; index += 2 {	
					res = append(res, byte(stringToInt(tmpsaldo[index:(index+2)])))
					fmt.Println(res, "inne i loopen")
				}
			}
			if len(tmpsaldo) % 2 != 0 {
				res = append(res, byte(stringToInt(tmpsaldo[(times):])))
			}
			fmt.Println(res, "sista raden i case1")
			connection.Write(res)
		case 2 : //whitdra
			byteamount := currentOperation[1:length]
			fmt.Println(byteamount, "byteamount")
			
			
			var tmp string
			for _, elem := range byteamount{
				if elem != 0 {
					tmp += string(elem)
				}
			}
			
			amount,_ := strconv.Atoi(tmp)
			
			user.mutex.Lock()
			userSaldo,_ := strconv.Atoi(user.saldo)
			//if user.enkod[randcode-1] == inputcode {
				if userSaldo - amount >= 0 {
					//user.enkod[randcode-1] = 1
					connection.Write(validate(true))
					userSaldo = userSaldo - amount
					user.saldo = strconv.Itoa(userSaldo) 
				}else {
					connection.Write(validate(false))
				}
			//} else {
			//	write(client, []byte("Wronge onetime code"))
			//}
			user.mutex.Unlock()
		case 3 : //deposit
			byteamount := currentOperation[1:length]
			fmt.Println(byteamount, "byteamount")
			
			var tmp string
			for _, elem := range byteamount{
				if elem != 0 {
					tmp += string(elem)
				}
			}
			
			amount,_ := strconv.Atoi(tmp)
			fmt.Println(amount)
			user.mutex.Lock()
			userSaldo,_ := strconv.Atoi(user.saldo)
			userSaldo = userSaldo + amount
			user.saldo = strconv.Itoa(userSaldo)
			user.mutex.Unlock()
			connection.Write(validate(true))
		case 4 : // Exit
			connection.Close()
			stillconnected = false
		default:
			connection.Write(validate(false))
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
			var intkodlist []string = make([]string, len(tmpenkod))
	

			for index, elem := range tmpenkod {
				intkodlist[index] = elem
			} 

			card_number := userdata[0]
			sifferkod :=  userdata[3]
			saldo := userdata[5]
			fmt.Println(saldo, userdata[5])
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
