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

	connections []net.Conn

	
	//updating = false
	mutex = sync.Mutex{}

	
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
	updateListener, err := net.Listen("tcp", ":" + strconv.Itoa(port + 1))

	
	go updater(updateListener)
	
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


func updater(listener net.Listener) {
	for {
		connection, err := listener.Accept()
		check(err)
		connections = append(connections, connection)
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

func updateDatabase() {
	var newDatabase string
	for _,user := range users {
		newDatabase += user.card_number + "\n"
		newDatabase += user.first_name + "\n"
		newDatabase += user.last_name + "\n"
		newDatabase += user.sifferkod + "\n"
		for _,codes := range user.enkod {
			newDatabase += codes + " "
		}
		newDatabase += "\n"
		newDatabase += user.saldo + "\n"
	}
	file, err := os.Create("databas.txt")
	if err != nil {
		fmt.Println("Master....I could not save the database.")
	}
	file.WriteString(newDatabase)
	file.Close()
}
		

func srvMaster(listener net.Listener){
	shutdown := false
	for shutdown == false {
		cmd := strings.TrimSpace(string(userInput()))
		switch cmd {
		case "shutdown":
			shutdown = true
			for _,cust := range masterList {
				cust.connection.Close()
			}
			listener.Close()
			updateDatabase()
		case "update" :
			fmt.Println("What lang do you wish to update?")
			lang := strings.TrimSpace(string(userInput())) //What langue master picked
			filename := lang + ".txt"
			content,_ := ioutil.ReadFile(filename)
			
			rest := make([]byte, 10-len(lang))
			for index := range rest {
				rest[index] = 32
			}

			
			sendlang := []byte(lang)
			sendlang = append(sendlang, rest...)

			fmt.Println(content)

			
			for _,cons := range connections { 
				//cons.connection.Write([]byte("hej"))
				write(cons, append([]byte{255}, append(sendlang, append([]byte(content), byte(4))...)...))
			}
		default:
			fmt.Println("Did not that understand command, please try again.")
		}	
	}	
}

func write(connection net.Conn, message []byte) {
	mutex.Lock()
	connection.Write(message)
	mutex.Unlock()
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
	tmp = fillup(tmp, 1, 10)
	return tmp
}

func fillup(array []byte, endIndex, startIndex int) []byte {

	for tmp := startIndex; tmp < endIndex; tmp++ {
		array[tmp] = byte(0)
	}

	return array
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
	fmt.Println("kom ut ur wait i srv", tmp, readFlag)
	return tmp, length
}

func removeZero(array []byte) []byte {
	tmp:= make([]byte, 0)
	for _,elem := range array {
		if elem == 0 {
			break
		}
		tmp = append(tmp, elem)
	}
	return tmp
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
			byteamount := removeZero(currentOperation[1:length])
			fmt.Println(byteamount, "byteamount in 2")
			
			amount,amounterr := strconv.Atoi(string(byteamount))
			fmt.Println(amounterr, amount)
			tmpCode := ""
			for index := range user.enkod {
				if user.enkod[index] != "-1" {
					tmpCode = strconv.Itoa(index + 1)
					break
				}
			}
			fmt.Println(tmpCode, "tmpCode")
			if tmpCode == "" {
				connection.Write(validate(false))
				break
			}else {
				connection.Write([]byte(tmpCode))
			}

			oneCodeInput := make([]byte, 10)
			connection.Read(oneCodeInput) //Read the given inputCode

			codeFound := false

			if oneCodeInput[0] != 103 {
				connection.Write(validate(false))
				break
			}

			oneCodeInput = removeZero(oneCodeInput[1:])
			fmt.Println(string(oneCodeInput))

			fmt.Println(oneCodeInput, "<-  oneCodeInput")
			for i := range user.enkod {
				fmt.Println(user.enkod[i])
				if user.enkod[i] == string(oneCodeInput) {
					user.enkod[i] = "-1" //remove the code
					codeFound = true
					connection.Write(validate(true))
				} 
			}
			//check if we found the code, else break and send false
			if codeFound == false {
				connection.Write(validate(false))
				break
			}
			
			user.mutex.Lock()
			//userSaldo,saldoerr := strconv.ParseInt(user.saldo, 10, 64)
			//fmt.Println(saldoerr)
			userSaldo,_ := strconv.Atoi(user.saldo)
			//saldoint,_ := strconv.ParseInt(string(userSaldo), 10, 64)
			//amountint,_ := strconv.ParseInt(string(amount), 10, 64)
			var checkSub int64
			checkSub =int64(userSaldo) - int64(amount)
			if checkSub >= 0 {
				connection.Write(validate(true))
				fmt.Println(checkSub, "<- userSal")
				user.saldo = strconv.FormatInt(checkSub, 10) 
			}else {
				fmt.Println("jag skriver ut error i userSlado-amount")
				connection.Write(validate(false))
			}
			user.mutex.Unlock()
		case 3 : //deposit
			byteamount := removeZero(currentOperation[1:length])
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
			var tmpSal int64
			userSaldo,_ := strconv.Atoi(user.saldo)
			tmpSal = int64(userSaldo) + int64(amount)
			user.saldo = strconv.FormatInt(tmpSal,10)
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
