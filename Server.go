package main

// Import packages.
import (
    "fmt"
    "net"
    "strconv"
    "errors"
)

func server(port int) {
    fmt.Println("Server startup!")
    stringPort := strconv.Itoa(port)
    // start listening to a part
    listener, err := net.Listen("tcp", ":" + stringPort)
    if (err != nil) {
        fmt.Println("Couldn't listen on port " + stringPort)
        return
    }
    for {
        connection, err := listener.Accept()
        fmt.Println("Connection established!")

         go forceShutDown(listener, connection) 

        if (err != nil) {
            fmt.Println("Failed to establish connection.")
            break
        } else {
            go handleClient(connection)
        }
    }
}

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


func handleClient(client net.Conn) {
	for {
		message, err := read(client)
	    
		if (err != nil) {
			fmt.Println(err)
			client.Close()
			break
		} else {
			fmt.Print(message)
			_, errWrite := client.Write([]byte("Message recieved!"))
			if (errWrite != nil) { 
                fmt.Println(errWrite) 
            }
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
        server(port)
    }
}
