/*
Webserver
Koden är skriven av Johannes Westlund & Eric Hallström
Senast modifierad : 2015-12-11 
*/
package main

// Importera nödvändiga paket.
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

// Globala variabler som används.
var (
    // Lista som innehåller alla olika users (bank konton) som banken har.
	users []User
	//Global reader som används för att läsa det användaren skriver i terminalen.
	reader *bufio.Reader = bufio.NewReader(os.Stdin)
	// Lista som sparar alla kunder som är inloggade.
	masterList []Customer
    // Lista med alla updaterings-connections som är etablerade.
	connections []net.Conn
)

/*User struct definierar all information om en user.*/
type User struct {
	card_number string
	first_name string
	last_name string
    // Lösenord.
	sifferkod string
	enkod []string
	saldo string
    // Mutex för att synkronisera transaktioner.
	mutex sync.Mutex
}
/* Customer beskriver en användare som är uppkopplad till serven. */
type Customer struct {
	connection net.Conn
	user *User
}
//Initierar serven till att använda flera processorer.
func init() {
	cpus := runtime.NumCPU() //Få hur många processorer som kan anvädas.
    //Gör så programmet använder så många processorer den kan.	
	runtime.GOMAXPROCS(cpus) 
}

/*Hanterar servern, 
input : Vilken port srv ska ligga på*/
func server(port int) {
	fmt.Println("Server startup!")
	stringPort := strconv.Itoa(port) 

	// Börja lyssna på porten.
	listener, err := net.Listen("tcp", ":" + stringPort)
    // Börja lyssna på updateringsporten.
	updateListener, err := net.Listen("tcp", ":" + strconv.Itoa(port + 1))
	// Skapa tråd som ska uppdatera klienterna när en uppdatering finns.
	go updater(updateListener)
	
	if (err != nil) {
		fmt.Println("Couldn't listen on port " + stringPort)
		return
	}
	fmt.Println("Creating server master.")
	//Skapa en tråd som tar hand om input från den som håller i serven.
	go srvMaster(listener)

	for {
        //Skapar en connection med en klient.
		connection, err := listener.Accept() 

		if (err != nil) {
			fmt.Println("Failed to establish connection.")
			break
		} else {
            //Skapa separat tråd för att ta hand om klienten.
			go func(){
				user := loginSetup(connection)
                //Lägg till i global lista av användare.
				masterList = append(masterList, Customer{connection, user})
				handleClient(connection, user)
			}()			
		}
	}
}

//Lyssnar på uppdaterings proten och lägger till nya connections allt eftersom.
func updater(listener net.Listener) {
	for {
		connection, err := listener.Accept()
		check(err)
		connections = append(connections, connection)
	}
}
//Läser det användaren ("server master") skriver.
func userInput() []byte {
    msg, err := reader.ReadBytes('\n')
	check(err)
    if len(msg) == 0 {
        msg = append(msg, 32)
    }
    return msg
}

//Uppdaterar databasen med nya värden. Sker vid en "shutdown" av serven.
func updateDatabase() {
	var newDatabase string
    // Samlar in den nya informationen från kontona.
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
    //Skriver över databas filen med de nya värdena.
	file, err := os.Create("databas.txt")
	if err != nil {
		fmt.Println("Master....I could not save the database.")
	}
	file.WriteString(newDatabase)
	file.Close()
}
		
//Funktionen som lyssnar på "server mastern" kommandon och
//utför dem.
func srvMaster(listener net.Listener){
    //Sant om vi vill stänga av serven.
	shutdown := false
	for shutdown == false {
        // Läser in input.
		cmd := strings.TrimSpace(string(userInput()))
		switch cmd {
        // Om kommandot var "shutdown" stängs serven ner.
		case "shutdown":
			shutdown = true
			for _,cust := range masterList {
				cust.connection.Close()
			}
			listener.Close()
			updateDatabase()
        // Om kommandon var "update" uppdaterar serven klienterna.
		case "update" :
			fmt.Println("What lang do you wish to update?")
            // Vilket språk vi vill uppdatera.
			lang := strings.TrimSpace(string(userInput())) 
            //Formatering av uppdateringen.
			filename := lang + ".txt"
			content,_ := ioutil.ReadFile(filename)
			
			rest := make([]byte, 10-len(lang))
			for index := range rest {
				rest[index] = 32
			}
			
			sendlang := []byte(lang)
			sendlang = append(sendlang, rest...)

			fmt.Println(content)
			//Sender över alla uppdaterings connections uppdateringen.
			for _,cons := range connections { 
				write(cons, append([]byte{255}, append(sendlang, append([]byte(content), byte(4))...)...))
			}
		default:
			fmt.Println("Did not that understand command, please try again.")
		}	
	}	
}

//Funktion som skriver över conection.
func write(connection net.Conn, message []byte) {
	connection.Write(message)
}
//Decode dekodar kortnummer och lösenord.
func decode(array []byte) string {
	var tmp string
    //Gör om varje byte till en sträng.
	for _, elem := range array {
		tmp += strconv.Itoa(int(elem))
	}
	return tmp
}
//Validate ger en 10byte false eller true som används för att indikera till
//klienten om det som skickades accepteras eller inte. 
func validate(approved bool) []byte {
    //Skapar en 10byte array.
	tmp := make([]byte, 10)
	if approved {
        //Accepterades.
		tmp[0] = 253
	} else {
        //Accepterades inte.
		tmp[0] = 252
	}
    //Fyller upp 10byte listan med nollor.
	tmp = fillup(tmp, 1, 10)
	return tmp
}

//Fyller ut arrayer från startIndex till endIndex med nollor
//för att protokollet ska hålla.
func fillup(array []byte, endIndex, startIndex int) []byte {

	for tmp := startIndex; tmp < endIndex; tmp++ {
		array[tmp] = byte(0)
	}

	return array
}

//Funktion som tar hand om login.
func loginSetup(connection net.Conn) *User{
    //Array för att hålla kortnummer	
	var card = make([]byte, 10)
	var userindex int
	var rightuser *User
	userFound := false
	
	for {
        //Läser in kortnr
		connection.Read(card)
		//Tar ut bytsen som håller nummret.
		cardInt := decode(card[1:9])
        //Skriver ut till server master så han kan se.
		fmt.Println(cardInt)
        // Om rätt opkod har angivits.
		if card[0]  == 100 {
            //Letar i användarna för att se om någon har det kortnummret.
			for index := range(users) {
                //Om en användare har det kortnummret.
				if(users[index].card_number == cardInt){
                    //Skicka tillbaka att kortnummer accepterats.
					connection.Write(validate(true))
					userindex = index
					userFound = true
					break
				}
			}
            // Om inte kortnummret fanns, skicka tillbaka false.
			if userFound == false {
				connection.Write(validate(false))
			}
        //Om inte opkoden var rätt, skicka false.
		}else {
			connection.Write(validate(false))
		}
        //Om vi hittat en giltig användare, bryt ut ur for loopen.
		if userFound {
			break
		}
	}
	//Array för att hålla lösenordet.
	var password = make([]byte, 10)
	for {
        //Läser in lösenordet.
		connection.Read(password)
        //Decodar lösenordet.
		sifferkod := decode(password[1:5])
        //Skriver lösenord för server master att se.
		fmt.Println(sifferkod)
        //Om lösenordet matchar den användare vi håller på att logga in som.
		if users[userindex].sifferkod == sifferkod {
            //Skicka true.
			connection.Write(validate(true))
            //Skapa en pekar till usern.
			rightuser = &users[userindex]
			break
		}else {
            //Skickar false.
			connection.Write(validate(false))
		}
	}
    //Returnerar pekaren till användaren.
	return rightuser
}

//Funktionen stringToInt gör om en sträng till en int.
func stringToInt(strint string) int {
	tmp,_ := strconv.Atoi(strint)
	return tmp
}

//Funktionen wait läser konstant till den får ett
//meddelande som har en opkod.
func wait(connection net.Conn) ([]byte,int) {
	tmp := make([]byte, 10)
	readFlag := false
    //Loop som går tills ett meddelande som har en opkod läses.
	for readFlag != true {
		connection.Read(tmp)
		if tmp[0] != 0 {
			readFlag = true
		}
	}
	return tmp, length
}

//Funktionen removeZero tar bort onödiga nollor i arrayen
//som inte innehåller någon information.
func removeZero(array []byte) []byte {
	tmp:= make([]byte, 0)
	for _,elem := range array {
        //Om vi stöter på en nolla vet vi att resterande element
        // i listan också är nollor på grund av fillout().
		if elem == 0 {
			break
		}
		tmp = append(tmp, elem)
	}
	return tmp
}

//handleClient hanterar klientens requests.
func handleClient(connection net.Conn, user *User) {
    //Bool variable som säger om connection fortfarande finns.
	stillconnected := true
	for stillconnected {
        //Väntar på svar från klienten om vad den vill göra.
		currentOperation, length := wait(connection)
        //Bryter ut opkoden från meddelandet.
		opcode := currentOperation[0]
        //Baserat på opkoden kan vi avgöra vad klienten vill göra.
		switch opcode {
		case 1 : //saldo
            //Tar ut saldo.
			tmpsaldo := (*user).saldo
			var res []byte
            //Försöker dela upp saldot i storlek om 2 siffror och stoppar in dem
            //på de olika byte positionerna.
			times := len(tmpsaldo)/2
			if times != 0 {
				for index := 0; index <= times; index += 2 {	
					res = append(res, byte(stringToInt(tmpsaldo[index:(index+2)])))
				}
			}
			if len(tmpsaldo) % 2 != 0 {
				res = append(res, byte(stringToInt(tmpsaldo[(times):])))
			}
            //Skicka saldot till klienten.
			connection.Write(res)
		case 2 : //whitdra
            //Ta ut hur stort belopp man ska ta ut.
			byteamount := removeZero(currentOperation[1:length])
			//Konvertera till int.
			amount,_ := strconv.Atoi(string(byteamount))
            //Ta fram vilken engångskod som användaren ska ge.
			tmpCode := ""
			for index := range user.enkod {
				if user.enkod[index] != "-1" {
					tmpCode = strconv.Itoa(index + 1)
					break
				}
			}
            //Om alla koder har använts skicka false.
			if tmpCode == "" {
				connection.Write(validate(false))
				break
            //Annars skickar vi vilken engångskod användaren ska ge.
			}else {
				connection.Write([]byte(tmpCode))
			}
            //Läs in svaret som användaren gett.
			oneCodeInput := make([]byte, 10)
			connection.Read(oneCodeInput)
            
			codeFound := false
            //Om inte rätt opkod har getts i meddelandet som skickades.
            //Skicka false.
			if oneCodeInput[0] != 103 {
				connection.Write(validate(false))
				break
			}
            //Annars, ta bort onödiga byte-nollor som inte innehåller någon info.
			oneCodeInput = removeZero(oneCodeInput[1:])
            //Leta upp om engångskoden är rätt.
			for i := range user.enkod {
				if user.enkod[i] == string(oneCodeInput) {
                    //Om vi hittade den så tar vi bort den.
					user.enkod[i] = "-1"
					codeFound = true
                    //Skickar att engångskoden hittades.
					connection.Write(validate(true))
				} 
			}
			//Om vi inte hittade koden skickar vi false.
			if codeFound == false {
				connection.Write(validate(false))
				break
			}
			//Låser mutexen kopplat till kontot för att synkronisera.
			user.mutex.Lock()
            //Gör om saldot till en int.
			userSaldo,_ := strconv.Atoi(user.saldo)
			var checkSub int64
            //Subtraherar saldor med hur mycket vi vill ta ut.
			checkSub =int64(userSaldo) - int64(amount)
            //Om saldot inte är negativt.
			if checkSub >= 0 {
                //Skicka true
				connection.Write(validate(true))
                //Skriv tillbaka saldot i minnet.
				user.saldo = strconv.FormatInt(checkSub, 10) 
			}else {
                //Annars skicka false.
				connection.Write(validate(false))
			}
            //Lås upp mutexen så att andra trådar kan använda kontot.
			user.mutex.Unlock()
		case 3 : //deposit
            //Tar bort onödiga nollor som inte innehåller någon info.
			byteamount := removeZero(currentOperation[1:length])
			//Gör om det till en sträng. Notera att varje byte
            //innehåller info om hur mycket som ska lägga in.
			var tmp string
			for _, elem := range byteamount{
				if elem != 0 {
					tmp += string(elem)
				}
			}
			//Gör om till int.
			amount,_ := strconv.Atoi(tmp)
            //Låser mutex för synkronisering.
			user.mutex.Lock()
			var tmpSal int64
            //Konvertera saldot till int.
			userSaldo,_ := strconv.Atoi(user.saldo)
            //Updatera saldo.
			tmpSal = int64(userSaldo) + int64(amount)
            //Lägg in nytt saldo i minnet.
			user.saldo = strconv.FormatInt(tmpSal,10)
            //Lås upp mutex.
			user.mutex.Unlock()
            //Skicka true.
			connection.Write(validate(true))
		case 4 : // Exit
            //Stäng connection.
			connection.Close()
            //Gör så vi kommer ur loopen.
			stillconnected = false
		default:
            //Skicka false om opkoden inte känns igen.
			connection.Write(validate(false))
		}
    }

}

//main är mainfunktionen.
func main() {
    var port int
    //Läser in vilken port som serven ska lyssna på.
    _, err := fmt.Scanf("%d", &port)
    if (err != nil) {
        fmt.Println("Couldn't read user argument.")
    } else {
        //Ladda in alla användarna.
	    findUser()
        //Starta server.
	    server(port)
    }
}

//check kollar error och hanterar dem.
func check(e error){
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

//findUser
func findUser() {
    //Öppna databasen
    filen,err := os.Open("databas.txt")
	check(err)
	//Skapa en scanner till filen.
	scanner := bufio.NewScanner(filen)

	var stop_read int = 0
    //Hur många rader i databasen som består av en user.
	var user_info int = 6 
    //Lista som håller en user.	
	var userdata []string = make([]string, 6)
	//Itererar genom filen och samlar upp användarna en åt gången.
	for i := 1; scanner.Scan(); i++ {
		userdata[stop_read] = scanner.Text()
		stop_read++	
        //Om vi läst in en hel användare (6 rader).
		if i % user_info == 0 {
            //Ta ut engångskoderna.
			tmpenkod := strings.Split(userdata[4], " ")
			var intkodlist []string = make([]string, len(tmpenkod))
            //Stoppa in koderna i en lista.
			for index, elem := range tmpenkod {
				intkodlist[index] = elem
			} 
            //Ta ut andra parametrar.
			card_number := userdata[0]
			sifferkod :=  userdata[3]
			saldo := userdata[5]
			fmt.Println(saldo, userdata[5])
            //Skapa en user objekt.
			user := User{card_number,
				userdata[1],
				userdata[2],
				sifferkod,
				intkodlist,
				saldo,
				sync.Mutex{}}
		    //Lägg till den i globala listan av användare.
			users = append(users, user)
			stop_read = 0
			userdata = make([]string, 6)
		}
	}
}
