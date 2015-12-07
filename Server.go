package main

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
<<<<<<< HEAD
    _, err := client.Read(holder)

	if (err != nil) {
        return "", errors.New("Error couldn't get how many bytes that will be sent.")
	}
	
    bytes, _ := strconv.Atoi(string(holder))
	fmt.Println(string(holder[:]))
	fmt.Println(bytes)
=======
    number, err := client.Read(holder)
    if (err != nil) {
        //fmt.Println("Error couldn't get how many bytes that will be sent.")
        return "", errors.New("Error couldn't get how many bytes that will be sent.")
    }

    bytes, _ := strconv.Atoi(string(holder[0:number]))
    bytes++
    fmt.Println(bytes)
>>>>>>> f71d9cea94f596a3754aeeceff7df897031dcd6d
    message := ""   
    
    holder = make([]byte, 10)

    for (bytes != 0) {
        fmt.Println("Enter loop!")
        letters, err := client.Read(holder)
        fmt.Println("Read from client.")
        if (err != nil) {
            return "", errors.New("Error when reading from client.")
        }
<<<<<<< HEAD
	    fmt.Println(holder[:])
        message += string(holder[:])
=======
        fmt.Println(holder)
        message += string(holder[0:letters])
        fmt.Print(message)
>>>>>>> f71d9cea94f596a3754aeeceff7df897031dcd6d
        bytes--
    }
    return message, nil
}

func handleClient(client net.Conn) {
	for {
		msg := make([]byte,100)
		fmt.Println(client.Read(msg))
		message, err := read(client)
	    
		if (err != nil) {
			fmt.Println(err)
			client.Close()
			break
		} else {
			fmt.Println(message)
			_, errWrite := client.Write([]byte("Message recieved!"))
			if (errWrite != nil) { fmt.Println(errWrite) }
		}
		if (message == "stop\n") {
			break
		}
		if (message == "close\n") {
			client.Close()
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
