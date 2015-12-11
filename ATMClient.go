package main
//ATMClient.go är skriven av Eric Hallström och Johannes Westlund
//Datum : 2015-12-11

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
	reader *bufio.Reader = bufio.NewReader(os.Stdin) //Read from terminal
	languages = make([]string,4) //Current preinstalled languages is just 4
	//standard output for choosing whatlanguage
	intro = "Please pick a language.\n" 
	end = "Please pick a real language.\n"

	//Lines hold all rows for output commands to clients.
	lines []string

	//custlang for each individual clients current language
	custlang string
)

//main function that will create the client.
func main() {
	//predef of 4 languages
	languages[0] = "english"
	languages[1] = "日本語"
	languages[2] = "deutsch"
	languages[3] = "svenska"

	//starting client
	client()
}

//Easyer way to check error, exit program if error
func check(err error) {
    if (err != nil) {
        fmt.Println(err)
        os.Exit(1)
    }
}

//read input from Stdin,
func validateLang() {
	
	//print all lang options
	fmt.Println(intro)
	for _, lang := range languages {
		fmt.Println(lang)
	}

	//local lines that will hold all lines in the choosen lang file
	lines = make([]string,0)
	//l will be our predefined languages
	l := languages

	//infinit loop for reading user input language, loop ends if correct lang was choosen 
	for {
		picked := strings.TrimSpace(userInput()) //read the input from Stdin, trim spaces 

		//looping through our predefined languages
		for _,lang := range(l) {
			fmt.Println(lang)
			if (picked == lang) {
				//fmt.Println("Found language!", lang, picked)
				custlang = lang //variable to hold current lang for client
				lang = lang + ".txt" //append .txt because we are reading from files
			
				file, err := os.Open(lang) //open the given languge file
				check(err) //check if correct

				scanner := bufio.NewScanner(file) //scanning through the file
				for scanner.Scan() {
					lines = append(lines, scanner.Text()) //appending each line in the given langue file to lines
				 }
				file.Close() //close the file
				check(scanner.Err()) //check for errors so nothing is left in the file to read
				break
			} 
		}
		//check so we actually got something in len, if we have, we have successfully changed language
		if (len(lines) != 0) {
			break
		} else {
			fmt.Println(end) //print error msg
		}
	}
}

//userInput reads from Stdin untill newline, if input was empty - append space to msg
func userInput() string {
	msg, err := reader.ReadBytes('\n')
	check(err)
	if len(msg) == 0 {
		msg = append(msg, 32)
	}
	return string(msg)
}


func stringToInt(strint string) int {
	tmp,_ := strconv.Atoi(strint)
	return tmp
}

//fills our return value from makeMsg up to 10 byte.
func fillup(array []byte, endIndex, startIndex int) []byte {

	for tmp := startIndex; tmp < endIndex; tmp++ {
		array[tmp] = byte(0)
	}
	return array
}

//makeMsg takes opt - code for what we want to send and a msg return value will always be 10 bytes
func makeMsg(opt int, msg string) []byte {
	
	msg = strings.TrimSpace(msg) //remove space from input
	var res = make([]byte, 10) //return array variable for what to send back to srv
	res[0] = byte(opt)    //opt code will always be on index zero
	
	switch opt {
	case 2 : //Withdrawl
		if len(msg) > 9 { //cant whithdrawl amounts more than length 9, 
			break
		}
		//convert input msg to bytes, each byte gets its own index in res
		for index := range msg {
			res[index+1] = byte(msg[index])
		}

		//if msg was less then 9 we fill upp the rest so we always send 10 bytes
		res = fillup(res, len(msg)+1, 10)
	case 3 : //deposit does same as case 2
		if len(msg) > 9 {
			break
		}
		for index := range msg {
			res[index+1] = byte(msg[index])
		}

		res = fillup(res, len(msg) +1, 10)
		
	case 100 : //cardnumber
		if len(msg) != 16 { //cardnumber must be 16 digits
			break
		}
		//each two digit gets it's own index in res to avoid when we are converintg numbers bigger then 255
		res[1] =  byte(stringToInt(msg[0:2]))
		res[2] =  byte(stringToInt(msg[2:4]))
		res[3] =  byte(stringToInt(msg[4:6]))
		res[4] =  byte(stringToInt(msg[6:8]))
		res[5] =  byte(stringToInt(msg[8:10]))
		res[6] =  byte(stringToInt(msg[10:12]))
		res[7] =  byte(stringToInt(msg[12:14]))
		res[8] =  byte(stringToInt(msg[14:16]))
		res = fillup(res, 9,10)
	case 101 : //password
		if len(msg) != 4 { //password must be length 4
			break	
		}
		//each digit in the password converts to bytes into res
		res[1] = byte(stringToInt(msg[0:1]))
		res[2] = byte(stringToInt(msg[1:2]))
		res[3] = byte(stringToInt(msg[2:3]))
		res[4] = byte(stringToInt(msg[3:4]))
		res = fillup(res, 5, 10)
	case 103 : //engångs koderna must be length 2 
		if len(msg) != 2 {
			break
		}
		res[1] = byte(msg[0])
		res[2] = byte(msg[1])
		res= fillup(res, 3, 10)
	}
	return res
}
//validate login setup from the user, needs cardnumber on 16 digits and sifferkod on 4 digits
//everything thats Write from loginSetup uses makeMsg to convert the msg to 10 byte array
//Each msg will start with an opt code and followed by different msg in byte 1-10
func loginSetUp(connection net.Conn) {
	fmt.Println(lines[0])  //first line in the given lang file

	//infinit loop to read in correct cardnumber
	for {
		fmt.Print(lines[1])
		card := strings.Replace(string(userInput()), " ", "", -1) //remooves space in cardnumber
		connection.Write(makeMsg(100, card))  //send msg to srv, first is opt code 100 for cardnumber to be validated
		ans := make([]byte, 10) //return value from srv
		connection.Read(ans) //read if it was validated from srv, 253 = true
		if ans[0] == 253 { 
			break
		} else {
			fmt.Println(lines[6]) //received opt code 252 = fail
		}
	}

	//infinit loop to read password
	for {
		fmt.Print(lines[2])
		password := userInput() //read input password from user
		connection.Write(makeMsg(101, password)) //makeMsg, opt code 101 for password, 
	    
		ans := make([]byte, 10)
		connection.Read(ans) //validate password from srv
		if ans[0] == 253 {
			break
		} else {
			fmt.Println(lines[6])
		}
	}
}

//decode converts byte array to string and returns it.
func decode(array []byte) string {
	var tmp string

	for _, elem := range array {
		tmp += strconv.Itoa(int(elem))
	}
	return tmp
}

//overrideFile is for when the srv is udpatning lang files
func overrideFile(filename, newcontent string) {
	file, fileErr := os.Create(filename) //open given file
	if fileErr != nil {
		fmt.Println(fileErr)
	}
	file.WriteString(newcontent) //write the new content to the file
	file.Close()
	return
}

//update the clients language file
func updateFile(connection net.Conn) {
	fmt.Println("Updating file...")
	lang := make([]byte, 10)
	var content string //content for the new file
	
	connection.Read(lang) //read what lang file to change

	
	
	tmp := make([]byte, 255) //tmp for holdning the content from the srv
	for {
		read,_ := connection.Read(tmp) //reads what the srv wants to update the file with
		if tmp[read-1] == 4 { //the update from the srv will end will 4 ascii = end of transmission
			break
		}
		content += string(tmp) //store the content of the file
	}
	

	//search for the input lang and if we can change it.
	exist := false
	for _,currentlang := range languages {
		if currentlang == strings.TrimSpace(string(lang)) {
			exist = true
		}
	}

	//the lang did exists! so we append the lang to the global languages list.
	if exist == false {
		languages = append(languages, strings.TrimSpace(string(lang)))
	}

	//We will override already existing languages and create a new file if its a totally new language
	filename := strings.TrimSpace(string(lang)) + ".txt"

	newcontent := string(content) //the content from the srv
	
	overrideFile(filename, newcontent) //overridefile will take all the content and replace all content with the new content


	// check if the clients current lang is the new lang
	if custlang == strings.TrimSpace(string(lang)) {
		tmplines := make([]string,0) //we must replace all old lines from the lang file
		file, err := os.Open(filename) //open the file
		check(err)
		
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			tmplines = append(tmplines, scanner.Text())
		}
		file.Close()
		*(&lines) = tmplines //replace all old lines with the new content 
		check(scanner.Err())
	}
	
}
/*
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
*/

//handlingRequests takes care of all users inputs from the client
func handlingRequests(connection net.Conn) {    
	for {
		//print instructions 
		fmt.Println(lines[3])
		fmt.Println(lines[4])
		fmt.Println(lines[5])
		
		input := strings.TrimSpace(userInput()) //read input from the user

		switch input {
		case "1": //saldo
			cmd := makeMsg(1,"") //opt code for checking saldo is 1
			connection.Write(cmd) //send opt code to srv
			saldo := make([]byte, 10) 
			read,_ := connection.Read(saldo) //srv returned saldo, read holds the length of the saldo
			//to avoid if saldo is zero...........
			if read != 0 { 
				fmt.Println(decode(saldo[0:read])) //turns saldo from bytearray to string
			}else {
				fmt.Println("0")
			}
		case "2": //withdraw
			var amount string //amount to withdraw
			for {
				fmt.Print(lines[7])
				amount = strings.TrimSpace(userInput()) // read amount to withdraw from user
				_, err := strconv.Atoi(amount) //check if amount was digits
				//cant withdraw more then 9 digits, cant withdraw 0 
				if len(amount) <= 9 && (len(amount) > 0) && (err == nil) {
					break
				} else {
					fmt.Println(lines[6]) // print error msg
				}
			}
			connection.Write(makeMsg(2, amount)) // write the amount to srv opt code 2 = withdrawl

			//srv wants a engångs kod from the client, srv stores a list of engångs koder.
			enKod := make([]byte, 10) // engångs koden 

			length, _ := connection.Read(enKod) //read what engångs kod the srv wants
			
			if enKod[0] == 252 { //if srv return inc enkod break and print error
				fmt.Println(lines[11])
				break
			}
			
			fmt.Println(lines[12] + " " + string(enKod[0:length])) //send it to the user

			inputCode := strings.TrimSpace(userInput()) // user enters engångs kod

			connection.Write(makeMsg(103, inputCode)) //Send inputCode to srv, opt code =  103 = withdrawl
			
			
			ans := make([]byte, 10)
			connection.Read(ans) //wait for srv to validate the engngskod
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
				// cant deposite more then 9 digits or 0
				if len(amount) <= 9 && len(amount) > 0 && err == nil {
					break
				} else {
					fmt.Println(lines[6])
				}
			}
			connection.Write(makeMsg(3, amount)) //write the amount to the srv to deposit

			ans := make([]byte, 10) //wait for validation from the srv
			connection.Read(ans)
			if ans[0] == 253 { // opt code 253 is true
				break
			} else {
				fmt.Print(lines[11])
			}
		case "4": //exit
			connection.Write(makeMsg(4, ""))
			fmt.Println(lines[8])
			os.Exit(0)
		case "5": //change language
			validateLang()	//validate new lang. 
		default:
			fmt.Println(lines[10]) //default print i did not understand that
		}
	}
}

//two different Dials, updateListener dial is constantly listing from srv updates 
func client(){
    // Connect to the server through tcp/IP.
	connection, err := net.Dial("tcp", ("127.0.0.1" + ":" + "9090"))
	updateListener , listErr := net.Dial("tcp", ("127.0.0.1" + ":" + "9091"))
	// If connection failed crash.
	check(err)
	check(listErr)
	//Create separate thread for updating  client.
	go update(updateListener)
	//Configure the language.
	validateLang()
	//Time to log in to the account.
	loginSetUp(connection)
	//handling requests from usr
	handlingRequests(connection)
}

// constantly listeing if srv sends opt ode 255, then we wants to update the lang file.
func update(connection net.Conn) {
	tmp := make([]byte, 1)
	for  {	 
		connection.Read(tmp)
		if tmp[0] == 255 {
			updateFile(connection)
		}
	}
}
